package keyservice

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound             = errors.New("Not found")
	ErrExpired              = errors.New("Expired")
	ErrInvalidEncryptedData = errors.New("Invalid encrypted data")
	ErrSignatureError       = errors.New("Signature error")
)

var keyExpireTime = 1 * time.Minute

// Key contains data of key
type Key struct {
	Version          uint16 `json:"n,omitempty"` // 密钥版本
	Value            string `json:"v,omitempty"` // 当前生效的密钥
	ValueWillExpired string `json:"o,omitempty"` // 即将过期的密钥
}

func (k *Key) Copy() *Key {
	return &Key{
		Version: k.Version,
		Value: k.Value,
		ValueWillExpired: k.ValueWillExpired,
	}
}

// KeyService provides cipher service
type KeyService struct {
	seedKey []byte
	storage Storage
	cache   *cacheWrapper

	signatureSize int
	versionSize   int
}

// NewKeyService return new *KeyService
// seedKey is the basic key for encrypt other keys
func NewKeyService(seedKey string, storage Storage, cache Cache) *KeyService {
	return &KeyService{
		seedKey:       []byte(seedKey),
		storage:       storage,
		cache:         newCacheWrapper(cache, newCipher(seedKey), keyExpireTime),
		versionSize:   2,
		signatureSize: 4,
	}
}

// 根据Key ID批量获取远程key
// map[keyID]key
func (sv *KeyService) GetKeys(ids []string) (keys map[string]*Key, err error) {
	keys = make(map[string]*Key)

	needToRefreshIDs := make([]string, 0)

	for _, id := range ids {
		if key, err := sv.cache.Load(id); err == nil || err == ErrExpired {
			keys[id] = key
			if err == ErrExpired {
				needToRefreshIDs = append(needToRefreshIDs, id)
			}
		} else {
			needToRefreshIDs = append(needToRefreshIDs, id)
		}
	}

	if len(needToRefreshIDs) > 0 {
		var rkeys map[string]*Key
		logger.Debugf("load keys=%s from storage", strings.Join(needToRefreshIDs, ","))
		rkeys, err = sv.storage.LoadMany(needToRefreshIDs)
		if err == nil {
			for id, key := range rkeys {
				keys[id] = key
				sv.cache.Store(id, key)
			}
		} else {
			logger.Errorf("load keys=%s from storage err=%v", strings.Join(needToRefreshIDs, ","), err)
		}
	}

	return
}

func (sv *KeyService) GetKey(id string) (key *Key) {
	keys, _ := sv.GetKeys([]string{id})

	if len(keys) > 0 {
		return keys[id]
	}

	return nil
}

// Encrypt encrypt content with key specified by keyID
func (sv *KeyService) Encrypt(content string, keyID string) (ret string, err error) {
	key := sv.GetKey(keyID)

	if key == nil {
		err = ErrNotFound
		return
	}

	cipher := newCipher(key.Value)
	cipherBytes, err := cipher.Encrypt([]byte(content))

	if err != nil {
		return
	}

	// cipher data byte size
	cl := len(cipherBytes)
	// seed key byte size
	kl := len(sv.seedKey)
	// signature byte size
	sl := sv.signatureSize
	// version bytes size
	vl := sv.versionSize
	// sig(4 bytes)+verion(2 bytes)+cipherBytes+seedKey
	bs := make([]byte, sl+vl+cl+kl)
	copy(bs[sl+vl:], cipherBytes)
	copy(bs[sl+vl+cl:], sv.seedKey)
	bs[sl] = byte(key.Version >> 8)
	bs[sl+1] = byte(key.Version & 0xff)
	sig := md5.Sum(bs[sl:])
	copy(bs[0:sl], shortSignature(sig, int(sl)))

	ret = base64.RawURLEncoding.EncodeToString(bs[0 : sl+vl+cl])

	return
}

func shortSignature(whole [16]byte, size int) []byte {
	if size > len(whole) {
		return whole[0:]
	}
	start := (len(whole) - size) / 2
	return whole[start : start+size]
}

func (sv *KeyService) Decrypt(content string, keyID string) (ret string, err error) {
	key := sv.GetKey(keyID)

	if key == nil {
		err = ErrNotFound
		return
	}

	cipherData, err := base64.RawURLEncoding.DecodeString(content)
	if err != nil {
		return
	}

	sl := sv.signatureSize
	vl := sv.versionSize
	kl := len(sv.seedKey)
	cl := len(cipherData)
	if cl < sl+vl+1 {
		err = ErrInvalidEncryptedData
		return
	}

	bs := make([]byte, cl+kl)
	copy(bs, cipherData)
	copy(bs[cl:], sv.seedKey)
	sig := md5.Sum(bs[sl:])

	if !bytes.Equal(bs[0:sl], shortSignature(sig, sl)) {
		err = ErrSignatureError
		return
	}

	version := (uint16(bs[sl]) << 8) | uint16(bs[sl+1])
	keyValue := key.Value
	if version != key.Version {
		// is previous version key
		if version == key.Version-1 && key.ValueWillExpired != "" {
			keyValue = key.ValueWillExpired
		}
	}

	cipher := newCipher(keyValue)

	s, err := cipher.Decrypt(cipherData[sl+vl:])

	if err != nil {
		return
	}

	ret = string(s)

	return
}

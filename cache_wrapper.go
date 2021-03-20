package keyservice

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type cacheItem struct {
	Value      *Key      `json:"v"`
	Expiration time.Time `json:"e"`
}

func (c *cacheItem) isExpired() bool {
	return time.Now().After(c.Expiration)
}

type cacheWrapper struct {
	cipher Cipher
	cache  Cache

	expireTime time.Duration
	buffer     sync.Map
}

func newCacheWrapper(cache Cache, cipher Cipher, cacheTime time.Duration) *cacheWrapper {
	return &cacheWrapper{
		cache:      cache,
		cipher:     cipher,
		expireTime: cacheTime,
	}
}

func (w *cacheWrapper) Load(id string) (key *Key, err error) {
	if v, ok := w.buffer.Load(id); ok {
		item := v.(*cacheItem)
		key = item.Value
		if item.isExpired() {
			err = ErrExpired
		}
		return
	}

	cipherData, err := w.cache.Load(id)
	if err != nil {
		return nil, err
	}
	if len(cipherData) == 0 {
		return nil, ErrNotFound
	}

	data, err := w.cipher.Decrypt(cipherData)

	if err != nil {
		err = errors.Wrap(err, "cacheWrapper.Load decrypt")
		logger.Error(err)
		return
	}

	val := &cacheItem{}
	err = json.Unmarshal(data, val)
	if err != nil {
		err = errors.Wrap(err, "cacheWrapper.Load unmarshal")
		logger.Error(err)
		return
	}

	w.buffer.Store(id, val)

	key = val.Value
	if val.isExpired() {
		err = ErrExpired
	}

	logger.Debugf("load cache %s", id)

	return key, err
}

func (w *cacheWrapper) Store(id string, key *Key) error {
	item := &cacheItem{
		Value:      key,
		Expiration: time.Now().Add(w.expireTime),
	}

	w.buffer.Store(id, item)

	logger.Debugf("store cache %s", id)

	data, err := json.Marshal(item)
	if err != nil {
		err = errors.Wrap(err, "cacheWrapper.Store marshal")
		logger.Error(err)
		return err
	}

	cipherData, err := w.cipher.Encrypt(data)
	if err != nil {
		err = errors.Wrap(err, "cacheWrapper.Store encrypt")
		logger.Error(err)
		return err
	}

	if err = w.cache.Store(id, cipherData); err != nil {
		err = errors.Wrap(err, "cacheWrapper.Store store")
		logger.Error(err)
		return err
	}

	return nil
}

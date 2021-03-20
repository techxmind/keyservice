package keyservice

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	ErrMethodNotImplemented = errors.New("method not implemented")
)

type Storage interface {
	Store(id string, key *Key) error
	LoadMany(ids []string) (map[string]*Key, error)
}

type FileStorage struct {
	mu sync.RWMutex
	path string
	secret string
	data map[string]*Key
	lastModified time.Time
}

func NewFileStorage(path string, secret string) (*FileStorage, error) {
	s := &FileStorage{
		path: path,
		secret: secret,
		data: make(map[string]*Key),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *FileStorage) Store(id string, key *Key) error {
	return ErrMethodNotImplemented
}

func (s *FileStorage) LoadMany(ids []string) (ret map[string]*Key, err error) {
	if s.modified() {
		if err := s.load(); err != nil {
			logger.Errorf("reload keystore file err:%s", err.Error())
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret = make(map[string]*Key)
	for _, id := range ids {
		if s.data[id] != nil {
			ret[id] = s.data[id].Copy()
		}
	}
	return ret, nil
}

func (s *FileStorage) modified() bool {
	info, err := os.Stat(s.path)
	if err != nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastModified.Before(info.ModTime())
}

func (s *FileStorage) load() error {
	info, err := os.Stat(s.path)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}
	cipherData, err := base64.RawURLEncoding.DecodeString(string(contents))
	if err != nil {
		return err
	}
	contents, err = AesDecrypt(cipherData, []byte(s.secret))
	if err != nil {
		return err
	}
	data := map[string]*Key{}
	if err = json.Unmarshal(contents, &data); err != nil {
		return err
	}
	s.mu.Lock()
	s.data = data
	s.lastModified = info.ModTime()
	s.mu.Unlock()
	return nil
}
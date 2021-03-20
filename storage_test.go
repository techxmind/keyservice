package keyservice

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var (
	_testKeyId1 = "test-key-1"
	_testKeyId2 = "test-key-2"
	_testKey1   = &Key{
		Value:   "test-key-1-value",
		Version: 1,
	}
	_testKey2 = &Key{
		Value:            "test-key-2-value-new",
		Version:          2,
		ValueWillExpired: "test-key-2-value",
	}
)

type testStorage struct {
	data            sync.Map
	loadStat        map[string]int
	loadErrNextTime bool
}

func newTestStorage() *testStorage {
	s := &testStorage{
		loadStat: make(map[string]int),
	}

	s.Store(_testKeyId1, _testKey1)
	s.Store(_testKeyId2, _testKey2)

	return s
}

func (s *testStorage) Store(id string, key *Key) error {
	s.data.Store(id, key)
	return nil
}

func (s *testStorage) LoadMany(ids []string) (ret map[string]*Key, err error) {
	ret = make(map[string]*Key)

	for _, id := range ids {
		s.loadStat[id] += 1
		if v, ok := s.data.Load(id); ok {
			ret[id] = v.(*Key)
		}
	}

	if s.loadErrNextTime {
		s.loadErrNextTime = false
		return nil, errors.New("loadmany error cause by test flag")
	}

	return ret, nil
}

func writeFileKeyStoreContents(file string, data map[string]*Key, secret string) (err error) {
	contents, err := json.Marshal(data)
	if err != nil {
		return
	}
	contents, err = AesEncrypt(contents, []byte(secret))
	if err != nil {
		return
	}
	contents = []byte(base64.RawURLEncoding.EncodeToString(contents))
	err = ioutil.WriteFile(file, contents, 0777)
	return
}

func TestNewFileStorage(t *testing.T) {
	file := filepath.Join(os.TempDir(), fmt.Sprintf("ks-%d", time.Now().UnixNano()+rand.Int63n(1000)))
	secret := "secret-key"
	data := map[string]*Key{
		"key1" : &Key{
			Version: 1,
			Value: "key1-value",
			ValueWillExpired: "",
		},
		"key2" : &Key{
			Version: 2,
			Value: "key2-value-2",
			ValueWillExpired: "key2-value-1",
		},
	}
	err := writeFileKeyStoreContents(file, data, secret)
	require.NoError(t, err)
	defer os.Remove(file)

	s, err := NewFileStorage(file, secret)
	require.NoError(t, err)
	keys, err := s.LoadMany([]string{"key1", "key2"})
	require.NoError(t, err)
	assert.EqualValues(t, keys, data)

	// modify file
	data["key1"] = &Key{
		Version: 2,
		Value: "key1-value-2",
		ValueWillExpired: "key1-value",
	}
	err = writeFileKeyStoreContents(file, data, secret)
	require.NoError(t, err)
	keys, err = s.LoadMany([]string{"key1", "key2"})
	require.NoError(t, err)
	assert.EqualValues(t, keys, data)
}
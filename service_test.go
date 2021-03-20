package keyservice

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	originKeyExpireTime := keyExpireTime
	keyExpireTime = 5 * time.Millisecond
	defer func() {
		keyExpireTime = originKeyExpireTime
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := newTestStorage()
	c := newTestCache()
	sv := NewKeyService(
		fmt.Sprintf("%f", r.Float64()),
		s,
		c,
	)

	keyNotExistsId := "key-not-exists-id"
	keys, err := sv.GetKeys([]string{_testKeyId1, _testKeyId2, keyNotExistsId})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keys))
	assert.EqualValues(t, _testKey1, keys[_testKeyId1])
	assert.EqualValues(t, _testKey2, keys[_testKeyId2])
	assert.Equal(t, 1, s.loadStat[_testKeyId1])
	assert.Equal(t, 1, s.loadStat[_testKeyId2])

	assert.Equal(t, 1, c.loadStat[_testKeyId1])
	assert.Equal(t, 0, c.loadHitStat[_testKeyId1])
	assert.Equal(t, 1, c.storeStat[_testKeyId1])

	key := sv.GetKey(_testKeyId1)
	assert.NotNil(t, key)
	assert.Equal(t, 1, c.loadStat[_testKeyId1])
	assert.Equal(t, 0, c.loadHitStat[_testKeyId1])
	assert.Equal(t, 1, s.loadStat[_testKeyId1])

	s.loadErrNextTime = true
	time.Sleep(keyExpireTime + 1*time.Millisecond)
	key = sv.GetKey(_testKeyId1)
	assert.NotNil(t, key)
	assert.Equal(t, 2, s.loadStat[_testKeyId1])
	assert.Equal(t, 1, c.storeStat[_testKeyId1])

	str := "hello,world!"
	encrypted, err := sv.Encrypt(str, _testKeyId1)
	assert.NotEqual(t, str, encrypted)
	t.Logf("encrypted by key[%s] %s => %s", _testKeyId1, str, encrypted)
	decrypted, err := sv.Decrypt(encrypted, _testKeyId1)
	assert.Equal(t, str, decrypted)

	decrypted, err = sv.Decrypt(encrypted+"?", _testKeyId1)
	assert.NotNil(t, err)
	t.Log(err)

	decrypted, err = sv.Decrypt(encrypted, _testKeyId2)
	assert.NotNil(t, err)
	t.Logf("decrypted:%s, err:%v", decrypted, err)
}

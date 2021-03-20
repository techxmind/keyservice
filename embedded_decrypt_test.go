package keyservice

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecryptEmbeddedString(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := newTestStorage()
	c := newTestCache()
	sv := NewKeyService(
		fmt.Sprintf("%f", r.Float64()),
		s,
		c,
	)

	secret := "hello"
	encrypted, err := sv.Encrypt(secret, _testKeyId1)
	assert.Nil(t, err)

	str := fmt.Sprintf(`{
		"secret" : "$$%s$$",
		"secret" : "$$%s$$",
	}`, encrypted, encrypted)

	decodeStr, err := sv.DecryptEmbeddedString(str, _testKeyId1)
	assert.Nil(t, err)

	expectedStr := strings.Replace(str, fmt.Sprintf("$$%s$$", encrypted), secret, -1)
	assert.Equal(t, expectedStr, decodeStr)
}

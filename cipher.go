package keyservice

type Cipher interface {
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
}

type stdCipher struct {
	key []byte
}

func newCipher(key string) Cipher {
	return &stdCipher{
		key: []byte(key),
	}
}

func (c *stdCipher) Decrypt(data []byte) ([]byte, error) {
	return AesDecrypt(data, c.key)
}

func (c *stdCipher) Encrypt(data []byte) ([]byte, error) {
	return AesEncrypt(data, c.key)
}

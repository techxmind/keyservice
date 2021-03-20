package keyservice

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"

	"github.com/pkg/errors"
)

func randomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

//benchmark: 955 ns/op
//性能一般，适合在一些对信息安全要求非常高的场景
func AesEncrypt(text, rawKey []byte) (ciphertext []byte, err error) {
	/*
	   $key = substr(sha1($key, true), 0, 16);
	       $iv = openssl_random_pseudo_bytes(16);
	       $ciphertext = openssl_encrypt($plaintext, 'AES-128-CBC', $key, OPENSSL_RAW_DATA, $iv);
	       $ciphertext_base64 = urlsafe_b64encode($iv.$ciphertext);
	       return $ciphertext_base64;
	*/
	var key, iv []byte
	hash := sha1.Sum(rawKey)
	key = (hash[:])[:aes.BlockSize]

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	iv, err = randomBytes(aes.BlockSize)
	if err != nil {
		return
	}

	blockSize := cipherBlock.BlockSize()
	textBytes := pKCS5Padding(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(cipherBlock, iv)
	cipherBytes := make([]byte, len(textBytes))
	blockMode.CryptBlocks(cipherBytes, textBytes)

	return append(iv, cipherBytes...), nil
}

//benchmark: 607 ns/op
func AesDecrypt(ciphertext, rawKey []byte) (text []byte, err error) {

	if len(ciphertext) < aes.BlockSize {
		err = errors.New("ciphertext too short")
		return
	}

	hash := sha1.Sum(rawKey)
	key := (hash[:])[:aes.BlockSize]

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	var cipherBlock cipher.Block
	cipherBlock, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	blockMode := cipher.NewCBCDecrypter(cipherBlock, iv)
	if len(ciphertext)%blockMode.BlockSize() != 0 {
		err = errors.New("crypto/cipher: input not full blocks")
		return
	}

	blockMode.CryptBlocks(ciphertext, ciphertext)

	text, err = pKCS5UnPadding(ciphertext)

	return
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if unpadding > length {
		return nil, errors.New("invalid data")
	}
	return origData[:(length - unpadding)], nil
}

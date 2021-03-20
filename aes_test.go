package keyservice

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesDecrypt(t *testing.T) {
	ast := assert.New(t)

	rawKey := []byte("world")
	text := []byte("hello")
	ciphertextBase64 := "wpj36h7TfJvkYfPN8hQ45zdvhJhSZ7XI7DloHb9IclA"
	ciphertext, _ := base64.RawURLEncoding.DecodeString(ciphertextBase64)
	result, err := AesDecrypt(ciphertext, rawKey)
	ast.Equal(nil, err)
	ast.Equal(text, result)
}

func TestAesEncrypt(t *testing.T) {
	ast := assert.New(t)

	rawKey := []byte("world")
	text := []byte("this is a long text contains 中文。")

	ciphertext, err := AesEncrypt(text, rawKey)
	ast.Equal(nil, err)
	decodetext, derr := AesDecrypt(ciphertext, rawKey)
	ast.Equal(nil, derr)
	ast.Equal(text, decodetext)

	rawKey2 := []byte("world2")
	ciphertext, _ = AesEncrypt(text, rawKey2)
	decodetext, derr = AesDecrypt(ciphertext, rawKey)
	if derr == nil && bytes.Equal(decodetext, text) {
		t.Errorf("decrypt err, return %v %v", decodetext, derr)
	}

	decodetext, derr = AesDecrypt(append(ciphertext, '?'), rawKey2)
	if derr == nil && bytes.Equal(decodetext, text) {
		t.Errorf("decrypt err, return %v %v", decodetext, derr)
	}
}

func BenchmarkAesDecrypt(b *testing.B) {
	key := []byte("thisistestkey")
	ciphertext, _ := base64.RawURLEncoding.DecodeString("XJorU1U1wyaSEef2tdLx5U--17mpjlQ1IqCymZeKXmKuwtx0uPOqcN91RTqUjsuwIOG2QX1jJ9JZzR6DxReDK9hl8sIEX5ieHiJ8rvNAoCiQlf8tFIgq1aOwn_8mEYys")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := AesDecrypt(ciphertext, key)
			if err != nil {
				b.Errorf("aes decrypt err:%v", err)
			}
		}
	})
}

func BenchmarkAesEncrypt(b *testing.B) {
	var key = []byte("thisistestkey")
	var text = []byte("空山不见人,但闻人语响。 返影入深林,复照青苔上")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := AesEncrypt(text, key)
			if err != nil {
				b.Errorf("aes encrypt err:%v", err)
			}
		}
	})
}

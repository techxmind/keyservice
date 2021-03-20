package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/techxmind/keyservice"
	"io/ioutil"
	"os"
)

var (

)

func main() {
	ksFile := flag.String("keystore", "", "keystore file")
	ksSourceFile := flag.String("source", "", "keystore source file that used to generate keystore file")
	secret := flag.String("secret", "", "secret key to encrypt/decrypt file contents")
	seedKey := flag.String("seedkey", "", "seed key for keyservcie")
	keyId := flag.String("key", "", "key id used to encrypt/decrypt")
	encryptText := flag.String("encrypt", "", "text needs to be encrypted")
	decryptText := flag.String("decrypt", "", "text needs to be decrypted")
	flag.Parse()

	if *ksFile == "" {
		error("keystore required")
	}

	if *secret == "" {
		error("secret required")
	}

	if *ksSourceFile != "" {
		if contents, err := ioutil.ReadFile(*ksSourceFile); err != nil {
			error("read source file error:%v", err)
		} else {
			contents, err = keyservice.AesEncrypt(contents, []byte(*secret))
			if err != nil {
				error("encrypt source file contents error:%v", err)
			}
			err = ioutil.WriteFile(*ksFile, []byte(base64.RawURLEncoding.EncodeToString(contents)), 0600)
			if err != nil {
				error("write keystore file error:%v", err)
			}
		}
	}

	storage, err := keyservice.NewFileStorage(*ksFile, *secret)
	if err != nil {
		error("NewFileStorage error:%v", err)
	}

	if *seedKey == "" {
		error("no seed key")
	}

	service := keyservice.NewKeyService(*seedKey, storage, keyservice.NoCache)

	if *keyId == "" {
		error("no key id specified")
	}

	if *encryptText != "" {
		value, err := service.Encrypt(*encryptText, *keyId)
		if err != nil {
			error("encrypt error:%v", err)
		}
		fmt.Printf("encrypt result: %s\n", value)
	}

	if *decryptText != "" {
		value, err := service.Decrypt(*decryptText, *keyId)
		if err != nil {
			error("decrypt error:%v", err)
		}
		fmt.Printf("decrypt result: %s\n", value)
	}
}

func error(tpl string, args ...interface{}) {
	fmt.Printf(tpl+"\n", args...)
	flag.Usage()
	os.Exit(1)
}
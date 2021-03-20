package keyservice

import (
	"regexp"
	"strings"
)

// 字符串中加密串匹配正则
var embeddedEncryptedDataRegexp = regexp.MustCompile(`\$\$([\w-]+)\$\$`)

// 将内嵌在字符串中的加密数据，替换成解密后的数据
// 内嵌格式(由$$包裹加密数据)： ...$$加密数据$$... => ...解密数据...
// 一般用于解密JSON配置中的密文
func (sv *KeyService) DecryptEmbeddedString(content string, keyID string) (ret string, err error) {
	matches := embeddedEncryptedDataRegexp.FindAllStringSubmatch(content, -1)

	if matches == nil {
		return content, nil
	}

	var encryptedText, text string
	ret = content
	for _, match := range matches {
		encryptedText = match[1]
		if err != nil {
			return
		}
		text, err = sv.Decrypt(encryptedText, keyID)
		if err != nil {
			return
		}
		ret = strings.ReplaceAll(ret, match[0], text)
	}

	return
}

// config plugin
//func ProcessConfigEmbeddedCipher(content []byte) []byte {
//ret, err := manager.DecryptEmbeddedString(
//string(content),
//defaultConfigKey,
//)

//if err != nil {
//logger.Errorw("parse config cipher", "err", err, "alarm", true)
//}

//return []byte(ret)
//}

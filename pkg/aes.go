package tripartite

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"time"
)

//加密字符串
func GcmEncrypt(plaintext string) (string, error) {
	if len(AUTHORIZATION) != 32 && len(AUTHORIZATION) != 24 && len(AUTHORIZATION) != 16 {
		return "", errors.New("the length of key is error")
	}

	if len(plaintext) < 1 {
		return "", errors.New("plaintext is null")
	}

	keyByte := []byte(AUTHORIZATION)
	plainByte := []byte(plaintext)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	seal := aesGcm.Seal(nonce, nonce, plainByte, nil)
	return base64.URLEncoding.EncodeToString(seal), nil
}

//解密字符串
func GcmDecrypt(cipherText string) (string, error) {
	if len(AUTHORIZATION) != 32 && len(AUTHORIZATION) != 24 && len(AUTHORIZATION) != 16 {
		return "", errors.New("the length of key is error")
	}

	if len(cipherText) < 1 {
		return "", errors.New("cipherText is null")
	}

	cipherByte, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(cipherByte) < 12 {
		return "", errors.New("cipherByte is error")
	}

	nonce, cipherByte := cipherByte[:12], cipherByte[12:]

	keyByte := []byte(AUTHORIZATION)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainByte, err := aesGcm.Open(nil, nonce, cipherByte, nil)
	if err != nil {
		return "", err
	}

	return string(plainByte), nil
}

//生成32位md5字串
func GetAesKey(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

func JSONMarshal(obj interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	// json.NewEncoder.Encode adds a final '\n', json.Marshal does not.
	// Let's keep the default json.Marshal behaviour.
	res := b.Bytes()
	if len(res) >= 1 && res[len(res)-1] == '\n' {
		res = res[:len(res)-1]
	}
	return res, nil
}

func getToken(user string) string {
	effectsTime := user + "|" + strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10) //effects time after
	token, _ := GcmEncrypt(effectsTime)
	return token
}

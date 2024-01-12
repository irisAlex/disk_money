package aes

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

const AUTHORIZATION = "9505C098192D4BCECE6C22F77E63BFB2"

// 加密字符串
func GcmEncrypt(plaintext string) ([]byte, error) {
	if len(AUTHORIZATION) != 32 && len(AUTHORIZATION) != 24 && len(AUTHORIZATION) != 16 {
		return []byte{}, errors.New("the length of key is error")
	}

	if len(plaintext) < 1 {
		return []byte{}, errors.New("plaintext is null")
	}

	keyByte := []byte(AUTHORIZATION)
	plainByte := []byte(plaintext)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return []byte{}, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return aesGcm.Seal(nonce, nonce, plainByte, nil), nil
	// return base64.URLEncoding.EncodeToString(seal), nil
}

// 解密字符串
func GcmDecrypt(cipherByte []byte) (string, error) {
	if len(AUTHORIZATION) != 32 && len(AUTHORIZATION) != 24 && len(AUTHORIZATION) != 16 {
		return "", errors.New("the length of key is error")
	}

	if len(cipherByte) < 1 {
		return "", errors.New("cipherText is null")
	}

	// cipherByte, err := base64.URLEncoding.DecodeString(cipherText)
	// if err != nil {
	// 	return "", err
	// }
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

// 生成32位md5字串
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

func GetToken(user string) string {
	effectsTime := user + "|" + strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10) //effects time after
	token, _ := GcmEncrypt(effectsTime)
	return base64.URLEncoding.EncodeToString(token)
}

func Md5(key *string) {
	hash := md5.New()

	// 将字符串转换为字节数组并写入哈希对象
	hash.Write([]byte(*key))

	// 计算哈希值
	hashInBytes := hash.Sum(nil)

	*key = hex.EncodeToString(hashInBytes)
}

func VerifyToken(token string) (string, error) {
	cipherByte, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("Token 不存在")
	}

	verfiyInfo, err := GcmDecrypt(cipherByte)
	if err != nil {
		return "", err
	}

	return verfiyInfo, nil
}

func EncryptCardKey(ty string) {

}

func DecryptCardKey(cardkey string) (int, error) {
	ck, err := GcmDecrypt([]byte(cardkey))
	if err != nil {
		return 0, errors.New("解密失败")
	}

	setMeal, _ := strconv.Atoi(string(ck[0]))
	return setMeal, nil

}

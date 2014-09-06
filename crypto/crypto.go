package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/dchest/pbkdf2"
)

type Encryptor struct {
	Key      []byte
	Salt     []byte
	Password []byte
}

func New(password, saltB64 string) Encryptor {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		panic(err)
	}

	key := pbkdf2.WithHMAC(sha256.New, []byte(password), salt, 1000, 32)

	return Encryptor{
		Key:      key,
		Salt:     salt,
		Password: []byte(password),
	}
}

func (e *Encryptor) EncryptAndEncode(data string) string {
	encData, err := e.Encrypt([]byte(data))
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(encData)
}

func (e *Encryptor) DecodeAndDecrypt(data string) string {
	encData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}

	decData, err := e.Decrypt(encData)
	if err != nil {
		panic(err)
	}

	return string(decData)
}

func (e *Encryptor) Encrypt(src []byte) ([]byte, error) {
	dst := make([]byte, len(src))
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)

	aesBlockEncrypter, err := aes.NewCipher([]byte(e.Key))
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)

	dst = append(dst, iv...)

	return dst, nil
}

func (e *Encryptor) Decrypt(src []byte) ([]byte, error) {
	dst := make([]byte, len(src)-aes.BlockSize)
	iv := src[len(src)-aes.BlockSize:]
	src = src[0 : len(src)-aes.BlockSize]

	aesBlockDecrypter, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)

	return dst, nil
}

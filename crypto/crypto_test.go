package crypto_test

import (
	"crypto/rand"
	"encoding/base64"
	. "totp/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crypto", func() {
	var encryptor Encryptor
	password := "iAmPassword"
	salt := make([]byte, 32)

	BeforeEach(func() {
		rand.Read(salt)
		saltB64 := base64.StdEncoding.EncodeToString(salt)
		encryptor = New(password, saltB64)
	})

	Describe("New(password)", func() {
		It("creates a key", func() {
			Expect(len(encryptor.Key)).To(Equal(32))
		})

		It("Sets the salt", func() {
			Expect(len(encryptor.Salt)).To(Equal(32))
		})
	})

	Describe("Encrypt And Decrypt", func() {
		rawData := "Hello World I am A String Much Much String"

		It("encrypts and decrypts data", func() {
			encData, err := encryptor.Encrypt([]byte(rawData))
			Expect(err).To(BeNil())

			decData, err := encryptor.Decrypt(encData)
			Expect(err).To(BeNil())

			Expect(string(decData)).To(Equal(rawData))
		})
	})

	Describe("Encrypt and Encode, and Decode and Decrypt", func() {
		rawData := "Hello World I am A String Much Much String"

		It("encrypts and decrypts data", func() {
			encDataB64 := encryptor.EncryptAndEncode(rawData)
			decData := encryptor.DecodeAndDecrypt(encDataB64)

			Expect(decData).To(Equal(rawData))
		})
	})
})

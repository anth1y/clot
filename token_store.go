package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bradylove/clot/crypto"
	"github.com/bradylove/clot/totp"
	"io/ioutil"
	"os"
)

type TokenStore struct {
	TokenStoreData TokenStoreData
	Filepath       string
	Encryptor      crypto.Encryptor
	Password       string
}

type TokenStoreData struct {
	SaltB64  string
	V1Tokens []*totp.Token
}

func NewTokenStore(path, password string) (store TokenStore) {
	store = TokenStore{
		Filepath: path,
		Password: password,
	}

	if store.Exists() {
		err := store.Load()
		if err != nil {
			panic(err)
		}

		for _, t := range store.TokenStoreData.V1Tokens {
			t.Secret = store.Encryptor.DecodeAndDecrypt(t.EncryptedSecret)
		}
	} else {
		salt := make([]byte, 32)
		rand.Read(salt)

		saltb64 := base64.StdEncoding.EncodeToString(salt)

		store.TokenStoreData = TokenStoreData{
			SaltB64:  saltb64,
			V1Tokens: []*totp.Token{},
		}

		store.Encryptor = crypto.New(store.Password, store.TokenStoreData.SaltB64)

		store.Save()
	}

	return store
}

func (s *TokenStore) Exists() bool {
	if _, err := os.Stat(s.Filepath); err == nil {
		return true
	}

	return false
}

func (s *TokenStore) Save() error {
	var buff bytes.Buffer

	err := json.NewEncoder(&buff).Encode(s.TokenStoreData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Filepath, buff.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *TokenStore) Load() error {
	data, err := ioutil.ReadFile(s.Filepath)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(data)
	var tsd TokenStoreData
	err = json.NewDecoder(buff).Decode(&tsd)
	if err != nil {
		return err
	}

	s.TokenStoreData = tsd

	s.Encryptor = crypto.New(s.Password, s.TokenStoreData.SaltB64)

	fmt.Println("Before decrypting secrets!")

	return nil
}

func (s *TokenStore) AddToken(t totp.Token) {
	t.EncryptedSecret = s.Encryptor.EncryptAndEncode(t.Secret)

	s.TokenStoreData.V1Tokens = append(s.TokenStoreData.V1Tokens, &t)
	s.Save()
}

func (s *TokenStore) TokenCount() int {
	return len(s.TokenStoreData.V1Tokens)
}

func (s *TokenStore) Tokens() []*totp.Token {
	return s.TokenStoreData.V1Tokens
}

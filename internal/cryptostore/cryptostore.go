package cryptostore

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	senc "github.com/jbenet/go-simple-encrypt"
)

type CryproStore[T any] interface {
	Store(record T, key string) error
	Load(key string) (*T, error)
}

func NewCryptoStore[T any](fileName string) CryproStore[T] {
	return &CryptoStoreImpl[T]{
		fileName: fileName,
	}
}

type CryptoStoreImpl[T any] struct {
	fileName string
}

func (c *CryptoStoreImpl[T]) Store(record T, key string) error {

	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return err
	}

	cipherReader, err := senc.Encrypt(keyToHash(key), bytes.NewReader(jsonRecord))
	if err != nil {
		return err
	}

	ciphertext, err := ioutil.ReadAll(cipherReader)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(c.fileName), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.fileName, []byte(ciphertext), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *CryptoStoreImpl[T]) Load(key string) (*T, error) {
	encryptedStr, err := ioutil.ReadFile(c.fileName)
	if err != nil {
		return nil, err
	}

	decReader, err := senc.Decrypt(keyToHash(key), bytes.NewReader(encryptedStr))
	if err != nil {
		return nil, err
	}

	dec, err := ioutil.ReadAll(decReader)
	if err != nil {
		return nil, err
	}

	var record T

	err = json.Unmarshal(dec, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func keyToHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

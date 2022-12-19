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

type CryptoStoreRecord struct {
	APIKey     string
	EmployeeID string
}

type CryproStore interface {
	Store(record CryptoStoreRecord, key string) error
	Load(key string) (CryptoStoreRecord, error)
}

func NewCryptoStore(fileName string) CryproStore {
	return &CryptoStoreImpl{
		fileName: fileName,
	}
}

type CryptoStoreImpl struct {
	fileName string
}

func (c *CryptoStoreImpl) Store(record CryptoStoreRecord, key string) error {
	// 1. serialize record

	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// 2. encrypt serialized record

	cipherReader, err := senc.Encrypt(keyToHash(key), bytes.NewReader(jsonRecord))
	if err != nil {
		return err
	}

	ciphertext, err := ioutil.ReadAll(cipherReader)
	if err != nil {
		return err
	}

	// 3. store encrypted record

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

func (c *CryptoStoreImpl) Load(key string) (CryptoStoreRecord, error) {
	// 1. read encrypted record

	encryptedStr, err := ioutil.ReadFile(c.fileName)
	if err != nil {
		return CryptoStoreRecord{}, err
	}

	// 2. decrypt encrypted record

	decReader, err := senc.Decrypt(keyToHash(key), bytes.NewReader(encryptedStr))
	if err != nil {
		return CryptoStoreRecord{}, err
	}

	dec, err := ioutil.ReadAll(decReader)
	if err != nil {
		return CryptoStoreRecord{}, err
	}

	// 3. deserialize decrypted record

	var record CryptoStoreRecord

	err = json.Unmarshal(dec, &record)
	if err != nil {
		return CryptoStoreRecord{}, err
	}

	return record, nil
}

func keyToHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

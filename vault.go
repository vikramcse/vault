package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/vikramcse/vault/encrypt"
)

type Vault struct {
	encodingKey string
	filepath    string
	keyValues   map[string]string
	mutex       sync.Mutex
}

// NewFileVault will crate a new vault value
func NewFileVault(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
		keyValues:   make(map[string]string),
	}
}

func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	defer f.Close()

	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}

	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}

	decryptedJson, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	r := strings.NewReader(decryptedJson)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)

	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}

	encryptedJson, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, encryptedJson)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}

	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("Secter: no value for the provided  key")
	}

	return value, nil
}

func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return err
	}

	v.keyValues[key] = value
	err = v.saveKeyValues()
	if err != nil {
		return err
	}
	return nil
}

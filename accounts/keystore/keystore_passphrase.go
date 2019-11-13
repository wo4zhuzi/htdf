package keystore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/crypto/keys/mintkey"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

type keyStorePassphrase struct {
	keysDirPath string
}

func (ks keyStorePassphrase) GetKey(addr, filename, auth string) (*Key, error) {
	keyJSON := new(Key)
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bz, &keyJSON)
	if keyJSON.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", keyJSON.Address, addr)
	}
	return keyJSON, nil
}

// StoreKey generates a key, encrypts with passphrase and stores in the given directory
func StoreKey(dir, passphrase string) (string, string, error) {
	_, secret, acc, err := storeNewKey(&keyStorePassphrase{dir}, passphrase)
	return acc.Address, secret, err
}

//
func StoreKeyEx(key *Key) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(keyFileName(key.Address), content)
}

//StoreKey store privkey with file name
func (ks keyStorePassphrase) StoreKey(filename string, k *Key) error {
	content, err := json.Marshal(k)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}

func writeKeyFile(file string, content []byte) error {
	name, err := writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}

	return os.Rename(name, file)
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700 //0700-junying-todo-20190422

	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	//f, err := ioutil.TempFile(defaultKeyStoreHome, file)
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

//GetPrivKey return a privkey
func GetPrivKey(acc accounts.Account, passphrase string, rootDir string) (crypto.PrivKey, error) {
	if rootDir == "" {
		rootDir = viper.GetString(cli.HomeFlag)
	}
	defaultKeyStoreHome := filepath.Join(rootDir, "keystores")
	ks := NewKeyStore(defaultKeyStoreHome)
	acc, err := ks.Find(acc)
	_, key, err := ks.getDecryptedKey(acc, passphrase)
	if err != nil {
		return nil, err
	}

	if key.Address != acc.Address {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", key.Address, acc.Address)
	}

	privKey, err := mintkey.UnarmorDecryptPrivKey(key.PrivKeyArmor, passphrase)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

//
func GetPrivKeyEx(accaddr sdk.AccAddress, passphrase string, rootDir string) (crypto.PrivKey, error) {
	bech32 := sdk.AccAddress.String(accaddr)
	account := accounts.Account{Address: bech32}
	return GetPrivKey(account, passphrase, rootDir)
}

func (ks keyStorePassphrase) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}

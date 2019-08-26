package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"

	tmcrypto "github.com/orientwalt/tendermint/crypto"
	"github.com/orientwalt/tendermint/types"
	"github.com/tendermint/go-amino"

	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/orientwalt/htdf/client/keys"
)

// junying-todo-20190429
// Unlock unlocks an account when address and password are given.
func Unlock(encrypted, passphrase string) (tmcrypto.PrivKey, error) {
	account := accounts.Account{Address: encrypted}
	privkey, err := keystore.GetPrivKey(account, passphrase, "")
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

// junying-todo-20190429
// UnlockByStdIn needs user to type password when bechaddr is given.
func UnlockByStdIn(bech32 string) (tmcrypto.PrivKey, error) {
	passphrase, err := keys.ReadPassphraseFromStdin(bech32)
	if err != nil {
		return nil, err
	}
	return Unlock(bech32, passphrase)
}

// LoadGenesisDoc reads and unmarshals GenesisDoc from the given file.
func LoadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc types.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

//
func WriteString(filepath string, msg string) error {
	err := ioutil.WriteFile(filepath, []byte(msg), 0644)
	return err
}

//
func ReadString(filepath string, lineNum int) (line string, lastLine int, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return line, lastLine, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

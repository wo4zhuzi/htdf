package server

import (
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/crypto/keys"

	clkeys "github.com/orientwalt/htdf/client/keys"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/accounts/keystore"
)

// GenerateCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateCoinKey() (sdk.AccAddress, string, error) {

	// generate a private key, with recovery phrase
	info, secret, err := clkeys.NewInMemoryKeyBase().CreateMnemonic(
		"name", keys.English, "pass", keys.Secp256k1)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	addr := info.GetPubKey().Address()
	return sdk.AccAddress(addr), secret, nil
}

// GenerateSaveCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateSaveCoinKey(clientRoot, keyName, keyPass string,
	overwrite bool) (sdk.AccAddress, string, error) {

	// get the keystore from the client
	keybase, err := clkeys.NewKeyBaseFromDir(clientRoot)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}

	// ensure no overwrite
	if !overwrite {
		_, err := keybase.Get(keyName)
		if err == nil {
			return sdk.AccAddress([]byte{}), "", fmt.Errorf(
				"key already exists, overwrite is disabled (clientRoot: %s)", clientRoot)
		}
	}

	// generate a private key, with recovery phrase
	info, secret, err := keybase.CreateMnemonic(keyName, keys.English, keyPass, keys.Secp256k1)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}

	return sdk.AccAddress(info.GetPubKey().Address()), secret, nil
}

// junying-todo-20190420
// GenerateSaveCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateSaveCoinKeyEx(clientRoot, keyPass string) (sdk.AccAddress, string, error) {
	defaultKeyStoreHome := filepath.Join(clientRoot, "keystores")
	// generate a private key, with recovery phrase
	bech32addr, secret, err := keystore.StoreKey(defaultKeyStoreHome, keyPass)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	accaddr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}
	return accaddr, secret, nil
}

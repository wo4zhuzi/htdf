package keystore

import (
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/accounts"
)

type keystoreWallet struct {
	account  accounts.Account
	keystore *KeyStore
}

// URL implements accounts.Wallet, returning the URL of the account within.
func (w *keystoreWallet) URL() accounts.URL {
	return w.account.URL
}

// Open implements accounts.Wallet, but is a noop for plain wallets since there
// is no connection or decryption step necessary to access the list of accounts.
func (w *keystoreWallet) Open(passphrase string) error { return nil }

// Close implements accounts.Wallet, but is a noop for plain wallets since is no
// meaningful open operation.
func (w *keystoreWallet) Close() error { return nil }

// Accounts implements accounts.Wallet, returning an account list consisting of
// a single account that the plain kestore wallet contains.
func (w *keystoreWallet) Accounts() []accounts.Account {
	return []accounts.Account{w.account}
}

// Contains implements accounts.Wallet, returning whether a particular account is
// or is not wrapped by this wallet instance.
func (w *keystoreWallet) Contains(account accounts.Account) bool {
	return account.Address == w.account.Address && (account.URL == (accounts.URL{}) || account.URL == w.account.URL)
}

// Derive implements accounts.Wallet, but is a noop for plain wallets since there
// is no notion of hierarchical account derivation for plain keystore accounts.
func (w *keystoreWallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, accounts.ErrNotSupported
}

// SignTxWithPassphrase implements accounts.Wallet, attempting to sign the given
// transaction with the given account using passphrase as extra authentication.
func (w *keystoreWallet) SignTx(a accounts.Account, passphrase string, txbuilder authtxb.TxBuilder, stdTx auth.StdTx) (auth.StdTx, error) {
	// Make sure the requested account is contained within
	if a.Address != w.account.Address {
		return stdTx, accounts.ErrUnknownAccount
	}
	if a.URL != (accounts.URL{}) && a.URL != w.account.URL {
		return stdTx, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTx(a, passphrase, txbuilder, stdTx)
}

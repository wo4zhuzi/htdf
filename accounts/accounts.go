package accounts

import (
	"github.com/orientwalt/htdf/x/auth"

	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
)

type Account struct {
	Address string `json:"address"`
	URL     URL
}

type KeyStoreWallets interface {
	Accounts() ([]Account, error)

	BuildAndSign(txbuilder authtxb.TxBuilder, addr string, passphrase string, msgs []sdk.Msg) ([]byte, error)

	Sign(txbuilder authtxb.TxBuilder, addr string, passphrase string, msg authtxb.StdSignMsg) ([]byte, error)

	SignStdTx(txbuilder authtxb.TxBuilder, stdTx auth.StdTx, addr string, passphrase string) (signedStdTx auth.StdTx, err error)

	FindPrivKey(addr string) (string, error)
}

type KeyStores interface {
	NewKey(passphrase string) (string, error)

	RecoverKey(strPrivKey string, passPhrase string) error

	RecoverKeyByMnemonic(mnemonic string, bip39Passphrase string, passPhrase string, account, index uint32) error
}

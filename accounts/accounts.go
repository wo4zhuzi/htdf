package accounts

import (
	"github.com/orientwalt/htdf/accounts/event"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
)

type Account struct {
	Address string `json:"address"`
	URL     URL
}

type Wallet interface {
	URL() URL

	Open(passphrase string) error

	Close() error

	Accounts() []Account

	Contains(account Account) bool

	Derive(path DerivationPath, pin bool) (Account, error)

	SignTx(account Account, passphrase string, txbuilder authtxb.TxBuilder, stdTx auth.StdTx) (signedStdTx auth.StdTx, err error)
}

type Backend interface {
	Wallets() []Wallet

	Subscribe(sink chan<- WalletEvent) event.Subscription
}

//WalletEventType represents the different event types that can be fired by
//the wallet subscription subsystem.
type WalletEventType int

const (
	// WalletArrived is fired when a new wallet is detected either via USB or via
	// a filesystem event in the keystore.
	WalletArrived WalletEventType = iota

	// WalletOpened is fired when a wallet is successfully opened with the purpose
	// of starting any background processes such as automatic key derivation.
	WalletOpened

	// WalletDropped
	WalletDropped
)

// WalletEvent is an event fired by an account backend when a wallet arrival or
// departure is detected.
type WalletEvent struct {
	Wallet Wallet          // Wallet instance arrived or departed
	Kind   WalletEventType // Event type that happened in the system
}

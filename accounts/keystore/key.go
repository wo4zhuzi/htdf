package keystore

import (
	"fmt"
	"time"

	"github.com/cosmos/go-bip39"
	"github.com/orientwalt/htdf/crypto/keys"
	"github.com/orientwalt/htdf/crypto/keys/hd"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/crypto/keys/mintkey"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/spf13/viper"
)

const (
	flagAccount string = "account"
	flagIndex   string = "index"
	//FlagPublicKey       string = "pubkey"
	mnemonicEntropySize int = 256
)

type keyStore interface {
	// Loads and decrypts the key from disk.
	GetKey(addr, filename, auth string) (*Key, error)
	// Writes and encrypts the key.
	StoreKey(filename string, k *Key) error
	// Joins filename with the key directory unless it is already absolute.
	JoinPath(filename string) string
}

// Key is the public information about a locally stored key
type Key struct {
	Address      string `json:"address"`
	PubKey       string `json:"pubkey"`
	PrivKeyArmor string `json:"privkey.armor"`
}

func newKey(pub crypto.PubKey, privArmor string) *Key {
	accAddr := sdk.AccAddress(pub.Address().Bytes())
	pubKey, err := sdk.Bech32ifyAccPub(pub)
	if err != nil {
		fmt.Println("newKey error for Bech32ifyAccPub !", err)
		return nil
	}
	return &Key{
		Address:      accAddr.String(),
		PubKey:       pubKey,
		PrivKeyArmor: privArmor,
	}
}

func storeNewKey(ks keyStore, passphrase string) (*Key, string, accounts.Account, error) {
	key, secret, err := generateKey(passphrase)
	if err != nil {
		return nil, "", accounts.Account{}, err
	}

	acc := accounts.Account{Address: key.Address, URL: accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))}}
	if err := ks.StoreKey(acc.URL.Path, key); err != nil {
		return nil, "", acc, err
	}
	return key, secret, acc, err
}

func setBytes(a *[32]byte, b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-32:]
	}
	copy(a[32-len(b):], b)
}

func recoverOldKey(ks keyStore, privateKey []byte, passphrase string) (*Key, accounts.Account, error) {
	var privateKey32 [32]byte
	setBytes(&privateKey32, privateKey)

	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("privateKey32=%v\n", privateKey32)

	key, err := importPrivateKey(privateKey32, passphrase)
	if err != nil {
		return nil, accounts.Account{}, err
	}

	acc := accounts.Account{Address: key.Address, URL: accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))}}
	if err := ks.StoreKey(acc.URL.Path, key); err != nil {
		return nil, acc, err
	}
	return key, acc, err
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address string>
func keyFileName(keyAddr string) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), keyAddr)
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

//
type Seed struct {
	mnemonic        string
	bip39Passphrase string
	fAccount        uint32
	index           uint32
}

func newSeed() (*Seed, error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	mnemc, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		return nil, err
	}
	seed := &Seed{
		mnemonic:        mnemc,
		bip39Passphrase: keys.DefaultBIP39Passphrase,
		fAccount:        uint32(viper.GetInt(flagAccount)),
		index:           uint32(viper.GetInt(flagIndex)),
	}

	return seed, nil
}

func generateKey(passphrase string) (*Key, string, error) {
	s, err := newSeed()
	if err != nil {
		return nil, "", err
	}
	hdPath := hd.NewFundraiserParams(s.fAccount, s.index)
	seed, err := bip39.NewSeedWithErrorChecking(s.mnemonic, s.bip39Passphrase)
	if err != nil {
		return nil, "", err
	}
	// create master key and derive first key:
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return nil, "", err
	}

	privkey := secp256k1.PrivKeySecp256k1(derivedPriv)

	privArmor := mintkey.EncryptArmorPrivKey(privkey, passphrase)

	pubkey := privkey.PubKey()

	key := newKey(pubkey, privArmor)

	return key, s.mnemonic, err
}

func importPrivateKey(inputPrivateKey [32]byte, passphrase string) (*Key, error) {

	privkey := secp256k1.PrivKeySecp256k1(inputPrivateKey)

	fmt.Printf("importPrivateKey|privkey=%v\n", privkey)

	privArmor := mintkey.EncryptArmorPrivKey(privkey, passphrase)

	pubkey := privkey.PubKey()

	key := newKey(pubkey, privArmor)

	return key, nil
}

//
func GenerateKeyEx(mnemonic, bip39Passphrase, passphrase string, account, index uint32) (*Key, error) {
	// create seed
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
	if err != nil {
		return nil, err
	}
	//
	hdPath := hd.NewFundraiserParams(account, index)
	// create master key and derive first key:
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	if err != nil {
		return nil, err
	}

	privkey := secp256k1.PrivKeySecp256k1(derivedPriv)

	privArmor := mintkey.EncryptArmorPrivKey(privkey, passphrase)

	pubkey := privkey.PubKey()

	key := newKey(pubkey, privArmor)

	return key, err
}

package keystore

import (
	"fmt"
	"testing"

	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
)

func TestKeyStore(t *testing.T) {
	password := "12345678"
	keydir := "keystores"

	ks := NewKeyStore(keydir)
	acc, err := ks.NewAccount(password)
	if err != nil {
		fmt.Print("New account error: ", err, "\n")
		return
	}
	fmt.Print("Account address:	", acc.Address)
	fmt.Print("Account URL:	", acc.URL)
}

func TestSignTx(t *testing.T) {

	password := "12345678"
	keydir := "keystores"

	accountNumber := uint64(0)
	sequence := uint64(1)
	gas := uint64(0)
	gasAjs := 1.0
	denom := "t_coin"
	memo := "TCoin"
	ID := "test_chain"

	cdc := codec.New()
	amount := sdk.NewInt(1)
	coin := sdk.NewCoin(denom, amount)
	c := []sdk.Coin{coin}
	coins := sdk.Coins(c)

	ks := NewKeyStore(keydir)
	acc, err := ks.NewAccount(password)
	if err != nil {
		fmt.Print("New account error: ", err, "\n")
		return
	}

	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc),
		accountNumber,
		sequence,
		gas,
		gasAjs,
		false,
		ID,
		memo,
		coins,
		nil,
	)

	var signedTx auth.StdTx
	signedTx, err = ks.SignTx(acc, password, txBldr, signedTx)
}

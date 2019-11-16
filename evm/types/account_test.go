package types

import (
	"fmt"
	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"

	"github.com/magiconair/properties/assert"
	"testing"
)

const (
	strAddr = "htdf1gh8yeqhxx7n29fx3fuaksqpdsau7gxc0rteptt"
)

func init() {

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(params.Bech32PrefixAccAddr, params.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(params.Bech32PrefixValAddr, params.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(params.Bech32PrefixConsAddr, params.Bech32PrefixConsPub)
	config.Seal()
}

func TestAccountUpdate(t *testing.T) {

	accAddr, err := sdk.AccAddressFromBech32(strAddr)
	assert.Equal(t, err, nil)

	baseAccount := auth.NewBaseAccountWithAddress(accAddr)
	coins := sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(12345678)), sdk.NewCoin("stake", sdk.NewInt(2345678))}
	baseAccount.Coins = coins

	fmt.Printf("baseAccount=%v\n", baseAccount)

	account := NewAccount(&baseAccount)

	fmt.Printf("account|coins=%v\n", account.GetCoins())
	account.SetBalance(sdk.NewInt(30000))

	fmt.Printf("account|coins=%v\n", account.GetCoins())

}

func TestAccountInsert(t *testing.T) {

	accAddr, err := sdk.AccAddressFromBech32(strAddr)
	assert.Equal(t, err, nil)

	baseAccount := auth.NewBaseAccountWithAddress(accAddr)
	account := NewAccount(&baseAccount)

	account.SetBalance(sdk.NewInt(3500350))
	fmt.Printf("account|coins=%v\n", account.GetCoins())
	assert.Equal(t, account.Balance().String() == "3500350", true)

}

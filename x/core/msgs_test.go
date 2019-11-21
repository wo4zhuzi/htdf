package htdfservice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
	bank "github.com/orientwalt/htdf/x/bank"
)

func TestMsgSendRoute(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = NewMsgSendDefault(addr1, addr2, coins)

	require.Equal(t, msg.Route(), "htdfservice")
	require.Equal(t, msg.Type(), "sendfrom")
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	atom123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	atom0 := sdk.NewCoins(sdk.NewInt64Coin("atom", 0))
	atom123eth123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 123))
	atom123eth0 := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 0)}

	var emptyAddr sdk.AccAddress

	cases := []struct {
		valid bool
		tx    MsgSend
	}{
		{true, NewMsgSendDefault(addr1, addr2, atom123)},       // valid send
		{true, NewMsgSendDefault(addr1, addr2, atom123eth123)}, // valid send with multiple coins
		{false, NewMsgSendDefault(addr1, addr2, atom0)},        // non positive coin
		{false, NewMsgSendDefault(addr1, addr2, atom123eth0)},  // non positive coin in multicoins
		{false, NewMsgSendDefault(emptyAddr, addr2, atom123)},  // empty from addr
		{false, NewMsgSendDefault(addr1, emptyAddr, atom123)},  // empty to addr
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = NewMsgSendDefault(addr1, addr2, coins)
	res := msg.GetSignBytes()

	expected := `{"type":"cosmos-sdk/MsgSend","value":{"amount":[{"amount":"10","denom":"atom"}],"from_address":"cosmos1d9h8qat57ljhcm","to_address":"cosmos1da6hgur4wsmpnjyg"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	var msg = NewMsgSendDefault(sdk.AccAddress([]byte("input1")), sdk.AccAddress{}, sdk.NewCoins())
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431]")
}

func TestInputValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	addr2 := sdk.AccAddress([]byte{7, 8})
	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	multiCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20))

	var emptyAddr sdk.AccAddress
	emptyCoins := sdk.NewCoins()
	emptyCoins2 := sdk.NewCoins(sdk.NewInt64Coin("eth", 0))
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txIn  bank.Input
	}{
		// auth works with different apps
		{true, bank.NewInput(addr1, someCoins)},
		{true, bank.NewInput(addr2, someCoins)},
		{true, bank.NewInput(addr2, multiCoins)},

		{false, bank.NewInput(emptyAddr, someCoins)},  // empty address
		{false, bank.NewInput(addr1, emptyCoins)},     // invalid coins
		{false, bank.NewInput(addr1, emptyCoins2)},    // invalid coins
		{false, bank.NewInput(addr1, someEmptyCoins)}, // invalid coins
		{false, bank.NewInput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txIn.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestOutputValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	addr2 := sdk.AccAddress([]byte{7, 8})
	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	multiCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20))

	var emptyAddr sdk.AccAddress
	emptyCoins := sdk.NewCoins()
	emptyCoins2 := sdk.NewCoins(sdk.NewInt64Coin("eth", 0))
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txOut bank.Output
	}{
		// auth works with different apps
		{true, bank.NewOutput(addr1, someCoins)},
		{true, bank.NewOutput(addr2, someCoins)},
		{true, bank.NewOutput(addr2, multiCoins)},

		{false, bank.NewOutput(emptyAddr, someCoins)},  // empty address
		{false, bank.NewOutput(addr1, emptyCoins)},     // invalid coins
		{false, bank.NewOutput(addr1, emptyCoins2)},    // invalid coins
		{false, bank.NewOutput(addr1, someEmptyCoins)}, // invalid coins
		{false, bank.NewOutput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txOut.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

/*
// what to do w/ this test?
func TestMsgSendSigners(t *testing.T) {
	signers := []sdk.AccAddress{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	inputs := make([]Input, len(signers))
	for i, signer := range signers {
		inputs[i] = bank.NewInput(signer, someCoins)
	}
	tx := NewMsgSendDefault(inputs, nil)

	require.Equal(t, signers, tx.Signers())
}
*/

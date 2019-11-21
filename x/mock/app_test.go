package mock

import (
	"testing"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	hscore "github.com/orientwalt/htdf/x/core"
	abci "github.com/orientwalt/tendermint/abci/types"
	"github.com/orientwalt/tendermint/crypto"
	"github.com/orientwalt/tendermint/crypto/ed25519"
	"github.com/stretchr/testify/require"
)

type (
	expectedBalance struct {
		addr  sdk.AccAddress
		coins sdk.Coins
	}

	appTestCase struct {
		expSimPass       bool
		expPass          bool
		msgs             []sdk.Msg
		accNums          []uint64
		accSeqs          []uint64
		privKeys         []crypto.PrivKey
		expectedBalances []expectedBalance
	}
)

var (
	priv1 = ed25519.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = ed25519.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	priv4 = ed25519.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

	coins     = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}
	halfCoins = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)}
	manyCoins = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), sdk.NewInt64Coin("barcoin", 1)}
	freeFee   = auth.NewStdFee(100000, 0)

	sendMsg1 = hscore.NewMsgSendDefault(addr1, addr2, coins)
	sendMsg2 = hscore.NewMsgSendDefault(addr1, addr2, manyCoins)
)

// initialize the mock application for this module
func getMockApp(t *testing.T) *App {
	mapp, err := getBenchmarkMockApp()
	require.NoError(t, err)
	return mapp
}

func TestMsgSendWithAccounts(t *testing.T) {
	mapp := getMockApp(t)
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 67)},
	}

	SetGenesis(mapp, []auth.Account{acc})

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 57)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{sendMsg1, sendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: false,
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}

	// bumping the tx nonce number without resigning should be an auth error
	mapp.BeginBlock(abci.RequestBeginBlock{})

	tx := GenTx([]sdk.Msg{sendMsg1}, []uint64{0}, []uint64{0}, priv1)
	// tx.Signatures[0].Sequence = 1

	res := mapp.Deliver(tx)
	require.Equal(t, sdk.CodeUnauthorized, res.Code, res.Log)
	require.Equal(t, sdk.CodespaceRoot, res.Codespace)

	// resigning the tx with the bumped sequence should work
	SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{sendMsg1, sendMsg2}, []uint64{0}, []uint64{1}, true, true, priv1)
}

func TestMsgSendMultipleOut(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}

	SetGenesis(mapp, []auth.Account{acc1, acc2})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 47)}},
				{addr3, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestSengMsgMultipleInOut(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}
	acc4 := &auth.BaseAccount{
		Address: addr4,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}

	SetGenesis(mapp, []auth.Account{acc1, acc2, acc4})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg1},
			accNums:    []uint64{0, 0},
			accSeqs:    []uint64{0, 0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1, priv4},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 32)}},
				{addr4, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 52)}},
				{addr3, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgSendDependent(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)},
	}

	SetGenesis(mapp, []auth.Account{acc1})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{sendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv2},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

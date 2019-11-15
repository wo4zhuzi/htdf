// nolint
package auth

import (
	abci "github.com/orientwalt/tendermint/abci/types"
	"github.com/orientwalt/tendermint/crypto"
	"github.com/orientwalt/tendermint/crypto/secp256k1"
	dbm "github.com/orientwalt/tendermint/libs/db"
	"github.com/orientwalt/tendermint/libs/log"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  AccountKeeper
	fck FeeCollectionKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	RegisterBaseAccount(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := NewAccountKeeper(cdc, authCapKey, pk.Subspace(DefaultParamspace), ProtoBaseAccount)
	fck := NewFeeCollectionKeeper(cdc, fckCapKey)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, DefaultParams())

	return testInput{cdc: cdc, ctx: ctx, ak: ak, fck: fck}
}

func newTestMsg(addrs ...sdk.AccAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func newStdFee() sdk.StdFee {
	return sdk.NewStdFee(50000, "250atom")
}

// coins to more than cover the fee
func newCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("atom", 10000000),
	}
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// junying-todo, 2019-11-14
// multi-msg to uni-msg
func newTestTx(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64) sdk.Tx {
	var sig StdSignature
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, msg, "")

	sign, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sig = StdSignature{PubKey: priv.PubKey(), Signature: sign}

	tx := NewStdTx(msg, sig, "")
	return tx
}

// junying-todo, 2019-11-14
// multi-msg to uni-msg
func newTestTxWithMemo(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, memo string) sdk.Tx {
	var sig StdSignature
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, msg, memo)

	sign, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sig = StdSignature{PubKey: priv.PubKey(), Signature: sign}

	tx := NewStdTx(msg, sig, memo)
	return tx
}

func newTestTxWithSignBytes(msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, signBytes []byte, memo string) sdk.Tx {
	var sig StdSignature

	sign, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sig = StdSignature{PubKey: priv.PubKey(), Signature: sign}

	tx := NewStdTx(msg, sig, memo)
	return tx
}

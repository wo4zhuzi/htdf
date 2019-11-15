package auth

import (
	"fmt"
	"testing"

	abci "github.com/orientwalt/tendermint/abci/types"
	"github.com/orientwalt/tendermint/crypto/ed25519"
	"github.com/orientwalt/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
)

var (
	priv = ed25519.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

// func TestStdTx(t *testing.T) {
// 	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
// 	fee := newStdFee()
// 	sigs := []StdSignature{}

// 	tx := NewStdTx(msgs, fee, sigs, "")
// 	require.Equal(t, msgs, tx.GetMsgs())
// 	require.Equal(t, sigs, tx.GetSignatures())

// 	feePayer := tx.GetSigners()[0]
// 	require.Equal(t, addr, feePayer)
// }

func TestStdSignBytes(t *testing.T) {
	type args struct {
		chainID  string
		accnum   uint64
		sequence uint64
		msg      sdk.Msg
		memo     string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args{"1234", 3, 6, sdk.Msg(sdk.NewTestMsg(addr)), "memo"},
			fmt.Sprintf("{\"account_number\":\"3\",\"chain_id\":\"1234\",\"fee\":{\"amount\":[{\"amount\":\"150\",\"denom\":\"atom\"}],\"gas\":\"50000\"},\"memo\":\"memo\",\"msgs\":[[\"%s\"]],\"sequence\":\"6\"}", addr),
		},
	}
	for i, tc := range tests {
		got := string(StdSignBytes(tc.args.chainID, tc.args.accnum, tc.args.sequence, tc.args.msg, tc.args.memo))
		require.Equal(t, tc.want, got, "Got unexpected result on test case i: %d", i)
	}
}

func TestTxValidateBasic(t *testing.T) {
	ctx := sdk.NewContext(nil, abci.Header{ChainID: "mychainid"}, false, log.NewNopLogger())

	// keys and addresses
	_, _, addr1 := keyPubAddr()
	priv2, _, addr2 := keyPubAddr()

	// msg and signatures
	msg := newTestMsg(addr1, addr2)

	// require to fail validation when signatures do not match expected signers

	accNums, seqs := []uint64{0, 1}, []uint64{0, 0}
	tx := newTestTx(ctx, msg, priv2, accNums[0], seqs[0])

	err := tx.ValidateBasic()
	require.Error(t, err)
	require.Equal(t, sdk.CodeUnauthorized, err.Result().Code)
	require.NoError(t, err)
}

func TestDefaultTxEncoder(t *testing.T) {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	encoder := DefaultTxEncoder(cdc)

	msg := sdk.Msg(sdk.NewTestMsg(addr))
	sig := StdSignature{
		PubKey:    nil,
		Signature: nil,
	}

	tx := NewStdTx(msg, sig, "")

	cdcBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)

	require.NoError(t, err)
	encoderBytes, err := encoder(tx)

	require.NoError(t, err)
	require.Equal(t, cdcBytes, encoderBytes)
}

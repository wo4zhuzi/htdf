/*
	modified by junying, 2019-11-13
	[StdTx]
	Msgs -> Msg
	Fee removed
	[Msg]
	Gas
	GasPrice
*/
package auth

import (
	"encoding/json"
	"fmt"

	"github.com/orientwalt/tendermint/crypto"
	"github.com/orientwalt/tendermint/crypto/multisig"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/params"
	"github.com/orientwalt/htdf/server/config"
	"github.com/orientwalt/htdf/types"
	sdk "github.com/orientwalt/htdf/types"
)

var (
	_ sdk.Tx = (*StdTx)(nil)

	maxGasWanted = uint64((1 << 63) - 1)
)

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msg       sdk.Msg      `json:"msg"`
	Signature StdSignature `json:"signature"`
	Memo      string       `json:"memo"`
}

func NewStdTx(msg sdk.Msg, sig StdSignature, memo string) StdTx {
	return StdTx{
		Msg:       msg,
		Signature: sig,
		Memo:      memo,
	}
}

// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsg() sdk.Msg { return tx.Msg }

// junying-todo, 2019-11-14
func (tx StdTx) GetFee() sdk.StdFee {
	return tx.Msg.GetFee()
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	fee := tx.GetFee()
	if fee.GasWanted > maxGasWanted {
		return sdk.ErrGasOverflow(fmt.Sprintf("invalid gas supplied; %d > %d", fee.GasWanted, maxGasWanted))
	}
	if fee.Amount.IsAnyNegative() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", fee.Amount))
	}

	// junying-todo, 2019-11-13
	// MinGasPrice Checking
	var gasprice = fee.GasPrice
	minGasPrices, err := types.ParseCoins(config.DefaultMinGasPrices)
	if err != nil {
		return sdk.ErrTxDecode("DefaultMinGasPrices decode error")
	}
	if !gasprice.IsAllGTE(minGasPrices) {
		return sdk.ErrInsufficientFee(fmt.Sprintf("gasprice must be greater than %s", config.DefaultMinGasPrices))
	}
	// junying-todo, 2019-11-13
	// Validate Msgs &
	// Check MinGas for staking txs
	var msg = tx.Msg
	if msg == nil {
		return sdk.ErrUnknownRequest("Tx.GetMsg() must return at least one message in list")
	}
	// Validate the Msg.
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	// Checking minimum gasprice condition for staking transactions
	if msg.Route() != "htdfservice" {
		if fee.GasWanted < params.TxStakingDefaultGas {
			return sdk.ErrInternal(fmt.Sprintf("staking tx gas must be greater than %d", params.TxStakingDefaultGas))
		}
	}

	// junying-todo-expected, 2019-11-14
	// sign exist checking
	// stdSig := tx.GetSignature()

	return nil
}

// countSubKeys counts the total number of keys for a multi-sig public key.
func countSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}

	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += countSubKeys(subkey)
	}

	return numKeys
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigner() sdk.AccAddress {
	return tx.GetMsg().GetSigner()
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignatures returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignature() StdSignature { return tx.Signature }

//__________________________________________________________

//__________________________________________________________

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64          `json:"account_number"`
	ChainID       string          `json:"chain_id"`
	Memo          string          `json:"memo"`
	Msg           json.RawMessage `json:"msg"`
	Sequence      uint64          `json:"sequence"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum uint64, sequence uint64, msg sdk.Msg, memo string) []byte {
	bz, err := msgCdc.MarshalJSON(StdSignDoc{
		AccountNumber: accnum,
		ChainID:       chainID,
		Memo:          memo,
		Msg:           json.RawMessage(msg.GetSignBytes()),
		Sequence:      sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// StdSignature represents a sig
type StdSignature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
}

// junying-todo, 2019-11-14
func (s *StdSignature) IsEmpty() bool {
	if s.Signature == nil {
		return true
	}
	return false
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = StdTx{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}
		// fmt.Println("DefaultTxDecoder:tx", tx)
		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}

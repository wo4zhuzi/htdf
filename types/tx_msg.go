package types

import (
	"encoding/json"

	"github.com/orientwalt/htdf/codec"
)

var msgCdc = codec.New()

// StdFee includes the amount of coins paid in fees and the maximum
// gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
type StdFee struct {
	Amount    Coins  `json:"amount"`
	GasWanted uint64 `json:"gaswanted"`
	GasPrice  Coins  `json:"gasprice"`
}

// junying-todo, 2019-11-07
// fee = gas * gasprice
func CalcFees(gaswanted uint64, gasprices Coins) Coins {
	Fees := make(Coins, len(gasprices))
	gaslimit := NewInt(int64(gaswanted))
	for i, gp := range gasprices {
		fee := gp.Amount.Mul(gaslimit)
		Fees[i] = NewCoin(gp.Denom, fee)
	}
	return Fees
}

// junying-todo, 2019-11-14
func ZeroFee() StdFee {
	zerocoins := Coins{NewInt64Coin("satoshi", 10000000)}
	return StdFee{
		Amount:    zerocoins,
		GasWanted: 0,
		GasPrice:  zerocoins,
	}
}

func NewStdFee(gaswanted uint64, gasprice string) StdFee {
	coins, err := ParseCoins(gasprice)
	if err != nil {
		return ZeroFee()
	}
	return StdFee{
		Amount:    CalcFees(gaswanted, coins),
		GasWanted: gaswanted,
		GasPrice:  coins,
	}
}

// Bytes for signing later
func (fee StdFee) Bytes() []byte {
	// normalize. XXX
	// this is a sign of something ugly
	// (in the lcd_test, client side its null,
	// server side its [])
	if len(fee.Amount) == 0 {
		fee.Amount = NewCoins()
	}
	bz, err := msgCdc.MarshalJSON(fee) // TODO
	if err != nil {
		panic(err)
	}
	return bz
}

// GasPrices returns the gas prices for a StdFee.
//
// NOTE: The gas prices returned are not the true gas prices that were
// originally part of the submitted transaction because the fee is computed
// as fee = ceil(gasWanted * gasPrices).
func (fee StdFee) GetGasPrice() uint64 {
	if len(fee.GasPrice) == 0 {
		return 0
	}
	return fee.GasPrice[0].Amount.Uint64()
}

// junying-todo, 2019-11-07
func (fee StdFee) GetAmount() Coins {
	return fee.Amount
}

// Transactions messages must fulfill the Msg
type Msg interface {

	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigner() AccAddress

	GetFee() StdFee
	SetFee(StdFee)
	// // added by junying,2019-11-13
	// SetGasWanted(uint64)
	// GetGasWanted() uint64

	// // added by junying,2019-11-13
	// SetGasPrice(string)
	// GetGasPrice() uint64
}

//__________________________________________________________

// Transactions objects must fulfill the Tx
type Tx interface {
	// Gets the all the transaction's messages.
	GetMsg() Msg

	// ValidateBasic does a simple and lightweight validation check that doesn't
	// require access to any other information.
	ValidateBasic() Error
}

//__________________________________________________________

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, Error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)

//__________________________________________________________

var _ Msg = (*TestMsg)(nil)

// msg type for testing
type TestMsg struct {
	addresses []AccAddress
	fee       StdFee
}

func NewTestMsg(addrs ...AccAddress) *TestMsg {
	return &TestMsg{
		addresses: addrs,
		fee:       ZeroFee(),
	}
}

//nolint
func (msg *TestMsg) Route() string { return "TestMsg" }

//
func (msg *TestMsg) Type() string { return "Test message" }

//
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.addresses)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bz)
}

//
func (msg *TestMsg) ValidateBasic() Error { return nil }

//
func (msg *TestMsg) GetSigner() AccAddress { return msg.addresses[0] }

//
func (msg *TestMsg) GetFee() StdFee { return msg.fee }

//
func (msg *TestMsg) SetFee(fee StdFee) { msg.fee = fee }

// //
// func (msg *TestMsg) GetGasWanted() uint64 { return msg.gaswanted }

// //
// func (msg *TestMsg) SetGasWanted(gaswanted uint64) { msg.gaswanted = gaswanted }

// //
// func (msg *TestMsg) GetGasPrice() uint64 {
// 	gasprice, err := ParseCoin(msg.gasprice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg *TestMsg) SetGasPrice(gasprice string) { msg.gasprice = gasprice }

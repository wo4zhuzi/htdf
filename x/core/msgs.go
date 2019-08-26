package htdfservice

import (
	"encoding/json"

	"github.com/orientwalt/htdf/types"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/orientwalt/htdf/types"
)

const (
	defaultGasLimit = 21000
	defaultGasPrice = 1
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgSendFrom defines a SendFrom message /////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgSendFrom struct {
	From     sdk.AccAddress
	To       sdk.AccAddress
	Amount   sdk.Coins
	Data     string
	Gas      uint64 //unit,  gallon
	GasPrice uint64 //unit,  satoshi/gallon
	GasLimit uint64 //unit,  gallon
}

var _ sdk.Msg = MsgSendFrom{}

// NewMsgSend is a constructor function for MsgSend
func NewMsgSendFrom(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins) MsgSendFrom {
	return MsgSendFrom{
		From:     fromaddr,
		To:       toaddr,
		Amount:   amount,
		Gas:      defaultGasLimit,
		GasPrice: defaultGasPrice,
		GasLimit: defaultGasLimit,
	}
}

func NewMsgSendFromForData(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins, data string, gas uint64, gasPrice uint64, gasLimit uint64) MsgSendFrom {
	return MsgSendFrom{
		From:     fromaddr,
		To:       toaddr,
		Amount:   amount,
		Data:     data,
		Gas:      gas,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}
}

// Route should return the name of the module
func (msg MsgSendFrom) Route() string { return "htdfservice" }

// Type should return the action
func (msg MsgSendFrom) Type() string { return "sendfrom" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSendFrom) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}

	if len(msg.Data) == 0 {
		// classic transfer

		// must have to address
		if msg.To.Empty() {
			return sdk.ErrInvalidAddress(msg.To.String())
		}

		// amount > 0
		if !msg.Amount.IsAllPositive() {
			return sdk.ErrInsufficientCoins("Amount must be positive")
		}
	} else {
		// access smart contract

		// amount must be 0
		if !msg.Amount.IsZero() {
			return sdk.NewError(types.Codespace, types.ErrCode_BeZeroAmount, "access smart contract, amount must be 0")
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendFrom) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSendFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetStringAddr defines whose fromaddr is required
func (msg MsgSendFrom) GetFromAddrStr() string {
	return sdk.AccAddress.String(msg.From)
}

//
func (msg MsgSendFrom) GetFromAddr() sdk.AccAddress {
	return msg.From
}

//
func (msg MsgSendFrom) FromAddress() common.Address {
	return types.ToEthAddress(msg.From)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgAdd defines a Add message ///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgAdd struct {
	SystemIssuer sdk.AccAddress
	Amount       sdk.Coins
}

var _ sdk.Msg = MsgAdd{}

// NewMsgAdd is a constructor function for Msgadd
func NewMsgAdd(addr sdk.AccAddress, amount sdk.Coins) MsgAdd {
	return MsgAdd{
		SystemIssuer: addr,
		Amount:       amount,
	}
}

// Route should return the name of the module
func (msg MsgAdd) Route() string { return "htdfservice" }

// Type should return the action
func (msg MsgAdd) Type() string { return "add" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAdd) ValidateBasic() sdk.Error {
	if msg.SystemIssuer.Empty() {
		return sdk.ErrInvalidAddress(msg.SystemIssuer.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAdd) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgAdd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SystemIssuer}
}

// GetStringAddr defines whose fromaddr is required
func (msg MsgAdd) GetSystemIssuerStr() string {
	return sdk.AccAddress.String(msg.SystemIssuer)
}

//
func (msg MsgAdd) GeSystemIssuer() sdk.AccAddress {
	return msg.SystemIssuer
}

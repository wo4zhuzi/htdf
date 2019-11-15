package htdfservice

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"

	"github.com/orientwalt/htdf/evm/vm"
	"github.com/orientwalt/htdf/server/config"
	"github.com/orientwalt/htdf/types"

	"github.com/ethereum/go-ethereum/common"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
)

// junying-todo, 2019-11-06
// from x/core/transition.go
// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if len(data) > 0 && homestead {
		gas = params.TxGasContractCreation // 53000 -> 60000
	} else {
		gas = params.TxGas // 21000 -> 30000
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/ethparams.TxDataNonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * ethparams.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/ethparams.TxDataZeroGas < z {
			return 0, vm.ErrOutOfGas
		}
		gas += z * ethparams.TxDataZeroGas
	}
	return gas, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgSendFrom defines a SendFrom message /////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgSendFrom struct {
	From   sdk.AccAddress
	To     sdk.AccAddress
	Amount sdk.Coins
	Data   string
	Fee    sdk.StdFee
	// GasPrice  string //unit,  satoshi/gallon
	// GasWanted uint64 //unit,  gallon
}

var _ sdk.Msg = MsgSendFrom{}

// NewMsgSend is a constructor function for MsgSend
// Normal Transaction
// Default GasLimit, Default GasPrice
func NewMsgSendFromDefault(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins) MsgSendFrom {
	return MsgSendFrom{
		From:   fromaddr,
		To:     toaddr,
		Amount: amount,
		Fee:    sdk.NewStdFee(params.TxGas, config.DefaultMinGasPrices),
		// GasPrice:  config.DefaultMinGasPrices,
		// GasWanted: params.TxGas,
	}
}

// Normal Transaction
// Default GasLimit, Customized GasPrice
func NewMsgSendFrom(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins, gasprice string, gaswanted uint64) MsgSendFrom {
	return MsgSendFrom{
		From:   fromaddr,
		To:     toaddr,
		Amount: amount,
		Fee:    sdk.NewStdFee(gaswanted, gasprice),
		// GasPrice:  gasprice,
		// GasWanted: gaswanted,
	}
}

// Contract Transaction
func NewMsgSendFromForData(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins, data string, gasprice string, gaswanted uint64) MsgSendFrom {
	return MsgSendFrom{
		From:   fromaddr,
		To:     toaddr,
		Amount: amount,
		Data:   data,
		Fee:    sdk.NewStdFee(gaswanted, gasprice),
		// GasPrice:  gasprice,
		// GasWanted: gaswanted,
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

		// junying-todo, 2019-11-12
		if msg.Fee.GasWanted < params.TxGas {
			return sdk.ErrOutOfGas(fmt.Sprintf("gas must be greather than %d", params.TxGas))
		}

	} else {
		// junying-todo, 2019-11-12
		inputCode, err := hex.DecodeString(msg.Data)
		if err != nil {
			return sdk.ErrTxDecode("decoding msg.data failed. you should check msg.data")
		}
		//Intrinsic gas calc
		itrsGas, err := IntrinsicGas(inputCode, true)
		if err != nil {
			return sdk.ErrOutOfGas("intrinsic out of gas")
		}
		if msg.Fee.GasWanted < itrsGas {
			return sdk.ErrOutOfGas(fmt.Sprintf("gas must be greather than %d to pass validating", itrsGas))
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
func (msg MsgSendFrom) GetSigner() sdk.AccAddress {
	return msg.From
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

// junying-todo, 2019-11-06
func (msg MsgSendFrom) GetData() string {
	return msg.Data
}

// junying -todo, 2019-11-14
//
func (msg MsgSendFrom) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgSendFrom) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgSendFrom) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgSendFrom) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgSendFrom) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgSendFrom) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

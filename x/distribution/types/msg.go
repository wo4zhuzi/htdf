//nolint
package types

import (
	"github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
)

// Verify interface at compile time
var _, _, _ sdk.Msg = &MsgSetWithdrawAddress{}, &MsgWithdrawDelegatorReward{}, &MsgWithdrawValidatorCommission{}

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgSetWithdrawAddress struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	WithdrawAddress  sdk.AccAddress `json:"withdraw_address"`
	Fee              sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		DelegatorAddress: delAddr,
		WithdrawAddress:  withdrawAddr,
		Fee:              sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//
func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }

//
func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetWithdrawAddress) GetSigner() sdk.AccAddress {
	return msg.DelegatorAddress
}

// get the bytes for the message signer to sign on
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgSetWithdrawAddress) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.WithdrawAddress.Empty() {
		return ErrNilWithdrawAddr(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgSetWithdrawAddress) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgSetWithdrawAddress) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgSetWithdrawAddress) GetGasWanted() uint64 { return msg.GasWanted }
// func (msg MsgSetWithdrawAddress) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }
// func (msg MsgSetWithdrawAddress) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }
// func (msg MsgSetWithdrawAddress) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegatorReward struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Fee              sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) MsgWithdrawDelegatorReward {
	return MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Fee:              sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//
func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }

//
func (msg MsgWithdrawDelegatorReward) Type() string { return "withdraw_delegator_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigner() sdk.AccAddress {
	return msg.DelegatorAddress
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgWithdrawDelegatorReward) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgWithdrawDelegatorReward) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgWithdrawDelegatorReward) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgWithdrawDelegatorReward) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgWithdrawDelegatorReward) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgWithdrawDelegatorReward) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

// msg struct for validator withdraw
type MsgWithdrawValidatorCommission struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Fee              sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgWithdrawValidatorCommission(valAddr sdk.ValAddress) MsgWithdrawValidatorCommission {
	return MsgWithdrawValidatorCommission{
		ValidatorAddress: valAddr,
		Fee:              sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//
func (msg MsgWithdrawValidatorCommission) Route() string { return ModuleName }

//
func (msg MsgWithdrawValidatorCommission) Type() string { return "withdraw_validator_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorCommission) GetSigner() sdk.AccAddress {
	return sdk.AccAddress(msg.ValidatorAddress.Bytes())
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawValidatorCommission) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgWithdrawValidatorCommission) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgWithdrawValidatorCommission) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgWithdrawValidatorCommission) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgWithdrawValidatorCommission) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgWithdrawValidatorCommission) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgWithdrawValidatorCommission) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

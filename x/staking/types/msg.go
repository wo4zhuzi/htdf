package types

import (
	"encoding/json"

	"github.com/orientwalt/tendermint/crypto"

	"github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgCreateValidator{}
	_ sdk.Msg = &MsgEditValidator{}
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgUndelegate{}
	_ sdk.Msg = &MsgBeginRedelegate{}
)

//______________________________________________________________________

// MsgCreateValidator - struct for bonding transactions
type MsgCreateValidator struct {
	Description       Description    `json:"description"`
	Commission        CommissionMsg  `json:"commission"`
	MinSelfDelegation sdk.Int        `json:"min_self_delegation"`
	ValidatorAddress  sdk.ValAddress `json:"validator_address"`
	PubKey            crypto.PubKey  `json:"pubkey"`
	Value             sdk.Coin       `json:"value"`
	Fee               sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

type msgCreateValidatorJSON struct {
	Description       Description    `json:"description"`
	Commission        CommissionMsg  `json:"commission"`
	MinSelfDelegation sdk.Int        `json:"min_self_delegation"`
	ValidatorAddress  sdk.ValAddress `json:"validator_address"`
	PubKey            string         `json:"pubkey"`
	Value             sdk.Coin       `json:"value"`
	Fee               sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

// Default way to create validator. Delegator address and validator address are the same
func NewMsgCreateValidator(
	valAddr sdk.ValAddress, pubKey crypto.PubKey, selfDelegation sdk.Coin,
	description Description, commission CommissionMsg, minSelfDelegation sdk.Int,
) MsgCreateValidator {

	return MsgCreateValidator{
		Description:       description,
		ValidatorAddress:  valAddr,
		PubKey:            pubKey,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
		Fee:               sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//nolint
func (msg MsgCreateValidator) Route() string { return RouterKey }

//
func (msg MsgCreateValidator) Type() string { return "create_validator" }

// junying-todo,2019-11-14
// now delegator must be validator address
// if not, return nil
// that's because multi-sign structure is removed already.

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgCreateValidator) GetSigner() sdk.AccAddress {
	return sdk.AccAddress(msg.ValidatorAddress)
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON
// serialization of the MsgCreateValidator type.
func (msg MsgCreateValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(msgCreateValidatorJSON{
		Description:       msg.Description,
		Commission:        msg.Commission,
		ValidatorAddress:  msg.ValidatorAddress,
		PubKey:            sdk.MustBech32ifyConsPub(msg.PubKey),
		Value:             msg.Value,
		MinSelfDelegation: msg.MinSelfDelegation,
		Fee:               msg.Fee,
		// GasPrice:          msg.GasPrice,
		// GasWanted:         msg.GasWanted,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface to provide custom
// JSON deserialization of the MsgCreateValidator type.
func (msg *MsgCreateValidator) UnmarshalJSON(bz []byte) error {
	var msgCreateValJSON msgCreateValidatorJSON
	if err := json.Unmarshal(bz, &msgCreateValJSON); err != nil {
		return err
	}

	msg.Description = msgCreateValJSON.Description
	msg.Commission = msgCreateValJSON.Commission
	msg.ValidatorAddress = msgCreateValJSON.ValidatorAddress
	var err error
	msg.PubKey, err = sdk.GetConsPubKeyBech32(msgCreateValJSON.PubKey)
	if err != nil {
		return err
	}
	msg.Value = msgCreateValJSON.Value
	msg.MinSelfDelegation = msgCreateValJSON.MinSelfDelegation
	msg.Fee = msgCreateValJSON.Fee
	// msg.GasPrice = msgCreateValJSON.GasPrice
	// msg.GasWanted = msgCreateValJSON.GasWanted
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgCreateValidator) ValidateBasic() sdk.Error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Value.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if msg.Description == (Description{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "description must be included")
	}
	if msg.Commission == (CommissionMsg{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "commission must be included")
	}
	if !msg.MinSelfDelegation.GT(sdk.ZeroInt()) {
		return ErrMinSelfDelegationInvalid(DefaultCodespace)
	}
	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrSelfDelegationBelowMinimum(DefaultCodespace)
	}

	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgCreateValidator) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgCreateValidator) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgCreateValidator) GetGasWanted() uint64 { return msg.GasWanted }
// func (msg MsgCreateValidator) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }
// func (msg MsgCreateValidator) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }
// func (msg MsgCreateValidator) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

// MsgEditValidator - struct for editing a validator
type MsgEditValidator struct {
	Description
	ValidatorAddress sdk.ValAddress `json:"address"`

	// We pass a reference to the new commission rate and min self delegation as it's not mandatory to
	// update. If not updated, the deserialized rate will be zero with no way to
	// distinguish if an update was intended.
	//
	// REF: #2373
	CommissionRate    *sdk.Dec   `json:"commission_rate"`
	MinSelfDelegation *sdk.Int   `json:"min_self_delegation"`
	Fee               sdk.StdFee `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgEditValidator(valAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *sdk.Int) MsgEditValidator {
	return MsgEditValidator{
		Description:       description,
		CommissionRate:    newRate,
		ValidatorAddress:  valAddr,
		MinSelfDelegation: newMinSelfDelegation,
		Fee:               sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
	}
}

//nolint
func (msg MsgEditValidator) Route() string { return RouterKey }

//
func (msg MsgEditValidator) Type() string { return "edit_validator" }

//
func (msg MsgEditValidator) GetSigner() sdk.AccAddress {
	return sdk.AccAddress(msg.ValidatorAddress)
}

// get the bytes for the message signer to sign on
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgEditValidator) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil validator address")
	}

	if msg.Description == (Description{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "transaction must include some information to modify")
	}

	if msg.MinSelfDelegation != nil && !(*msg.MinSelfDelegation).GT(sdk.ZeroInt()) {
		return ErrMinSelfDelegationInvalid(DefaultCodespace)
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.LT(sdk.ZeroDec()) {
			return sdk.NewError(DefaultCodespace, CodeInvalidInput, "commission rate must be between 0 and 1, inclusive")
		}
	}

	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgEditValidator) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgEditValidator) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgEditValidator) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgEditValidator) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgEditValidator) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgEditValidator) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin       `json:"amount"`
	Fee              sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Amount:           amount,
		Fee:              sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//nolint
func (msg MsgDelegate) Route() string { return RouterKey }

//
func (msg MsgDelegate) Type() string { return "delegate" }

//
func (msg MsgDelegate) GetSigner() sdk.AccAddress {
	return msg.DelegatorAddress
}

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgDelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgDelegate) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgDelegate) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgDelegate) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgDelegate) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgDelegate) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgDelegate) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgBeginRedelegate struct {
	DelegatorAddress    sdk.AccAddress `json:"delegator_address"`
	ValidatorSrcAddress sdk.ValAddress `json:"validator_src_address"`
	ValidatorDstAddress sdk.ValAddress `json:"validator_dst_address"`
	Amount              sdk.Coin       `json:"amount"`
	Fee                 sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgBeginRedelegate(delAddr sdk.AccAddress, valSrcAddr,
	valDstAddr sdk.ValAddress, amount sdk.Coin) MsgBeginRedelegate {

	return MsgBeginRedelegate{
		DelegatorAddress:    delAddr,
		ValidatorSrcAddress: valSrcAddr,
		ValidatorDstAddress: valDstAddr,
		Amount:              amount,
		Fee:                 sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//nolint
func (msg MsgBeginRedelegate) Route() string { return RouterKey }

//
func (msg MsgBeginRedelegate) Type() string { return "begin_redelegate" }

//
func (msg MsgBeginRedelegate) GetSigner() sdk.AccAddress {
	return msg.DelegatorAddress
}

// get the bytes for the message signer to sign on
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgBeginRedelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorSrcAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.ValidatorDstAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSharesAmount(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgBeginRedelegate) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgBeginRedelegate) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgBeginRedelegate) GetGasWanted() uint64 { return msg.GasWanted }
// func (msg MsgBeginRedelegate) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }
// func (msg MsgBeginRedelegate) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }
// func (msg MsgBeginRedelegate) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

// MsgUndelegate - struct for unbonding transactions
type MsgUndelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin       `json:"amount"`
	Fee              sdk.StdFee     `json:"fee"`
	// GasWanted        uint64         `json:"gas_wanted"`
	// GasPrice         string         `json:"gas_price"`
}

func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgUndelegate {
	return MsgUndelegate{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Amount:           amount,
		Fee:              sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:        uint64(10000),
		// GasPrice:         config.DefaultMinGasPrices,
	}
}

//nolint
func (msg MsgUndelegate) Route() string { return RouterKey }

//
func (msg MsgUndelegate) Type() string { return "begin_unbonding" }

//
func (msg MsgUndelegate) GetSigner() sdk.AccAddress { return msg.DelegatorAddress }

// get the bytes for the message signer to sign on
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUndelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSharesAmount(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgUndelegate) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgUndelegate) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgUndelegate) GetGasWanted() uint64 { return msg.GasWanted }
// func (msg MsgUndelegate) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }
// func (msg MsgUndelegate) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }
// func (msg MsgUndelegate) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

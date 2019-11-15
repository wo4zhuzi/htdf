package crisis

import (
	"github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
)

// MsgVerifyInvariant - message struct to verify a particular invariance
type MsgVerifyInvariant struct {
	Sender              sdk.AccAddress `json:"sender"`
	InvariantModuleName string         `json:"invariant_module_name"`
	InvariantRoute      string         `json:"invariant_route"`
	Fee                 sdk.StdFee     `json:"fee"`
	// GasWanted           uint64         `json:"gas_wanted"`
	// GasPrice            string         `json:"gas_price"`
}

// ensure Msg interface compliance at compile time
var _ sdk.Msg = &MsgVerifyInvariant{}

// NewMsgVerifyInvariant creates a new MsgVerifyInvariant object
func NewMsgVerifyInvariant(sender sdk.AccAddress, invariantModuleName,
	invariantRoute string) MsgVerifyInvariant {

	return MsgVerifyInvariant{
		Sender:              sender,
		InvariantModuleName: invariantModuleName,
		InvariantRoute:      invariantRoute,
		Fee:                 sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
		// GasWanted:           uint64(10000),
		// GasPrice:            config.DefaultMinGasPrices,
	}
}

//nolint
func (msg MsgVerifyInvariant) Route() string { return ModuleName }

//
func (msg MsgVerifyInvariant) Type() string { return "verify_invariant" }

// get the bytes for the message signer to sign on
func (msg MsgVerifyInvariant) GetSigner() sdk.AccAddress { return msg.Sender }

// GetSignBytes gets the sign bytes for the msg MsgVerifyInvariant
func (msg MsgVerifyInvariant) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgVerifyInvariant) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return ErrNilSender(DefaultCodespace)
	}
	return nil
}

// FullInvariantRoute - get the messages full invariant route
func (msg MsgVerifyInvariant) FullInvariantRoute() string {
	return msg.InvariantModuleName + "/" + msg.InvariantRoute
}

// junying -todo, 2019-11-14
//
func (msg MsgVerifyInvariant) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgVerifyInvariant) SetFee(fee sdk.StdFee) { msg.Fee = fee }

// func (msg MsgVerifyInvariant) GetGasWanted() uint64 { return msg.GasWanted }

// //
// func (msg MsgVerifyInvariant) SetGasWanted(gaswanted uint64) { msg.GasWanted = gaswanted }

// //
// func (msg MsgVerifyInvariant) GetGasPrice() uint64 {
// 	gasprice, err := types.ParseCoin(msg.GasPrice)
// 	if err != nil {
// 		return 0
// 	}
// 	amount := gasprice.Amount
// 	return amount.Uint64()
// }

// //
// func (msg MsgVerifyInvariant) SetGasPrice(gasprice string) { msg.GasPrice = gasprice }

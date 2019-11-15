package slashing

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
)

var cdc = codec.New()

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr sdk.ValAddress `json:"address"` // address of the validator operator
	Fee           sdk.StdFee     `json:"fee"`
}

func NewMsgUnjail(validatorAddr sdk.ValAddress) MsgUnjail {
	return MsgUnjail{
		ValidatorAddr: validatorAddr,
		Fee:           sdk.NewStdFee(uint64(10000), config.DefaultMinGasPrices),
	}
}

//nolint
func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return "unjail" }
func (msg MsgUnjail) GetSigner() sdk.AccAddress {
	return sdk.AccAddress(msg.ValidatorAddr)
}

// get the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := cdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	return nil
}

// junying -todo, 2019-11-14
//
func (msg MsgUnjail) GetFee() sdk.StdFee { return msg.Fee }

//
func (msg MsgUnjail) SetFee(fee sdk.StdFee) { msg.Fee = fee }

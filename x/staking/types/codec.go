package types

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "htdf/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "htdf/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "htdf/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "htdf/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "htdf/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}

package bank

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	// cdc.RegisterConcrete(MsgSend{}, "htdf/MsgSend", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

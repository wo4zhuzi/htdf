package htdfservice

import (
	"github.com/orientwalt/htdf/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSendFrom{}, "htdfservice/send", nil)
	cdc.RegisterConcrete(MsgAdd{}, "htdfservice/add", nil)
}

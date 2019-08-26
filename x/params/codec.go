package params

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*ParamSet)(nil), nil)
}

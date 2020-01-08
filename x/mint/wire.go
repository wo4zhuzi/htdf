package mint

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	// Not Register mint codec in app, deprecated now
	//cdc.RegisterConcrete(Minter{}, "htdf/mint/Minter", nil)
	cdc.RegisterConcrete(&Params{}, "htdf/mint/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

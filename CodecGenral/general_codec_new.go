package CodecGenral

import (
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
)

var InstCodecNew = codec.New()

func init() {
	RegisterNew(InstCodecNew)
}

// Register concrete types on codec codec
func RegisterNew(cdc *codec.Codec) {

	RegisterOld(cdc)

	//cdc.RegisterConcrete(&mint.Params{}, "htdf/mint/Params", nil)
	//cdc.RegisterConcrete(&sdk.Dec{}, "htdf/mint/rewards", nil)
	cdc.RegisterConcrete(&sdk.Int{}, "htdf/mint/Int", nil)
}

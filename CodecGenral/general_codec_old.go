package CodecGenral

import (
	"github.com/orientwalt/htdf/codec"
	newevmtypes "github.com/orientwalt/htdf/evm/types"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	htdfservice "github.com/orientwalt/htdf/x/core"
	"github.com/orientwalt/htdf/x/crisis"
	distr "github.com/orientwalt/htdf/x/distribution"
	"github.com/orientwalt/htdf/x/gov"
	"github.com/orientwalt/htdf/x/guardian"
	"github.com/orientwalt/htdf/x/mint"
	"github.com/orientwalt/htdf/x/params"
	"github.com/orientwalt/htdf/x/service"
	"github.com/orientwalt/htdf/x/slashing"
	"github.com/orientwalt/htdf/x/upgrade"

	stake "github.com/orientwalt/htdf/x/staking"
)

var InstCodecOld = codec.New()

func init() {
	RegisterOld(InstCodecOld)
}

func RegisterOld(cdc *codec.Codec) {
	newevmtypes.RegisterCodec(cdc)
	htdfservice.RegisterCodec(cdc)
	params.RegisterCodec(cdc) // only used by querier
	mint.RegisterCodec(cdc)   // only used by querier
	// bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	crisis.RegisterCodec(cdc)
}

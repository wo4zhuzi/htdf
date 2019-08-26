package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/crypto/keys"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}

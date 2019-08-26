package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}

package rest

import (
	"fmt"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"

	svrConfig "github.com/orientwalt/htdf/server/config"

	"github.com/gorilla/mux"
)

const (
	restName = "custom"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		r.HandleFunc(fmt.Sprintf("/%s/send", storeName), SendTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/create", storeName), CreateTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/sign", storeName), SignTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	}
	r.HandleFunc(fmt.Sprintf("/%s/broadcast", storeName), BroadcastTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}

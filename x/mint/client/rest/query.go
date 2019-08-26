package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/x/mint"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/minting/parameters",
		queryParamsHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/inflation",
		queryInflationHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/annual-provisions",
		queryAnnualProvisionsHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryParameters)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryInflationHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryInflation)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryAnnualProvisionsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryAnnualProvisions)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

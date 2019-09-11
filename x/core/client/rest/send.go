package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/accounts/keystore"
	htdfRest "github.com/orientwalt/htdf/accounts/rest"
	hsign "github.com/orientwalt/htdf/accounts/signs"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/crypto/keys/keyerror"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/utils/unit_convert"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/bank"
	htdfservice "github.com/orientwalt/htdf/x/core"
	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to"`
	Amount   sdk.Coins    `json:"amount"`
	Data     string       `json:"data"`
	GasPrice string       `json:"gas_price"` // uint: HTDF/gallon
	GasLimit string       `json:"gas_limit"` // unit: gallon
}

var msgCdc = codec.New()

func init() {
	bank.RegisterCodec(msgCdc)
}

// SendTxRequestHandlerFn - http request handler to send coins to a address.
func SendTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SendReq
		var mreq htdfRest.SendShiftReq
		if !rest.ReadRESTReq(w, r, cdc, &mreq) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq.From = mreq.BaseReq.From
		req.BaseReq.Memo = mreq.BaseReq.Memo
		req.BaseReq.ChainID = mreq.BaseReq.ChainID
		req.BaseReq.AccountNumber = mreq.BaseReq.AccountNumber
		req.BaseReq.Sequence = mreq.BaseReq.Sequence
		req.BaseReq.Fees = unit_convert.BigCoinsToDefaultCoins(mreq.BaseReq.Fees)
		req.BaseReq.GasPrices = mreq.BaseReq.GasPrices
		req.BaseReq.Gas = mreq.BaseReq.Gas
		req.BaseReq.GasAdjustment = mreq.BaseReq.GasAdjustment
		req.BaseReq.Simulate = mreq.BaseReq.Simulate
		req.To = mreq.To
		req.Data = mreq.Data
		req.GasPrice = mreq.GasPrice
		req.GasLimit = mreq.GasLimit
		fmt.Printf("req.BaseReq.Fees=%v\n", req.BaseReq.Fees)

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		toAddr, err := sdk.AccAddressFromBech32(req.To)
		if err != nil {

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		gas, err := strconv.ParseUint(req.BaseReq.Gas, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// when access smart contract, extract gas field
		var gasPrice, gasLimit uint64
		if len(req.Data) > 0 {
			gasPrice, err = strconv.ParseUint(unit_convert.BigAmountToDefaultAmount(req.GasPrice), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			gasLimit, err = strconv.ParseUint(req.GasLimit, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		fmt.Printf("gas=%d|gasPrice=%d|gasLimit=%d\n", gas, gasPrice, gasLimit)

		msg := htdfservice.NewMsgSendFromForData(fromAddr, toAddr, unit_convert.BigCoinsToDefaultCoins(mreq.Amount), req.Data, gas, gasPrice, gasLimit)
		CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, mreq.BaseReq.Password, []sdk.Msg{msg}, cdc)

	}
}

//-----------------------------------------------------------------------------
// Building / Sending utilities

// CompleteAndBroadcastTxREST implements a utility function that facilitates
// sending a series of messages in a signed tx. In addition, it will handle
// tx gas simulation and estimation.
//
// NOTE: Also see CompleteAndBroadcastTxCLI.
func CompleteAndBroadcastTxREST(w http.ResponseWriter, cliCtx context.CLIContext,
	baseReq rest.BaseReq, password string, msgs []sdk.Msg, cdc *codec.Codec) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := client.ParseGas(baseReq.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), baseReq.AccountNumber,
		baseReq.Sequence, gas, gasAdj, baseReq.Simulate,
		baseReq.ChainID, baseReq.Memo, baseReq.Fees, baseReq.GasPrices,
	)

	// get fromaddr
	fromaddr := msgs[0].(htdfservice.MsgSendFrom).GetFromAddr()

	txBldr, err = hscorecli.PrepareTxBuilder(txBldr, cliCtx, fromaddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if baseReq.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, client.ErrInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if baseReq.Simulate {
			rest.WriteSimulationResponse(w, cdc, txBldr.Gas())
			return
		}
	}

	bech32 := sdk.AccAddress.String(fromaddr)
	account := accounts.Account{Address: bech32}
	privkey, err := keystore.GetPrivKey(account, password, "")

	if err != nil {
		return
	}

	txBytes, err := hsign.BuildAndSign(txBldr, privkey, msgs)
	if keyerror.IsErrKeyNotFound(err) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		rest.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	cliCtx = cliCtx.WithBroadcastMode("sync")
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
}

package htdfservice

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	vmcore "github.com/orientwalt/htdf/evm/core"
	"github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"
	appParams "github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
)

//
type SendTxResp struct {
	ErrCode         sdk.CodeType `json:"code"`
	ErrMsg          string       `json:"message"`
	ContractAddress string       `json:"contract_address"`
	EvmOutput       string       `json:"evm_output"`
}

//
func (rsp SendTxResp) String() string {
	rsp.ErrMsg = sdk.GetErrMsg(rsp.ErrCode)
	data, _ := json.Marshal(&rsp)
	return string(data)
}

// New HTDF Message Handler
// connected to handler.go
// HandleMsgSendFrom, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgSendFrom:
			return HandleMsgSendFrom(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp SendTxResp
	fmt.Printf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCode_Param
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// junying-todo, 2019-08-26
func HandleMsgSendFrom(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSendFrom) sdk.Result {
	// initialize
	var sendTxResp SendTxResp
	var gasUsed uint64
	var evmOutput string
	var err error

	if !msg.To.Empty() {
		// open smart contract
		evmOutput, gasUsed, err = HandleOpenContract(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		if err != nil {
			sendTxResp.ErrCode = sdk.ErrCode_OpenContract
		}
		sendTxResp.EvmOutput = evmOutput
	} else {
		// create smart contract
		evmOutput, gasUsed, err = HandleCreateContract(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		if err != nil {
			sendTxResp.ErrCode = sdk.ErrCode_CreateContract
		}
		sendTxResp.ContractAddress = evmOutput
	}
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: gasUsed}
}

//
func FeeCollecting(ctx sdk.Context,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	stateDB *state.CommitStateDB,
	gasused uint64,
	gasprice *big.Int) {
	gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(gasused), gasprice)
	fmt.Printf("gasUsed=%s\n", gasUsed.String())
	feeCollectionKeeper.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
	stateDB.Commit(false)
}

// junying-todo, 2019-08-26
func HandleOpenContract(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSendFrom) (evmOutput string, gasUsed uint64, err error) {

	fmt.Printf("Handling MsgSendFrom with No Contract.\n")
	fmt.Println(" HandleOpenContract0:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, keyStorage, keyCode)
	if err != nil {
		evmOutput = fmt.Sprintf("newStateDB error\n")
		return
	}

	fromAddress := sdk.ToEthAddress(msg.From)
	toAddress := sdk.ToEthAddress(msg.To)

	fmt.Printf("fromAddr|appFormat=%s|ethFormat=%s|\n", msg.From.String(), fromAddress.String())
	fmt.Printf("toAddress|appFormat=%s|ethFormat=%s|\n", msg.To.String(), toAddress.String())

	fmt.Printf("fromAddress|testBalance=%v\n", stateDB.GetBalance(fromAddress))
	fmt.Printf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evmCtx := vmcore.NewEVMContext(msg, &fromAddress, uint64(ctx.BlockHeight()))
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	contractRef := vm.AccountRef(fromAddress)

	inputCode, err := hex.DecodeString(msg.Data)
	if err != nil {
		evmOutput = fmt.Sprintf("DecodeString error\n")
		return
	}

	fmt.Printf("inputCode=%s\n", hex.EncodeToString(inputCode))

	transferAmount := msg.Amount.AmountOf(sdk.DefaultDenom).BigInt()

	fmt.Printf("transferAmount: %d\n", transferAmount)
	st := NewStateTransition(evm, msg, stateDB)

	fmt.Printf("gasPrice=%d|gasLimit=%d\n", msg.GasPrice, msg.GasLimit)

	// commented by junying, 2019-08-22
	// subtract gaslimit*gasprice from sender
	err = st.buyGas()
	if err != nil {
		evmOutput = fmt.Sprintf("buyGas error|err=%s\n", err)
		return
	}

	// Intrinsic gas calc
	// commented by junying, 2019-08-22
	// default non-contract tx gas: 21000
	// default contract tx gas: 53000 + f(tx.data)
	itrsGas, err := auth.IntrinsicGas(inputCode, true)
	fmt.Printf("itrsGas|gas=%d\n", itrsGas)
	// commented by junying, 2019-08-22
	// check if tx.gas >= calculated gas
	err = st.useGas(itrsGas)
	if err != nil {
		evmOutput = fmt.Sprintf("useGas error|err=%s\n", err)
		return
	}

	// commented by junying, 2019-08-22
	// 1. cantransfer check
	// 2. create receiver account if no exists
	// 3. execute contract & calculate gas
	fmt.Println(" HandleOpenContract1:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	outputs, gasLeftover, err := evm.Call(contractRef, toAddress, inputCode, st.gas, transferAmount)
	fmt.Println(" HandleOpenContract2:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	if err != nil {
		fmt.Printf("evm call error|err=%s\n", err)
		// junying-todo, 2019-11-05
		gasUsed = msg.GasLimit
		evmOutput = fmt.Sprintf("evm call error|err=%s\n", err)
	} else {
		st.gas = gasLeftover
		// junying-todo, 2019-08-22
		// refund(add) remaining to sender
		st.refundGas()
		fmt.Printf("gasUsed=%d\n", st.gasUsed())
		gasUsed = st.gasUsed()
		evmOutput = hex.EncodeToString(outputs)
	}
	FeeCollecting(ctx, feeCollectionKeeper, stateDB, gasUsed, st.gasPrice)
	return
}

// junying-todo, 2019-08-26
func HandleCreateContract(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSendFrom) (evmOutput string, gasUsed uint64, err error) {

	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, keyStorage, keyCode)
	if err != nil {
		evmOutput = fmt.Sprintf("newStateDB error\n")
		return
	}
	fromAddress := sdk.ToEthAddress(msg.From)
	toAddress := sdk.ToEthAddress(msg.To)

	fmt.Printf("fromAddr|appFormat=%s|ethFormat=%s|\n", msg.From.String(), fromAddress.String())
	fmt.Printf("toAddress|appFormat=%s|ethFormat=%s|\n", msg.To.String(), toAddress.String())
	fmt.Printf("fromAddress|Balance=%v\n", stateDB.GetBalance(fromAddress))

	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	fmt.Printf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	evmCtx := vmcore.NewEVMContext(msg, &fromAddress, uint64(ctx.BlockHeight()))
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	contractRef := vm.AccountRef(fromAddress)

	fmt.Printf("blockHeight=%d|IsHomestead=%v|IsDAOFork=%v|IsEIP150=%v|IsEIP155=%v|IsEIP158=%v|IsByzantium=%v\n", ctx.BlockHeight(), evm.ChainConfig().IsHomestead(evm.BlockNumber),
		evm.ChainConfig().IsDAOFork(evm.BlockNumber), evm.ChainConfig().IsEIP150(evm.BlockNumber),
		evm.ChainConfig().IsEIP155(evm.BlockNumber), evm.ChainConfig().IsEIP158(evm.BlockNumber),
		evm.ChainConfig().IsByzantium(evm.BlockNumber))

	inputCode, err := hex.DecodeString(msg.Data)
	if err != nil {
		evmOutput = fmt.Sprintf("DecodeString error\n")
		return
	}

	fmt.Printf("inputCode=%s\n", hex.EncodeToString(inputCode))

	st := NewStateTransition(evm, msg, stateDB)

	fmt.Printf("gasPrice=%d|gasLimit=%d\n", msg.GasPrice, msg.GasLimit)

	err = st.buyGas()
	if err != nil {
		evmOutput = fmt.Sprintf("buyGas error|err=%s\n", err)
		return
	}

	//Intrinsic gas calc
	itrsGas, err := auth.IntrinsicGas(inputCode, true)
	fmt.Printf("itrsGas|gas=%d\n", itrsGas)
	err = st.useGas(itrsGas)
	if err != nil {
		evmOutput = fmt.Sprintf("useGas error|err=%s\n", err)
		return
	}

	_, contractAddr, gasLeftover, err := evm.Create(contractRef, inputCode, st.gas, big.NewInt(0))
	if err != nil {
		fmt.Printf("evm Create error|err=%s\n", err)
		// junying-todo, 2019-11-05
		gasUsed = msg.GasLimit
		evmOutput = fmt.Sprintf("evm Create error|err=%s\n", err)
	} else {
		st.gas = gasLeftover
		st.refundGas()
		gasUsed = st.gasUsed()
		evmOutput = sdk.ToAppAddress(contractAddr).String()
	}

	fmt.Printf("Create contract ok,contractAddr|appFormat=%s|ethFormat=%s\n", sdk.ToAppAddress(contractAddr).String(), contractAddr.String())
	FeeCollecting(ctx, feeCollectionKeeper, stateDB, gasUsed, st.gasPrice)
	return
}

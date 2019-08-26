package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/orientwalt/htdf/app/protocol"
	"github.com/spf13/viper"

	"github.com/orientwalt/tendermint/libs/log"
	tmtypes "github.com/orientwalt/tendermint/types"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/x/auth"

	//"github.com/orientwalt/htdf/x/mint"

	htdfservice "github.com/orientwalt/htdf/x/core"

	sdk "github.com/orientwalt/htdf/types"
	abci "github.com/orientwalt/tendermint/abci/types"
	dbm "github.com/orientwalt/tendermint/libs/db"

	v0 "github.com/orientwalt/htdf/app/v0"
	cfg "github.com/orientwalt/tendermint/config"
	cmn "github.com/orientwalt/tendermint/libs/common"
)

const (
	appName = "HtdfServiceApp"

	appPrometheusNamespace = "htdf"
	//
	RouterKey = "htdfservice"
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"

	FlagReplay = "replay-last-block"

	DefaultCacheSize = 100 // Multistore saves last 100 blocks

	DefaultSyncableHeight = 10000 // Multistore saves a snapshot every 10000 blocks
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.hscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.hsd")
)

// Extended ABCI application
type HtdfServiceApp struct {
	*BaseApp
	// cdc *codec.Codec

	invCheckPeriod uint
}

// NewHtdfServiceApp is a constructor function for htdfServiceApp
func NewHtdfServiceApp(logger log.Logger, config *cfg.InstrumentationConfig, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *HtdfServiceApp {

	cdc := MakeLatestCodec()

	bApp := NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	var app = &HtdfServiceApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}
	protocolKeeper := sdk.NewProtocolKeeper(protocol.KeyMain)
	fmt.Print("/---------protocolKeeper----------/", protocolKeeper, "\n")
	engine := protocol.NewProtocolEngine(protocolKeeper)
	app.SetProtocolEngine(&engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())

	var err error
	if viper.GetBool(FlagReplay) {
		lastHeight := Replay(app.logger)
		err = app.LoadVersion(lastHeight, protocol.KeyMain, true)
	} else {
		err = app.LoadLatestVersion(protocol.KeyMain)
	} // app is now sealed
	if err != nil {
		cmn.Exit(err.Error())
	}
	//Duplicate prometheus config
	appPrometheusConfig := *config
	//Change namespace to appName
	appPrometheusConfig.Namespace = appPrometheusNamespace
	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, app.invCheckPeriod, &appPrometheusConfig))
	//engine.Add(v1.NewProtocolV0(1, logger, protocolKeeper, app.checkInvariant, app.trackCoinFlow, &appPrometheusConfig))
	// engine.Add(v2.NewProtocolV1(2, ...))
	fmt.Print("----->	", app.GetKVStore(protocol.KeyMain), "\n")
	loaded, current := engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())
	engine.GetCurrentProtocol().InitMetrics(app.cms)
	return app
}

func (app *HtdfServiceApp) ExportOrReplay(replayHeight int64) (replay bool, height int64) {
	lastBlockHeight := app.BaseApp.LastBlockHeight()
	if replayHeight > lastBlockHeight {
		replayHeight = lastBlockHeight
	}

	if lastBlockHeight-replayHeight <= DefaultCacheSize {
		err := app.LoadVersion(replayHeight, protocol.KeyMain, false)
		if err != nil {
			cmn.Exit(err.Error())
		}
		return false, replayHeight
	}

	loadHeight := app.replayToHeight(replayHeight, app.logger)
	err := app.LoadVersion(loadHeight, protocol.KeyMain, true)
	if err != nil {
		cmn.Exit(err.Error())
	}
	app.logger.Info(fmt.Sprintf("Load store at %d, start to replay to %d", loadHeight, replayHeight))
	return true, replayHeight

}

// export the state of iris for a genesis file
func (app *HtdfServiceApp) ExportAppStateAndValidators(forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, []string{})
}

// load a particular height
func (app *HtdfServiceApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.KeyMain, false)
}

// MakeCodec generates the necessary codecs for Amino
func MakeLatestCodec() *codec.Codec {
	var cdc = v0.MakeLatestCodec() // replace with latest protocol version
	return cdc
}

func (app *HtdfServiceApp) replayToHeight(replayHeight int64, logger log.Logger) int64 {
	loadHeight := int64(0)
	logger.Info("Please make sure the replay height is smaller than the latest block height.")
	if replayHeight >= DefaultSyncableHeight {
		loadHeight = replayHeight - replayHeight%DefaultSyncableHeight
	} else {
		// version 1 will always be kept
		loadHeight = 1
	}
	logger.Info("This replay operation will change the application store, backup your node home directory before proceeding!!")
	logger.Info("Are you sure to proceed? (y/n)")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	confirm := strings.ToLower(strings.TrimSpace(input))
	if confirm != "y" && confirm != "yes" {
		cmn.Exit("Replay operation aborted.")
	}
	return loadHeight
}

func (app *HtdfServiceApp) openContract(globalCtx sdk.Context, msgSendFrom htdfservice.MsgSendFrom) (err error, evmOutput string) {

	// fmt.Printf("open contract\n")

	// ctx := app.NewContext(false, abci.Header{})
	// stateDB, err := evmstate.NewCommitStateDB(ctx, &app.accountMapper, app.keyStorage, app.keyCode)
	// if err != nil {
	// 	fmt.Printf("newStateDB error\n")
	// 	return err, ""
	// }

	// fromAddress := apptypes.ToEthAddress(msgSendFrom.From)
	// toAddress := apptypes.ToEthAddress(msgSendFrom.To)

	// fmt.Printf("fromAddr|appFormat=%s|ethFormat=%s|\n", msgSendFrom.From.String(), fromAddress.String())
	// fmt.Printf("toAddress|appFormat=%s|ethFormat=%s|\n", msgSendFrom.To.String(), toAddress.String())

	// fmt.Printf("fromAddress|testBalance=%v\n", stateDB.GetBalance(fromAddress))
	// fmt.Printf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	// config := appParams.MainnetChainConfig
	// logConfig := vm.LogConfig{}
	// structLogger := vm.NewStructLogger(&logConfig)
	// vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	// evmCtx := ec.NewEVMContext(msgSendFrom, &fromAddress, uint64(globalCtx.BlockHeight()))
	// evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	// contractRef := vm.AccountRef(fromAddress)

	// inputCode, err := hex.DecodeString(msgSendFrom.Data)
	// if err != nil {
	// 	fmt.Printf("DecodeString error\n")
	// 	return err, ""
	// }

	// fmt.Printf("inputCode=%s\n", hex.EncodeToString(inputCode))

	// transferAmount := msgSendFrom.Amount.AmountOf(unit_convert.DefaultDenom).BigInt()
	// if len(msgSendFrom.Data) > 0 {
	// 	transferAmount = big.NewInt(0) //when contract invoke , can't transfer; otherwise panic happen
	// }

	// st := NewStateTransition(evm, msgSendFrom, stateDB)

	// fmt.Printf("gas=%d|gasPrice=%d|gasLimit=%d\n", msgSendFrom.Gas, msgSendFrom.GasPrice, msgSendFrom.GasLimit)

	// err = st.buyGas()
	// if err != nil {
	// 	fmt.Printf("buyGas error|err=%s\n", err)
	// 	return err, ""
	// }

	// //Intrinsic gas calc
	// itrsGas, err := IntrinsicGas(inputCode, true, true)
	// fmt.Printf("itrsGas|gas=%d\n", itrsGas)
	// err = st.useGas(itrsGas)
	// if err != nil {
	// 	fmt.Printf("useGas error|err=%s\n", err)
	// 	return err, ""
	// }

	// outputs, gasLeftover, vmerr := evm.Call(contractRef, toAddress, inputCode, st.gas, transferAmount)
	// if err != nil {
	// 	fmt.Printf("evm call error|err=%s\n", vmerr)
	// 	return vmerr, ""
	// }

	// st.gas = gasLeftover

	// st.refundGas()

	// fmt.Printf("gasUsed=%d\n", st.gasUsed())

	// // gasUsedValue
	// gasUsedValue := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
	// fmt.Printf("gasUsedValue=%s\n", gasUsedValue.String())
	// app.feeCollectionKeeper.AddCollectedFees(globalCtx, sdk.Coins{sdk.NewCoin(unit_convert.DefaultDenom, sdk.NewIntFromBigInt(gasUsedValue))})

	// fmt.Printf("evm call end|outputs=%x\n", outputs)

	// stateDB.Commit(false)

	// return nil, hex.EncodeToString(outputs)
	return err, ""
}

func (app *HtdfServiceApp) CreateContract(globalCtx sdk.Context, msgSendFrom htdfservice.MsgSendFrom) (err error, retContractAddr string) {

	// ctx := app.NewContext(false, abci.Header{})
	// stateDB, err := evmstate.NewCommitStateDB(ctx, &app.accountKeeper, app.keyStorage, app.keyCode)
	// if err != nil {
	// 	fmt.Printf("newStateDB error\n")
	// 	return err, ""
	// }

	// fromAddress := apptypes.ToEthAddress(msgSendFrom.From)
	// toAddress := apptypes.ToEthAddress(msgSendFrom.To)

	// fmt.Printf("fromAddr|appFormat=%s|ethFormat=%s|\n", msgSendFrom.From.String(), fromAddress.String())
	// fmt.Printf("toAddress|appFormat=%s|ethFormat=%s|\n", msgSendFrom.To.String(), toAddress.String())
	// fmt.Printf("fromAddress|Balance=%v\n", stateDB.GetBalance(fromAddress))

	// config := appParams.MainnetChainConfig
	// logConfig := vm.LogConfig{}
	// structLogger := vm.NewStructLogger(&logConfig)
	// vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	// fmt.Printf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	// evmCtx := ec.NewEVMContext(msgSendFrom, &fromAddress, uint64(globalCtx.BlockHeight()))
	// evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	// contractRef := vm.AccountRef(fromAddress)

	// fmt.Printf("blockHeight=%d|IsHomestead=%v|IsDAOFork=%v|IsEIP150=%v|IsEIP155=%v|IsEIP158=%v|IsByzantium=%v\n", globalCtx.BlockHeight(), evm.ChainConfig().IsHomestead(evm.BlockNumber),
	// 	evm.ChainConfig().IsDAOFork(evm.BlockNumber), evm.ChainConfig().IsEIP150(evm.BlockNumber),
	// 	evm.ChainConfig().IsEIP155(evm.BlockNumber), evm.ChainConfig().IsEIP158(evm.BlockNumber),
	// 	evm.ChainConfig().IsByzantium(evm.BlockNumber))

	// inputCode, err := hex.DecodeString(msgSendFrom.Data)
	// if err != nil {
	// 	fmt.Printf("DecodeString error\n")
	// 	return err, ""
	// }

	// fmt.Printf("inputCode=%s\n", hex.EncodeToString(inputCode))

	// st := NewStateTransition(evm, msgSendFrom, stateDB)

	// fmt.Printf("gas=%d|gasPrice=%d|gasLimit=%d\n", msgSendFrom.Gas, msgSendFrom.GasPrice, msgSendFrom.GasLimit)

	// err = st.buyGas()
	// if err != nil {
	// 	fmt.Printf("buyGas error|err=%s\n", err)
	// 	return err, ""
	// }

	// //Intrinsic gas calc
	// itrsGas, err := IntrinsicGas(inputCode, true, true)
	// fmt.Printf("itrsGas|gas=%d\n", itrsGas)
	// err = st.useGas(itrsGas)
	// if err != nil {
	// 	fmt.Printf("useGas error|err=%s\n", err)
	// 	return err, ""
	// }

	// _, contractAddr, gasLeftover, vmerr := evm.Create(contractRef, inputCode, st.gas, big.NewInt(0))
	// if vmerr != nil {
	// 	fmt.Printf("evm Create error|err=%s\n", vmerr)
	// 	return vmerr, ""
	// }
	// st.gas = gasLeftover

	// st.refundGas()

	// fmt.Printf("gasUsed=%d\n", st.gasUsed())

	// // gasUsedValue
	// gasUsedValue := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
	// fmt.Printf("gasUsedValue=%s\n", gasUsedValue.String())
	// app.feeCollectionKeeper.AddCollectedFees(globalCtx, sdk.Coins{sdk.NewCoin(unit_convert.DefaultDenom, sdk.NewIntFromBigInt(gasUsedValue))})

	// fmt.Printf("Create contract ok,contractAddr|appFormat=%s|ethFormat=%s\n", apptypes.ToAppAddress(contractAddr).String(), contractAddr.String())

	// stateDB.Commit(false)

	// return nil, apptypes.ToAppAddress(contractAddr).String()
	return err, ""
}

func (app *HtdfServiceApp) Transition(ctx sdk.Context, inputMsg sdk.Msg) (result sdk.Result) {
	// var sendTxResp apptypes.SendTxResp

	// switch msg := inputMsg.(type) {

	// case htdfservice.MsgSendFrom:

	// 	// classic transfer
	// 	if len(msg.Data) == 0 {
	// 		return htdfservice.HandleMsgSendFrom(ctx, app.bankKeeper, msg)
	// 	}

	// 	fmt.Printf("FeeTotal1=%v\n", app.feeCollectionKeeper.GetCollectedFees(ctx))

	// 	if !msg.To.Empty() {
	// 		//open smart contract
	// 		fmt.Printf("openContract\n")
	// 		genErr, evmOutput := app.openContract(ctx, msg)
	// 		if genErr != nil {
	// 			fmt.Printf("openContract error|err=%s\n", genErr)
	// 			sendTxResp.ErrCode = apptypes.ErrCode_OpenContract
	// 			return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
	// 		}

	// 		sendTxResp.EvmOutput = evmOutput
	// 		return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}

	// 	} else {
	// 		// new smart contract
	// 		fmt.Printf("create contract\n")
	// 		err, contractAddr := app.CreateContract(ctx, msg)
	// 		if err != nil {
	// 			fmt.Printf("CreateContract error|err=%s\n", err)
	// 			sendTxResp.ErrCode = apptypes.ErrCode_CreateContract
	// 			return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
	// 		}

	// 		sendTxResp.ContractAddress = contractAddr
	// 		return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
	// 	}

	// 	fmt.Printf("FeeTotal2=%v\n", app.feeCollectionKeeper.GetCollectedFees(ctx))

	// case htdfservice.MsgAdd:
	// 	return htdfservice.HandleMsgAdd(ctx, app.bankKeeper, msg)
	// default:
	// 	fmt.Printf("msgType error|mstType=%v\n", msg.Type())
	// 	sendTxResp.ErrCode = apptypes.ErrCode_Param
	// 	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
	// }

	return sdk.Result{}
}

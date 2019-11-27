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

	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/x/auth"

	//"github.com/orientwalt/htdf/x/mint"

	sdk "github.com/orientwalt/htdf/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"

	v0 "github.com/orientwalt/htdf/app/v0"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
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
	engine.Add(v0.NewProtocolV0(1, logger, protocolKeeper, app.invCheckPeriod, &appPrometheusConfig))
	//engine.Add(v2.NewProtocolV1(2, ...))
	fmt.Print("KeyMain----->	", app.GetKVStore(protocol.KeyMain), "\n")
	loaded, current := engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())
	engine.GetCurrentProtocol().InitMetrics(app.cms)
	fmt.Print("keystorage----->	", app.GetKVStore(protocol.KeyStorage), "\n")
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

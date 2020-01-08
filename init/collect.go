package init

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/orientwalt/htdf/app"
	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/server"
	"github.com/orientwalt/htdf/x/auth"
)

const (
	flagGenTxDir = "gentx-dir"
)

type initConfig struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey crypto.PubKey
}

// nolint
func CollectGenTxsCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(111111111111111111)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			name := viper.GetString(client.FlagName)
			nodeID, valPubKey, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}
			fmt.Println(2222222222222222222)
			genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
			if err != nil {
				return err
			}
			fmt.Println(3333333333333333333)
			genTxsDir := viper.GetString(flagGenTxDir)
			if genTxsDir == "" {
				genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
			}
			fmt.Println(44444444444444444)
			toPrint := newPrintInfo(config.Moniker, genDoc.ChainID, nodeID, genTxsDir, json.RawMessage(""))
			initCfg := newInitConfig(genDoc.ChainID, genTxsDir, name, nodeID, valPubKey)
			fmt.Println(55555555555555555)
			appMessage, err := genAppStateFromConfig(cdc, config, initCfg, genDoc)
			if err != nil {
				return err
			}
			fmt.Println(6666666666666666)
			toPrint.AppMessage = appMessage

			// print out some key information
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagGenTxDir, "",
		"override default \"gentx\" directory from which collect and execute "+
			"genesis transactions; default [--home]/config/gentx/")
	return cmd
}

func genAppStateFromConfig(
	cdc *codec.Codec, config *cfg.Config, initCfg initConfig, genDoc types.GenesisDoc,
) (appState json.RawMessage, err error) {

	genFile := config.GenesisFile()
	var (
		appGenTxs       []auth.StdTx
		persistentPeers string
		genTxs          []json.RawMessage
		jsonRawTx       json.RawMessage
	)
	fmt.Println(77777777777)
	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err = v0.CollectStdTxs(
		cdc, config.Moniker, initCfg.GenTxsDir, genDoc,
	)
	if err != nil {
		return
	}
	fmt.Println(8888888888888)
	genTxs = make([]json.RawMessage, len(appGenTxs))
	config.P2P.PersistentPeers = persistentPeers
	fmt.Println(9999999999999)
	for i, stdTx := range appGenTxs {
		jsonRawTx, err = cdc.MarshalJSON(stdTx)
		if err != nil {
			return
		}
		genTxs[i] = jsonRawTx
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	fmt.Println(0000000000000)
	appState, err = v0.HtdfAppGenStateJSON(cdc, genDoc, genTxs)
	if err != nil {
		return
	}
	fmt.Println("AAAAAAAAAAAA")
	err = ExportGenesisFile(genFile, initCfg.ChainID, nil, appState)
	return
}

// junying-todo-20190517
// remove persistentPeers
func genAppStateFromConfigEx(
	cdc *codec.Codec, config *cfg.Config, initCfg initConfig, genDoc types.GenesisDoc,
) (appState json.RawMessage, err error) {

	genFile := config.GenesisFile()
	var (
		appGenTxs       []auth.StdTx
		persistentPeers string
		genTxs          []json.RawMessage
		jsonRawTx       json.RawMessage
	)

	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err = v0.CollectStdTxsEx(
		cdc, config.Moniker, initCfg.GenTxsDir, genDoc,
	)
	if err != nil {
		return
	}

	genTxs = make([]json.RawMessage, len(appGenTxs))
	config.P2P.PersistentPeers = persistentPeers

	for i, stdTx := range appGenTxs {
		jsonRawTx, err = cdc.MarshalJSON(stdTx)
		if err != nil {
			return
		}
		genTxs[i] = jsonRawTx
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	appState, err = v0.HtdfAppGenStateJSON(cdc, genDoc, genTxs)
	if err != nil {
		return
	}

	err = ExportGenesisFile(genFile, initCfg.ChainID, nil, appState)
	return
}

func newInitConfig(chainID, genTxsDir, name, nodeID string,
	valPubKey crypto.PubKey) initConfig {

	return initConfig{
		ChainID:   chainID,
		GenTxsDir: genTxsDir,
		Name:      name,
		NodeID:    nodeID,
		ValPubKey: valPubKey,
	}
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

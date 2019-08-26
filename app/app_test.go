package app

import (
	"os"
	"testing"

	"github.com/orientwalt/htdf/x/bank"

	"github.com/orientwalt/tendermint/libs/db"
	"github.com/orientwalt/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/x/auth"
	distr "github.com/orientwalt/htdf/x/distribution"
	"github.com/orientwalt/htdf/x/gov"
	"github.com/orientwalt/htdf/x/mint"
	"github.com/orientwalt/htdf/x/slashing"
	"github.com/orientwalt/htdf/x/staking"

	abci "github.com/orientwalt/tendermint/abci/types"
)

func setGenesis(gapp *HtdfServiceApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	genesisState := NewGenesisState(
		genaccs,
		auth.DefaultGenesisState(),
		bank.DefaultGenesisState(),
		staking.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	gapp.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	gapp.Commit()

	return nil
}

func TestGaiadExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewHtdfServiceApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewHtdfServiceApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

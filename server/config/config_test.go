package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	fmt.Printf("minGasPrices=%v\n", cfg.GetMinGasPrices())
	require.True(t, cfg.GetMinGasPrices().IsZero())
}

func TestSetMinimumFees(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinGasPrices(sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)})
	require.Equal(t, "5.000000000000000000foo", cfg.MinGasPrices)
}

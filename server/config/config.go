package config

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
)

//
const (
	DefaultMinGasPrices     = "20.0satoshi"
	ValueSecurityLevel_High = "high"
	ValueSecurityLevel_Low  = "low"

	ValueDebugApi_On = "ON"
)

var (
	// Api Security Level
	//    high : disable operate type API, like new account, send tx ,and so on; only query type API is enable
	//    low  : enable all API
	//    high(level) is default;
	ApiSecurityLevel string
)

func init() {
	ApiSecurityLevel = ValueSecurityLevel_High

}

// BaseConfig defines the server's basic configuration
type BaseConfig struct {
	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. A transaction's fees must meet the minimum of any denomination
	// specified in this config (e.g. 0.01photino;0.0001stake).
	MinGasPrices string `mapstructure:"minimum-gas-prices"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`
}

// SetMinGasPrices sets the validator's minimum gas prices.
func (c *Config) SetMinGasPrices(gasPrices sdk.DecCoins) {
	c.MinGasPrices = gasPrices.String()
}

// GetMinGasPrices returns the validator's minimum gas prices based on the set
// configuration.
func (c *Config) GetMinGasPrices() sdk.DecCoins {
	gasPrices, err := sdk.ParseDecCoins(c.MinGasPrices)
	if err != nil {
		panic(fmt.Sprintf("invalid minimum gas prices: %v", err))
	}

	return gasPrices
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig{
			MinGasPrices: DefaultMinGasPrices,
		},
	}
}

package params

import "fmt"
import sdk "github.com/orientwalt/htdf/types"

// evm gas estimation
const (
	TxGas                 uint64 = 30000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	TxGasContractCreation uint64 = 60000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.

	DefaultMinGasPriceUint = 100
)

//
var DefaultMinGasPriceStr = fmt.Sprintf("%d%s", DefaultMinGasPriceUint, sdk.DefaultDenom)

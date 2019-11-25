package params

// evm gas estimation
const (
	DefaultMsgGas                 uint64 = 30000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	DefaultMsgGasContractCreation uint64 = 60000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.

	DefaultMinGasPrice = 100 //unit: satoshi

	DefaultTxGas = DefaultMsgGas
)

package params

// evm gas estimation
const (
	TxGas                 uint64 = 30000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	TxGasContractCreation uint64 = 60000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.
)

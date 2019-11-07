package htdfservice

import (
	"errors"
	"fmt"
	"math/big"

	ethcore "github.com/ethereum/go-ethereum/core"
	evmstate "github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"

	apptypes "github.com/orientwalt/htdf/types"
)

//
type StateTransition struct {
	gpGasLimit *ethcore.GasPool
	msg        MsgSendFrom
	gas        uint64   //unit: gallon
	gasPrice   *big.Int //unit: satoshi/gallon
	initialGas uint64
	value      *big.Int
	data       []byte
	stateDB    vm.StateDB
	evm        *vm.EVM
}

func NewStateTransition(evm *vm.EVM, msg MsgSendFrom, stateDB *evmstate.CommitStateDB) *StateTransition {
	return &StateTransition{
		gpGasLimit: new(ethcore.GasPool).AddGas(msg.GasLimit),
		evm:        evm,
		stateDB:    stateDB,
		msg:        msg,
		gasPrice:   big.NewInt(int64(msg.GasPrice)),
	}
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return vm.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	st.gas = st.msg.GasLimit
	st.initialGas = st.gas
	fmt.Printf("msgGas=%d\n", st.initialGas)

	eaSender := apptypes.ToEthAddress(st.msg.From)

	msgGasVal := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.GasLimit), st.gasPrice)
	fmt.Printf("msgGasVal=%s\n", msgGasVal.String())

	if st.stateDB.GetBalance(eaSender).Cmp(msgGasVal) < 0 {
		return errors.New("insufficient balance for gas")
	}

	// try call subGas method, to check gas limit
	if err := st.gpGasLimit.SubGas(st.msg.GasLimit); err != nil {
		fmt.Printf("SubGas error|err=%s\n", err)
		return err
	}

	st.stateDB.SubBalance(eaSender, msgGasVal)
	return nil
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.stateDB.GetRefund() {
		refund = st.stateDB.GetRefund()
	}

	st.gas += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	eaSender := apptypes.ToEthAddress(st.msg.From)

	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.stateDB.AddBalance(eaSender, remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gpGasLimit.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

func (st *StateTransition) tokenUsed() uint64 {
	return new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice).Uint64()
}

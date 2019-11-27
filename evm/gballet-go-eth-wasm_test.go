package evm

import (

	//"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	infra "github.com/orientwalt/htdf/evm/test_infra"
	"github.com/orientwalt/htdf/evm/vm"
	"github.com/orientwalt/htdf/params"

	"math/big"

	//"os"
	"testing"
)

// TODO:current test code , is base go-ethereum V1.8.0
//	when this evm package is stable ,need to update to new version, like  V1.8.23

func TestGballetGoEthWasm(t *testing.T) {
	//abiFileName := "./testdata/coin_sol_Coin.abi"
	//binFileName := "./testdata/test.wasm"
	binFileName := "../tests/evm/wasm/add-ex-main.wasm"
	//binFileName := "./testdata/add.wasm"
	data := infra.LoadRaw(binFileName)

	//dataPath := "/tmp/htdfTmpTestData_gballet-eth-wasm"
	//os.Remove(dataPath)
	//mdb, err := ethdb.NewLDBDatabase(dataPath, 100, 100)
	//infra.Must(err)
	//fmt.Println("mdb=%v\n", mdb)

	//db := state.NewDatabase(mdb)
	//root := common.Hash{}
	db := ethdb.NewMemDatabase()
	//infra.Must(err)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(db))
	infra.Must(err)

	//set balance
	statedb.GetOrNewStateObject(infra.FromAddress)
	statedb.GetOrNewStateObject(infra.ToAddress)
	statedb.AddBalance(infra.FromAddress, big.NewInt(1e18))
	testBalance := statedb.GetBalance(infra.FromAddress)
	fmt.Println("init testBalance =", testBalance)
	infra.Must(err)

	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	fmt.Printf("statedb=%v|vmconfig=%v\n", statedb, vmConfig)

	//var vmTest tests.VMTest
	//
	//vmTest.NewEVM(statedb,vmConfig)

	initialCall := true
	canTransfer := func(db vm.StateDB, address common.Address, amount *big.Int) bool {
		if initialCall {
			initialCall = false
			return true
		}
		return core.CanTransfer(db, address, amount)
	}
	transfer := func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {}

	cc := infra.ChainContext{}
	header := cc.GetHeader(infra.TestHash, 0)

	context := vm.Context{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		GetHash:     infra.VmTestBlockHash,
		Origin:      infra.FromAddress,
		Coinbase:    infra.FromAddress,
		BlockNumber: header.Number,
		Time:        big.NewInt(int64(header.Time)),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    big.NewInt(1000),
	}
	vmConfig.NoRecursion = true

	evm := vm.NewEVM(context, statedb, params.MainnetChainConfig, vmConfig)

	fmt.Printf("evm=%v\n", evm)

	contractRef := vm.AccountRef(infra.FromAddress)
	contractCode, contractAddr, gasLeftover, vmerr := evm.Create(contractRef, data, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	infra.Must(vmerr)
	fmt.Printf("getcode:%x\n%x\n", contractCode, statedb.GetCode(contractAddr))
	fmt.Printf("gasLeftover=%v\n", gasLeftover)
	fmt.Printf("after create|contractAddr=%s\n", contractAddr.String())

	statedb.SetBalance(infra.FromAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance = statedb.GetBalance(infra.FromAddress)
	fmt.Println("after create contract, testBalance =", testBalance)
	//abiObj := infra.LoadAbi(abiFileName)
	//
	//input, err := abiObj.Pack("minter")
	//infra.Must(err)

	input := []byte{0x01, 0x00, 0x00, 0x01}

	outputs, gasLeftover, vmerr := evm.Call(contractRef, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	infra.Must(vmerr)

	fmt.Printf("minter is %x\n", common.BytesToAddress(outputs))
	fmt.Printf("call address %x\n", contractRef)
	//
	//sender := common.BytesToAddress(outputs)
	//
	//if !bytes.Equal(sender.Bytes(), infra.FromAddress.Bytes()) {
	//	fmt.Println("caller are not equal to minter!!")
	//	os.Exit(-1)
	//}
	//
	//senderAcc := vm.AccountRef(sender)
	//
	//input, err = abiObj.Pack("mint", sender, big.NewInt(1000000))
	//infra.Must(err)
	//outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//
	//// get balance
	//input, err = abiObj.Pack("balances", sender)
	//infra.Must(err)
	//outputs, gasLeftover, vmerr = evm.Call(contractRef, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//fmt.Printf("contract balance, after mint|minterAddress=%s|Balance=%x\n", sender.String(), outputs)
	//
	//statedb.SetBalance(infra.FromAddress, big.NewInt(0).SetUint64(gasLeftover))
	//testBalance = evm.StateDB.GetBalance(infra.FromAddress)
	//
	//input, err = abiObj.Pack("send", infra.ToAddress, big.NewInt(11))
	//outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//
	////send
	//input, err = abiObj.Pack("send", infra.ToAddress, big.NewInt(19))
	//infra.Must(err)
	//outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//
	//// get balance
	//input, err = abiObj.Pack("balances", infra.ToAddress)
	//infra.Must(err)
	//outputs, gasLeftover, vmerr = evm.Call(contractRef, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//fmt.Printf("contract balance, after send|toAddress=%s|Balance=%x\n", infra.ToAddress.String(), outputs)
	//
	//// get balance
	//input, err = abiObj.Pack("balances", sender)
	//infra.Must(err)
	//outputs, gasLeftover, vmerr = evm.Call(contractRef, contractAddr, input, statedb.GetBalance(infra.FromAddress).Uint64(), big.NewInt(0))
	//infra.Must(vmerr)
	//fmt.Printf("contract balance, after send|minterAddress=%s|Balance=%x\n", sender.String(), outputs)

}

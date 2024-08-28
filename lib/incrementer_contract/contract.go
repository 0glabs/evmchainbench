// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package incrementer_contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IncrementerContractMetaData contains all meta data concerning the IncrementerContract contract.
var IncrementerContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040525f5f553480156011575f5ffd5b5061012f8061001f5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806306661abd146034578063d09de08a14604e575b5f5ffd5b603a6056565b604051604591906089565b60405180910390f35b6054605b565b005b5f5481565b60015f5f828254606a919060cd565b92505081905550565b5f819050919050565b6083816073565b82525050565b5f602082019050609a5f830184607c565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f60d5826073565b915060de836073565b925082820190508082111560f35760f260a0565b5b9291505056fea2646970667358221220a1970bdec0c2b49c3c30fdbaed20c0502fd1556f2399533d35a29502cb31674064736f6c634300081b0033",
}

// IncrementerContractABI is the input ABI used to generate the binding from.
// Deprecated: Use IncrementerContractMetaData.ABI instead.
var IncrementerContractABI = IncrementerContractMetaData.ABI

// IncrementerContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IncrementerContractMetaData.Bin instead.
var IncrementerContractBin = IncrementerContractMetaData.Bin

// DeployIncrementerContract deploys a new Ethereum contract, binding an instance of IncrementerContract to it.
func DeployIncrementerContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IncrementerContract, error) {
	parsed, err := IncrementerContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IncrementerContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IncrementerContract{IncrementerContractCaller: IncrementerContractCaller{contract: contract}, IncrementerContractTransactor: IncrementerContractTransactor{contract: contract}, IncrementerContractFilterer: IncrementerContractFilterer{contract: contract}}, nil
}

// IncrementerContract is an auto generated Go binding around an Ethereum contract.
type IncrementerContract struct {
	IncrementerContractCaller     // Read-only binding to the contract
	IncrementerContractTransactor // Write-only binding to the contract
	IncrementerContractFilterer   // Log filterer for contract events
}

// IncrementerContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type IncrementerContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncrementerContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IncrementerContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncrementerContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IncrementerContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncrementerContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IncrementerContractSession struct {
	Contract     *IncrementerContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IncrementerContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IncrementerContractCallerSession struct {
	Contract *IncrementerContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// IncrementerContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IncrementerContractTransactorSession struct {
	Contract     *IncrementerContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// IncrementerContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type IncrementerContractRaw struct {
	Contract *IncrementerContract // Generic contract binding to access the raw methods on
}

// IncrementerContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IncrementerContractCallerRaw struct {
	Contract *IncrementerContractCaller // Generic read-only contract binding to access the raw methods on
}

// IncrementerContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IncrementerContractTransactorRaw struct {
	Contract *IncrementerContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIncrementerContract creates a new instance of IncrementerContract, bound to a specific deployed contract.
func NewIncrementerContract(address common.Address, backend bind.ContractBackend) (*IncrementerContract, error) {
	contract, err := bindIncrementerContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IncrementerContract{IncrementerContractCaller: IncrementerContractCaller{contract: contract}, IncrementerContractTransactor: IncrementerContractTransactor{contract: contract}, IncrementerContractFilterer: IncrementerContractFilterer{contract: contract}}, nil
}

// NewIncrementerContractCaller creates a new read-only instance of IncrementerContract, bound to a specific deployed contract.
func NewIncrementerContractCaller(address common.Address, caller bind.ContractCaller) (*IncrementerContractCaller, error) {
	contract, err := bindIncrementerContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IncrementerContractCaller{contract: contract}, nil
}

// NewIncrementerContractTransactor creates a new write-only instance of IncrementerContract, bound to a specific deployed contract.
func NewIncrementerContractTransactor(address common.Address, transactor bind.ContractTransactor) (*IncrementerContractTransactor, error) {
	contract, err := bindIncrementerContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IncrementerContractTransactor{contract: contract}, nil
}

// NewIncrementerContractFilterer creates a new log filterer instance of IncrementerContract, bound to a specific deployed contract.
func NewIncrementerContractFilterer(address common.Address, filterer bind.ContractFilterer) (*IncrementerContractFilterer, error) {
	contract, err := bindIncrementerContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IncrementerContractFilterer{contract: contract}, nil
}

// bindIncrementerContract binds a generic wrapper to an already deployed contract.
func bindIncrementerContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IncrementerContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IncrementerContract *IncrementerContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IncrementerContract.Contract.IncrementerContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IncrementerContract *IncrementerContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IncrementerContract.Contract.IncrementerContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IncrementerContract *IncrementerContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IncrementerContract.Contract.IncrementerContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IncrementerContract *IncrementerContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IncrementerContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IncrementerContract *IncrementerContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IncrementerContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IncrementerContract *IncrementerContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IncrementerContract.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_IncrementerContract *IncrementerContractCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IncrementerContract.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_IncrementerContract *IncrementerContractSession) Count() (*big.Int, error) {
	return _IncrementerContract.Contract.Count(&_IncrementerContract.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_IncrementerContract *IncrementerContractCallerSession) Count() (*big.Int, error) {
	return _IncrementerContract.Contract.Count(&_IncrementerContract.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_IncrementerContract *IncrementerContractTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IncrementerContract.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_IncrementerContract *IncrementerContractSession) Increment() (*types.Transaction, error) {
	return _IncrementerContract.Contract.Increment(&_IncrementerContract.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_IncrementerContract *IncrementerContractTransactorSession) Increment() (*types.Transaction, error) {
	return _IncrementerContract.Contract.Increment(&_IncrementerContract.TransactOpts)
}

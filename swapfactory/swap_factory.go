// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swapfactory

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
)

// SwapFactorySwap is an auto generated low-level Go binding around an user-defined struct.
type SwapFactorySwap struct {
	Owner        common.Address
	Claimer      common.Address
	PubKeyClaim  [32]byte
	PubKeyRefund [32]byte
	Timeout0     *big.Int
	Timeout1     *big.Int
	Value        *big.Int
	Nonce        *big.Int
}

// Secp256k1MetaData contains all meta data concerning the Secp256k1 contract.
var Secp256k1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"scalar\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"qKeccak\",\"type\":\"uint256\"}],\"name\":\"mulVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b32d1b4f": "mulVerify(uint256,uint256)",
	},
	Bin: "0x61018861003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c8063b32d1b4f1461003a575b600080fd5b61004d610048366004610130565b610061565b604051901515815260200160405180910390f35b600080600181601b7f79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179870014551231950b75fc4402da1732fc9bebe197f79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179889096040805160008152602081018083529590955260ff909316928401929092526060830152608082015260a0016020604051602081039080840390855afa15801561010e573d6000803e3d6000fd5b5050604051601f1901516001600160a01b038581169116149250505092915050565b6000806040838503121561014357600080fd5b5050803592602090910135915056fea264697066735822122018fbbf5a3edd6fe1ca293ad0ef235a1f40ed6e60a36c7e039f7ee76db17f5eb064736f6c634300080e0033",
}

// Secp256k1ABI is the input ABI used to generate the binding from.
// Deprecated: Use Secp256k1MetaData.ABI instead.
var Secp256k1ABI = Secp256k1MetaData.ABI

// Deprecated: Use Secp256k1MetaData.Sigs instead.
// Secp256k1FuncSigs maps the 4-byte function signature to its string representation.
var Secp256k1FuncSigs = Secp256k1MetaData.Sigs

// Secp256k1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Secp256k1MetaData.Bin instead.
var Secp256k1Bin = Secp256k1MetaData.Bin

// DeploySecp256k1 deploys a new Ethereum contract, binding an instance of Secp256k1 to it.
func DeploySecp256k1(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Secp256k1, error) {
	parsed, err := Secp256k1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Secp256k1Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Secp256k1{Secp256k1Caller: Secp256k1Caller{contract: contract}, Secp256k1Transactor: Secp256k1Transactor{contract: contract}, Secp256k1Filterer: Secp256k1Filterer{contract: contract}}, nil
}

// Secp256k1 is an auto generated Go binding around an Ethereum contract.
type Secp256k1 struct {
	Secp256k1Caller     // Read-only binding to the contract
	Secp256k1Transactor // Write-only binding to the contract
	Secp256k1Filterer   // Log filterer for contract events
}

// Secp256k1Caller is an auto generated read-only Go binding around an Ethereum contract.
type Secp256k1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Secp256k1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Secp256k1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Secp256k1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Secp256k1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Secp256k1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Secp256k1Session struct {
	Contract     *Secp256k1        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Secp256k1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Secp256k1CallerSession struct {
	Contract *Secp256k1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Secp256k1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Secp256k1TransactorSession struct {
	Contract     *Secp256k1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Secp256k1Raw is an auto generated low-level Go binding around an Ethereum contract.
type Secp256k1Raw struct {
	Contract *Secp256k1 // Generic contract binding to access the raw methods on
}

// Secp256k1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Secp256k1CallerRaw struct {
	Contract *Secp256k1Caller // Generic read-only contract binding to access the raw methods on
}

// Secp256k1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Secp256k1TransactorRaw struct {
	Contract *Secp256k1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSecp256k1 creates a new instance of Secp256k1, bound to a specific deployed contract.
func NewSecp256k1(address common.Address, backend bind.ContractBackend) (*Secp256k1, error) {
	contract, err := bindSecp256k1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Secp256k1{Secp256k1Caller: Secp256k1Caller{contract: contract}, Secp256k1Transactor: Secp256k1Transactor{contract: contract}, Secp256k1Filterer: Secp256k1Filterer{contract: contract}}, nil
}

// NewSecp256k1Caller creates a new read-only instance of Secp256k1, bound to a specific deployed contract.
func NewSecp256k1Caller(address common.Address, caller bind.ContractCaller) (*Secp256k1Caller, error) {
	contract, err := bindSecp256k1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Secp256k1Caller{contract: contract}, nil
}

// NewSecp256k1Transactor creates a new write-only instance of Secp256k1, bound to a specific deployed contract.
func NewSecp256k1Transactor(address common.Address, transactor bind.ContractTransactor) (*Secp256k1Transactor, error) {
	contract, err := bindSecp256k1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Secp256k1Transactor{contract: contract}, nil
}

// NewSecp256k1Filterer creates a new log filterer instance of Secp256k1, bound to a specific deployed contract.
func NewSecp256k1Filterer(address common.Address, filterer bind.ContractFilterer) (*Secp256k1Filterer, error) {
	contract, err := bindSecp256k1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Secp256k1Filterer{contract: contract}, nil
}

// bindSecp256k1 binds a generic wrapper to an already deployed contract.
func bindSecp256k1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Secp256k1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Secp256k1 *Secp256k1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Secp256k1.Contract.Secp256k1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Secp256k1 *Secp256k1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Secp256k1.Contract.Secp256k1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Secp256k1 *Secp256k1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Secp256k1.Contract.Secp256k1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Secp256k1 *Secp256k1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Secp256k1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Secp256k1 *Secp256k1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Secp256k1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Secp256k1 *Secp256k1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Secp256k1.Contract.contract.Transact(opts, method, params...)
}

// MulVerify is a free data retrieval call binding the contract method 0xb32d1b4f.
//
// Solidity: function mulVerify(uint256 scalar, uint256 qKeccak) pure returns(bool)
func (_Secp256k1 *Secp256k1Caller) MulVerify(opts *bind.CallOpts, scalar *big.Int, qKeccak *big.Int) (bool, error) {
	var out []interface{}
	err := _Secp256k1.contract.Call(opts, &out, "mulVerify", scalar, qKeccak)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MulVerify is a free data retrieval call binding the contract method 0xb32d1b4f.
//
// Solidity: function mulVerify(uint256 scalar, uint256 qKeccak) pure returns(bool)
func (_Secp256k1 *Secp256k1Session) MulVerify(scalar *big.Int, qKeccak *big.Int) (bool, error) {
	return _Secp256k1.Contract.MulVerify(&_Secp256k1.CallOpts, scalar, qKeccak)
}

// MulVerify is a free data retrieval call binding the contract method 0xb32d1b4f.
//
// Solidity: function mulVerify(uint256 scalar, uint256 qKeccak) pure returns(bool)
func (_Secp256k1 *Secp256k1CallerSession) MulVerify(scalar *big.Int, qKeccak *big.Int) (bool, error) {
	return _Secp256k1.Contract.MulVerify(&_Secp256k1.CallOpts, scalar, qKeccak)
}

// SwapFactoryMetaData contains all meta data concerning the SwapFactory contract.
var SwapFactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"claimKey\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"refundKey\",\"type\":\"bytes32\"}],\"name\":\"New\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"swapID\",\"type\":\"bytes32\"}],\"name\":\"Ready\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"claimer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyRefund\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeout_0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeout_1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structSwapFactory.Swap\",\"name\":\"_swap\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"is_ready\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_pubKeyClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_pubKeyRefund\",\"type\":\"bytes32\"},{\"internalType\":\"addresspayable\",\"name\":\"_claimer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_timeoutDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"new_swap\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"claimer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyRefund\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeout_0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeout_1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structSwapFactory.Swap\",\"name\":\"_swap\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"claimer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"pubKeyRefund\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeout_0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeout_1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structSwapFactory.Swap\",\"name\":\"_swap\",\"type\":\"tuple\"}],\"name\":\"set_ready\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swaps\",\"outputs\":[{\"internalType\":\"enumSwapFactory.Stage\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7069c7f3": "claim((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256),bytes32)",
		"268a3bd4": "is_ready(bytes32)",
		"d749b6c4": "new_swap(bytes32,bytes32,address,uint256,uint256)",
		"262cd8da": "refund((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256),bytes32)",
		"3e7a7b55": "set_ready((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256))",
		"eb84e7f2": "swaps(bytes32)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610c7e806100206000396000f3fe6080604052600436106100555760003560e01c8063262cd8da1461005a578063268a3bd41461007c5780633e7a7b55146100b15780637069c7f3146100d1578063d749b6c4146100f1578063eb84e7f214610112575b600080fd5b34801561006657600080fd5b5061007a610075366004610a85565b61014f565b005b34801561008857600080fd5b5061009c610097366004610ab2565b6103bc565b60405190151581526020015b60405180910390f35b3480156100bd57600080fd5b5061007a6100cc366004610acb565b6103ea565b3480156100dd57600080fd5b5061007a6100ec366004610a85565b610540565b6101046100ff366004610aef565b610794565b6040519081526020016100a8565b34801561011e57600080fd5b5061014261012d366004610ab2565b60006020819052908152604090205460ff1681565b6040516100a89190610b4c565b6000826040516020016101629190610b74565b60408051601f19818403018152918152815160209283012060008181529283905291205490915060ff1660038160038111156101a0576101a0610b36565b141580156101c0575060008160038111156101bd576101bd610b36565b14155b61020d5760405162461bcd60e51b81526020600482015260196024820152781cddd85c081a5cc8185b1c9958591e4818dbdb5c1b195d1959603a1b60448201526064015b60405180910390fd5b83516001600160a01b031633146102765760405162461bcd60e51b815260206004820152602760248201527f726566756e64206d7573742062652063616c6c65642062792074686520737761604482015266381037bbb732b960c91b6064820152608401610204565b8360a00151421015806102a857508360800151421080156102a8575060028160038111156102a6576102a6610b36565b145b61031a5760405162461bcd60e51b815260206004820152603f60248201527f697427732074686520636f756e74657270617274792773207475726e2c20756e60448201527f61626c6520746f20726566756e642c2074727920616761696e206c61746572006064820152608401610204565b6103288385606001516108d3565b60408051838152602081018590527e7c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f910160405180910390a1835160c08501516040516001600160a01b039092169181156108fc0291906000818181858888f1935050505015801561039d573d6000803e3d6000fd5b50506000908152602081905260409020805460ff191660031790555050565b6000600260008381526020819052604090205460ff1660038111156103e3576103e3610b36565b1492915050565b6000816040516020016103fd9190610b74565b60408051601f1981840301815291905280516020909101209050600160008281526020819052604090205460ff16600381111561043c5761043c610b36565b146104895760405162461bcd60e51b815260206004820152601c60248201527f73776170206973206e6f7420696e2050454e44494e47207374617465000000006044820152606401610204565b81516001600160a01b031633146104f15760405162461bcd60e51b815260206004820152602660248201527f6f6e6c79207468652073776170206f776e65722063616e2063616c6c207365746044820152655f726561647960d01b6064820152608401610204565b60008181526020818152604091829020805460ff1916600217905590518281527f5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f910160405180910390a15050565b6000826040516020016105539190610b74565b60408051601f19818403018152918152815160209283012060008181529283905291205490915060ff16600381600381111561059157610591610b36565b141580156105b1575060008160038111156105ae576105ae610b36565b14155b6105f95760405162461bcd60e51b81526020600482015260196024820152781cddd85c081a5cc8185b1c9958591e4818dbdb5c1b195d1959603a1b6044820152606401610204565b83602001516001600160a01b0316336001600160a01b03161461065e5760405162461bcd60e51b815260206004820152601760248201527f6f6e6c7920636c61696d65722063616e20636c61696d210000000000000000006044820152606401610204565b8360800151421015806106825750600281600381111561068057610680610b36565b145b6106c45760405162461bcd60e51b8152602060048201526013602482015272746f6f206561726c7920746f20636c61696d2160681b6044820152606401610204565b8360a00151421061070c5760405162461bcd60e51b8152602060048201526012602482015271746f6f206c61746520746f20636c61696d2160701b6044820152606401610204565b61071a8385604001516108d3565b60408051838152602081018590527f38d6042dbdae8e73a7f6afbabd3fbe0873f9f5ed3cd71294591c3908c2e65fee910160405180910390a183602001516001600160a01b03166108fc8560c001519081150290604051600060405180830381858888f1935050505015801561039d573d6000803e3d6000fd5b604080516101008101825260006080820181905260a0820181905260c0820181905260e082018190523382526001600160a01b0386166020830152918101879052606081018690526107e68442610bef565b60808201526107f6846002610c07565b6108009042610bef565b60a08201523460c082015260e08101839052604051600090610826908390602001610b74565b60408051601f19818403018152919052805160209091012090506000808281526020819052604090205460ff16600381111561086457610864610b36565b1461086e57600080fd5b60408051828152602081018a90529081018890527f0444ef8b0e734743401a0d6e7ce42bd7581151a42627a3949cebe31fd9adf89b9060600160405180910390a16000818152602081905260409020805460ff19166001179055979650505050505050565b60405163b32d1b4f60e01b8152600481018390526024810182905273__$a2f8518bf039aad522402b98cf6b7ac660$__9063b32d1b4f90604401602060405180830381865af415801561092a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061094e9190610c26565b6109b95760405162461bcd60e51b815260206004820152603660248201527f70726f76696465642073656372657420646f6573206e6f74206d6174636820746044820152756865206578706563746564207075626c6963206b657960501b6064820152608401610204565b5050565b80356001600160a01b03811681146109d457600080fd5b919050565b60006101008083850312156109ed57600080fd5b6040519081019067ffffffffffffffff82118183101715610a1e57634e487b7160e01b600052604160045260246000fd5b81604052809250610a2e846109bd565b8152610a3c602085016109bd565b602082015260408401356040820152606084013560608201526080840135608082015260a084013560a082015260c084013560c082015260e084013560e0820152505092915050565b6000806101208385031215610a9957600080fd5b610aa384846109d9565b94610100939093013593505050565b600060208284031215610ac457600080fd5b5035919050565b60006101008284031215610ade57600080fd5b610ae883836109d9565b9392505050565b600080600080600060a08688031215610b0757600080fd5b8535945060208601359350610b1e604087016109bd565b94979396509394606081013594506080013592915050565b634e487b7160e01b600052602160045260246000fd5b6020810160048310610b6e57634e487b7160e01b600052602160045260246000fd5b91905290565b60006101008201905060018060a01b038084511683528060208501511660208401525060408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015160c083015260e083015160e083015292915050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610c0257610c02610bd9565b500190565b6000816000190483118215151615610c2157610c21610bd9565b500290565b600060208284031215610c3857600080fd5b81518015158114610ae857600080fdfea26469706673582212203104aa47f057880b33eb9c0c896dd05cfa5e891f958b1880d81b5e0ffc1a0a8964736f6c634300080e0033",
}

// SwapFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use SwapFactoryMetaData.ABI instead.
var SwapFactoryABI = SwapFactoryMetaData.ABI

// Deprecated: Use SwapFactoryMetaData.Sigs instead.
// SwapFactoryFuncSigs maps the 4-byte function signature to its string representation.
var SwapFactoryFuncSigs = SwapFactoryMetaData.Sigs

// SwapFactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SwapFactoryMetaData.Bin instead.
var SwapFactoryBin = SwapFactoryMetaData.Bin

// DeploySwapFactory deploys a new Ethereum contract, binding an instance of SwapFactory to it.
func DeploySwapFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SwapFactory, error) {
	parsed, err := SwapFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	secp256k1Addr, _, _, _ := DeploySecp256k1(auth, backend)
	SwapFactoryBin = strings.Replace(SwapFactoryBin, "__$a2f8518bf039aad522402b98cf6b7ac660$__", secp256k1Addr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SwapFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapFactory{SwapFactoryCaller: SwapFactoryCaller{contract: contract}, SwapFactoryTransactor: SwapFactoryTransactor{contract: contract}, SwapFactoryFilterer: SwapFactoryFilterer{contract: contract}}, nil
}

// SwapFactory is an auto generated Go binding around an Ethereum contract.
type SwapFactory struct {
	SwapFactoryCaller     // Read-only binding to the contract
	SwapFactoryTransactor // Write-only binding to the contract
	SwapFactoryFilterer   // Log filterer for contract events
}

// SwapFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwapFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapFactorySession struct {
	Contract     *SwapFactory      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapFactoryCallerSession struct {
	Contract *SwapFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SwapFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapFactoryTransactorSession struct {
	Contract     *SwapFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SwapFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwapFactoryRaw struct {
	Contract *SwapFactory // Generic contract binding to access the raw methods on
}

// SwapFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapFactoryCallerRaw struct {
	Contract *SwapFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// SwapFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapFactoryTransactorRaw struct {
	Contract *SwapFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapFactory creates a new instance of SwapFactory, bound to a specific deployed contract.
func NewSwapFactory(address common.Address, backend bind.ContractBackend) (*SwapFactory, error) {
	contract, err := bindSwapFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapFactory{SwapFactoryCaller: SwapFactoryCaller{contract: contract}, SwapFactoryTransactor: SwapFactoryTransactor{contract: contract}, SwapFactoryFilterer: SwapFactoryFilterer{contract: contract}}, nil
}

// NewSwapFactoryCaller creates a new read-only instance of SwapFactory, bound to a specific deployed contract.
func NewSwapFactoryCaller(address common.Address, caller bind.ContractCaller) (*SwapFactoryCaller, error) {
	contract, err := bindSwapFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapFactoryCaller{contract: contract}, nil
}

// NewSwapFactoryTransactor creates a new write-only instance of SwapFactory, bound to a specific deployed contract.
func NewSwapFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SwapFactoryTransactor, error) {
	contract, err := bindSwapFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapFactoryTransactor{contract: contract}, nil
}

// NewSwapFactoryFilterer creates a new log filterer instance of SwapFactory, bound to a specific deployed contract.
func NewSwapFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SwapFactoryFilterer, error) {
	contract, err := bindSwapFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapFactoryFilterer{contract: contract}, nil
}

// bindSwapFactory binds a generic wrapper to an already deployed contract.
func bindSwapFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapFactory *SwapFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapFactory.Contract.SwapFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapFactory *SwapFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapFactory.Contract.SwapFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapFactory *SwapFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapFactory.Contract.SwapFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapFactory *SwapFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapFactory *SwapFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapFactory *SwapFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapFactory.Contract.contract.Transact(opts, method, params...)
}

// IsReady is a free data retrieval call binding the contract method 0x268a3bd4.
//
// Solidity: function is_ready(bytes32 _swapID) view returns(bool)
func (_SwapFactory *SwapFactoryCaller) IsReady(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var out []interface{}
	err := _SwapFactory.contract.Call(opts, &out, "is_ready", _swapID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsReady is a free data retrieval call binding the contract method 0x268a3bd4.
//
// Solidity: function is_ready(bytes32 _swapID) view returns(bool)
func (_SwapFactory *SwapFactorySession) IsReady(_swapID [32]byte) (bool, error) {
	return _SwapFactory.Contract.IsReady(&_SwapFactory.CallOpts, _swapID)
}

// IsReady is a free data retrieval call binding the contract method 0x268a3bd4.
//
// Solidity: function is_ready(bytes32 _swapID) view returns(bool)
func (_SwapFactory *SwapFactoryCallerSession) IsReady(_swapID [32]byte) (bool, error) {
	return _SwapFactory.Contract.IsReady(&_SwapFactory.CallOpts, _swapID)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint8)
func (_SwapFactory *SwapFactoryCaller) Swaps(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var out []interface{}
	err := _SwapFactory.contract.Call(opts, &out, "swaps", arg0)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint8)
func (_SwapFactory *SwapFactorySession) Swaps(arg0 [32]byte) (uint8, error) {
	return _SwapFactory.Contract.Swaps(&_SwapFactory.CallOpts, arg0)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint8)
func (_SwapFactory *SwapFactoryCallerSession) Swaps(arg0 [32]byte) (uint8, error) {
	return _SwapFactory.Contract.Swaps(&_SwapFactory.CallOpts, arg0)
}

// Claim is a paid mutator transaction binding the contract method 0x7069c7f3.
//
// Solidity: function claim((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactoryTransactor) Claim(opts *bind.TransactOpts, _swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.contract.Transact(opts, "claim", _swap, _s)
}

// Claim is a paid mutator transaction binding the contract method 0x7069c7f3.
//
// Solidity: function claim((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactorySession) Claim(_swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.Contract.Claim(&_SwapFactory.TransactOpts, _swap, _s)
}

// Claim is a paid mutator transaction binding the contract method 0x7069c7f3.
//
// Solidity: function claim((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactoryTransactorSession) Claim(_swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.Contract.Claim(&_SwapFactory.TransactOpts, _swap, _s)
}

// NewSwap is a paid mutator transaction binding the contract method 0xd749b6c4.
//
// Solidity: function new_swap(bytes32 _pubKeyClaim, bytes32 _pubKeyRefund, address _claimer, uint256 _timeoutDuration, uint256 _nonce) payable returns(bytes32)
func (_SwapFactory *SwapFactoryTransactor) NewSwap(opts *bind.TransactOpts, _pubKeyClaim [32]byte, _pubKeyRefund [32]byte, _claimer common.Address, _timeoutDuration *big.Int, _nonce *big.Int) (*types.Transaction, error) {
	return _SwapFactory.contract.Transact(opts, "new_swap", _pubKeyClaim, _pubKeyRefund, _claimer, _timeoutDuration, _nonce)
}

// NewSwap is a paid mutator transaction binding the contract method 0xd749b6c4.
//
// Solidity: function new_swap(bytes32 _pubKeyClaim, bytes32 _pubKeyRefund, address _claimer, uint256 _timeoutDuration, uint256 _nonce) payable returns(bytes32)
func (_SwapFactory *SwapFactorySession) NewSwap(_pubKeyClaim [32]byte, _pubKeyRefund [32]byte, _claimer common.Address, _timeoutDuration *big.Int, _nonce *big.Int) (*types.Transaction, error) {
	return _SwapFactory.Contract.NewSwap(&_SwapFactory.TransactOpts, _pubKeyClaim, _pubKeyRefund, _claimer, _timeoutDuration, _nonce)
}

// NewSwap is a paid mutator transaction binding the contract method 0xd749b6c4.
//
// Solidity: function new_swap(bytes32 _pubKeyClaim, bytes32 _pubKeyRefund, address _claimer, uint256 _timeoutDuration, uint256 _nonce) payable returns(bytes32)
func (_SwapFactory *SwapFactoryTransactorSession) NewSwap(_pubKeyClaim [32]byte, _pubKeyRefund [32]byte, _claimer common.Address, _timeoutDuration *big.Int, _nonce *big.Int) (*types.Transaction, error) {
	return _SwapFactory.Contract.NewSwap(&_SwapFactory.TransactOpts, _pubKeyClaim, _pubKeyRefund, _claimer, _timeoutDuration, _nonce)
}

// Refund is a paid mutator transaction binding the contract method 0x262cd8da.
//
// Solidity: function refund((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactoryTransactor) Refund(opts *bind.TransactOpts, _swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.contract.Transact(opts, "refund", _swap, _s)
}

// Refund is a paid mutator transaction binding the contract method 0x262cd8da.
//
// Solidity: function refund((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactorySession) Refund(_swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.Contract.Refund(&_SwapFactory.TransactOpts, _swap, _s)
}

// Refund is a paid mutator transaction binding the contract method 0x262cd8da.
//
// Solidity: function refund((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap, bytes32 _s) returns()
func (_SwapFactory *SwapFactoryTransactorSession) Refund(_swap SwapFactorySwap, _s [32]byte) (*types.Transaction, error) {
	return _SwapFactory.Contract.Refund(&_SwapFactory.TransactOpts, _swap, _s)
}

// SetReady is a paid mutator transaction binding the contract method 0x3e7a7b55.
//
// Solidity: function set_ready((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap) returns()
func (_SwapFactory *SwapFactoryTransactor) SetReady(opts *bind.TransactOpts, _swap SwapFactorySwap) (*types.Transaction, error) {
	return _SwapFactory.contract.Transact(opts, "set_ready", _swap)
}

// SetReady is a paid mutator transaction binding the contract method 0x3e7a7b55.
//
// Solidity: function set_ready((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap) returns()
func (_SwapFactory *SwapFactorySession) SetReady(_swap SwapFactorySwap) (*types.Transaction, error) {
	return _SwapFactory.Contract.SetReady(&_SwapFactory.TransactOpts, _swap)
}

// SetReady is a paid mutator transaction binding the contract method 0x3e7a7b55.
//
// Solidity: function set_ready((address,address,bytes32,bytes32,uint256,uint256,uint256,uint256) _swap) returns()
func (_SwapFactory *SwapFactoryTransactorSession) SetReady(_swap SwapFactorySwap) (*types.Transaction, error) {
	return _SwapFactory.Contract.SetReady(&_SwapFactory.TransactOpts, _swap)
}

// SwapFactoryClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the SwapFactory contract.
type SwapFactoryClaimedIterator struct {
	Event *SwapFactoryClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SwapFactoryClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapFactoryClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SwapFactoryClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SwapFactoryClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapFactoryClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapFactoryClaimed represents a Claimed event raised by the SwapFactory contract.
type SwapFactoryClaimed struct {
	SwapID [32]byte
	S      [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x38d6042dbdae8e73a7f6afbabd3fbe0873f9f5ed3cd71294591c3908c2e65fee.
//
// Solidity: event Claimed(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) FilterClaimed(opts *bind.FilterOpts) (*SwapFactoryClaimedIterator, error) {

	logs, sub, err := _SwapFactory.contract.FilterLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return &SwapFactoryClaimedIterator{contract: _SwapFactory.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x38d6042dbdae8e73a7f6afbabd3fbe0873f9f5ed3cd71294591c3908c2e65fee.
//
// Solidity: event Claimed(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *SwapFactoryClaimed) (event.Subscription, error) {

	logs, sub, err := _SwapFactory.contract.WatchLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapFactoryClaimed)
				if err := _SwapFactory.contract.UnpackLog(event, "Claimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimed is a log parse operation binding the contract event 0x38d6042dbdae8e73a7f6afbabd3fbe0873f9f5ed3cd71294591c3908c2e65fee.
//
// Solidity: event Claimed(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) ParseClaimed(log types.Log) (*SwapFactoryClaimed, error) {
	event := new(SwapFactoryClaimed)
	if err := _SwapFactory.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapFactoryNewIterator is returned from FilterNew and is used to iterate over the raw logs and unpacked data for New events raised by the SwapFactory contract.
type SwapFactoryNewIterator struct {
	Event *SwapFactoryNew // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SwapFactoryNewIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapFactoryNew)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SwapFactoryNew)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SwapFactoryNewIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapFactoryNewIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapFactoryNew represents a New event raised by the SwapFactory contract.
type SwapFactoryNew struct {
	SwapID    [32]byte
	ClaimKey  [32]byte
	RefundKey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNew is a free log retrieval operation binding the contract event 0x0444ef8b0e734743401a0d6e7ce42bd7581151a42627a3949cebe31fd9adf89b.
//
// Solidity: event New(bytes32 swapID, bytes32 claimKey, bytes32 refundKey)
func (_SwapFactory *SwapFactoryFilterer) FilterNew(opts *bind.FilterOpts) (*SwapFactoryNewIterator, error) {

	logs, sub, err := _SwapFactory.contract.FilterLogs(opts, "New")
	if err != nil {
		return nil, err
	}
	return &SwapFactoryNewIterator{contract: _SwapFactory.contract, event: "New", logs: logs, sub: sub}, nil
}

// WatchNew is a free log subscription operation binding the contract event 0x0444ef8b0e734743401a0d6e7ce42bd7581151a42627a3949cebe31fd9adf89b.
//
// Solidity: event New(bytes32 swapID, bytes32 claimKey, bytes32 refundKey)
func (_SwapFactory *SwapFactoryFilterer) WatchNew(opts *bind.WatchOpts, sink chan<- *SwapFactoryNew) (event.Subscription, error) {

	logs, sub, err := _SwapFactory.contract.WatchLogs(opts, "New")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapFactoryNew)
				if err := _SwapFactory.contract.UnpackLog(event, "New", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNew is a log parse operation binding the contract event 0x0444ef8b0e734743401a0d6e7ce42bd7581151a42627a3949cebe31fd9adf89b.
//
// Solidity: event New(bytes32 swapID, bytes32 claimKey, bytes32 refundKey)
func (_SwapFactory *SwapFactoryFilterer) ParseNew(log types.Log) (*SwapFactoryNew, error) {
	event := new(SwapFactoryNew)
	if err := _SwapFactory.contract.UnpackLog(event, "New", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapFactoryReadyIterator is returned from FilterReady and is used to iterate over the raw logs and unpacked data for Ready events raised by the SwapFactory contract.
type SwapFactoryReadyIterator struct {
	Event *SwapFactoryReady // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SwapFactoryReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapFactoryReady)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SwapFactoryReady)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SwapFactoryReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapFactoryReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapFactoryReady represents a Ready event raised by the SwapFactory contract.
type SwapFactoryReady struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReady is a free log retrieval operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 swapID)
func (_SwapFactory *SwapFactoryFilterer) FilterReady(opts *bind.FilterOpts) (*SwapFactoryReadyIterator, error) {

	logs, sub, err := _SwapFactory.contract.FilterLogs(opts, "Ready")
	if err != nil {
		return nil, err
	}
	return &SwapFactoryReadyIterator{contract: _SwapFactory.contract, event: "Ready", logs: logs, sub: sub}, nil
}

// WatchReady is a free log subscription operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 swapID)
func (_SwapFactory *SwapFactoryFilterer) WatchReady(opts *bind.WatchOpts, sink chan<- *SwapFactoryReady) (event.Subscription, error) {

	logs, sub, err := _SwapFactory.contract.WatchLogs(opts, "Ready")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapFactoryReady)
				if err := _SwapFactory.contract.UnpackLog(event, "Ready", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReady is a log parse operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 swapID)
func (_SwapFactory *SwapFactoryFilterer) ParseReady(log types.Log) (*SwapFactoryReady, error) {
	event := new(SwapFactoryReady)
	if err := _SwapFactory.contract.UnpackLog(event, "Ready", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SwapFactoryRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the SwapFactory contract.
type SwapFactoryRefundedIterator struct {
	Event *SwapFactoryRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SwapFactoryRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapFactoryRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SwapFactoryRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SwapFactoryRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapFactoryRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapFactoryRefunded represents a Refunded event raised by the SwapFactory contract.
type SwapFactoryRefunded struct {
	SwapID [32]byte
	S      [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) FilterRefunded(opts *bind.FilterOpts) (*SwapFactoryRefundedIterator, error) {

	logs, sub, err := _SwapFactory.contract.FilterLogs(opts, "Refunded")
	if err != nil {
		return nil, err
	}
	return &SwapFactoryRefundedIterator{contract: _SwapFactory.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *SwapFactoryRefunded) (event.Subscription, error) {

	logs, sub, err := _SwapFactory.contract.WatchLogs(opts, "Refunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapFactoryRefunded)
				if err := _SwapFactory.contract.UnpackLog(event, "Refunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefunded is a log parse operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 swapID, bytes32 s)
func (_SwapFactory *SwapFactoryFilterer) ParseRefunded(log types.Log) (*SwapFactoryRefunded, error) {
	event := new(SwapFactoryRefunded)
	if err := _SwapFactory.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

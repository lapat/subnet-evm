// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ StatefulPrecompileConfig = &ContractTestPrecompileConfig{}
	// Singleton StatefulPrecompiledContract for minting native assets by permissioned callers.
	ContractTestPrecompile StatefulPrecompiledContract = createTestFunctionPrecompile(ContractTestPrecompileAddress)

	testFunctionSignature = CalculateFunctionSelector("testFunction(address,uint256)") // address, amount

	ErrCannotTest = errors.New("non-enabled cannot test")

	mintInputLen = common.HashLength + common.HashLength
)

// ContractTestPrecompileConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type ContractTestPrecompileConfig struct {
	AllowListConfig
}

// Address returns the address of the native minter contract.
func (c *ContractTestPrecompileConfig) Address() common.Address {
	return ContractTestPrecompileAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *ContractTestPrecompileConfig) Configure(state StateDB) {
	c.AllowListConfig.Configure(state, ContractTestPrecompileAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the native minter.
func (c *ContractTestPrecompileConfig) Contract() StatefulPrecompiledContract {
	return ContractTestPrecompile
}

// GetContractNativeMinterStatus returns the role of [address] for the minter list.
func GetContractNativeMinterStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, ContractTestPrecompileAddress, address)
}

// SetContractNativeMinterStatus sets the permissions of [address] to [role] for the
// minter list. assumes [role] has already been verified as valid.
func SetContractNativeMinterStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, ContractTestPrecompileAddress, address, role)
}

// PackMintInput packs [address] and [amount] into the appropriate arguments for minting operation.
func PackMintInput(address common.Address, amount *big.Int) ([]byte, error) {
	// function selector (4 bytes) + input(hash for address + hash for amount)
	fullLen := selectorLen + mintInputLen
	input := make([]byte, fullLen)
	copy(input[:selectorLen], testFunctionSignature)
	copy(input[selectorLen:selectorLen+common.HashLength], address.Hash().Bytes())
	amount.FillBytes(input[selectorLen+common.HashLength : fullLen])
	return input, nil
}

// UnpackMintInput attempts to unpack [input] into the arguments to the mint precompile
// assumes that [input] does not include selector (omits first 4 bytes in PackMintInput)
func UnpackMintInput(input []byte) (common.Address, *big.Int, error) {
	if len(input) != mintInputLen {
		return common.Address{}, nil, fmt.Errorf("invalid input length for minting: %d", len(input))
	}
	to := common.BytesToAddress(input[:common.HashLength])
	assetAmount := new(big.Int).SetBytes(input[common.HashLength : common.HashLength+common.HashLength])
	return to, assetAmount, nil
}

// createMintNativeCoin checks if the caller is permissioned for minting operation.
// The execution function parses the [input] into native coin amount and receiver address.
func createTestFunction(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, MintGasCost); err != nil {
		return nil, 0, err
	}

	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}

	to, amount, err := UnpackMintInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, ContractTestPrecompileAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotTest, caller)
	}

	// if there is no address in the state, create one.
	if !stateDB.Exist(to) {
		stateDB.CreateAccount(to)
	}

	stateDB.AddBalance(to, amount)
	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// createNativeMinterPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr] and a native coin minter.
func createTestFunctionPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	setAdmin := newStatefulPrecompileFunction(setAdminSignature, createAllowListRoleSetter(precompileAddr, AllowListAdmin))
	setEnabled := newStatefulPrecompileFunction(setEnabledSignature, createAllowListRoleSetter(precompileAddr, AllowListEnabled))
	setNone := newStatefulPrecompileFunction(setNoneSignature, createAllowListRoleSetter(precompileAddr, AllowListNoRole))
	read := newStatefulPrecompileFunction(readAllowListSignature, createReadAllowList(precompileAddr))

	mint := newStatefulPrecompileFunction(testFunctionSignature, createTestFunction)

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{setAdmin, setEnabled, setNone, read, mint})
	return contract
}

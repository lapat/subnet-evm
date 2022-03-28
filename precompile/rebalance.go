// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	//"errors"
	"fmt"
	"math/big"
	"encoding/json"

	//"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"

	"io/ioutil"
   "log"
  "net/http"

)

// Enum constants for valid AllowListRole
//type AllowListRole common.Hash

var (
	_ StatefulPrecompileConfig = &ContractRebalanceConfig{}
	ContractRebalancePrecompile StatefulPrecompiledContract = createRebalancePrecompile(ContractRebalanceAddress)

	getCoinGeckoSignature = CalculateFunctionSelector("getCoingeckoPrice(address)")

	// Error returned when an invalid write is attempted
	//ErrCannotModifyAllowList = errors.New("non-admin cannot modify allow list")

	//allowListInputLen = common.HashLength
	rebalInputLen = common.HashLength// + common.HashLength

)



type ContractRebalanceConfig struct {
	AllowListConfig
}

// Address returns the address of the native minter contract.
func (c *ContractRebalanceConfig) Address() common.Address {
	return ContractRebalanceAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *ContractRebalanceConfig) Configure(state StateDB) {
	c.AllowListConfig.Configure(state, ContractRebalanceAddress)
}


//func (c *AllowListConfig) Timestamp() *big.Int { return c.BlockTimestamp }


// Configure configures [state] with the desired admins based on [c].
//func (c *ContractRebalanceConfig) Configure(state StateDB) {
	//c.AllowListConfig.Configure(state, ContractRebalanceAddress)
//}

func (c *ContractRebalanceConfig) Contract() StatefulPrecompiledContract {
	return ContractRebalancePrecompile
}



func SetContractRebalanceStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, ContractRebalanceAddress, address, role)
}

func GetContractRebalanceStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, ContractNativeMinterAddress, address)
}


func UnpackRebalanceInput(input []byte) (common.Address, error) {
	if len(input) != rebalInputLen {
		return common.Address{}, fmt.Errorf("invalid input length for : %d", len(input))
	}
	to := common.BytesToAddress(input[:common.HashLength])
	//assetAmount := new(big.Int).SetBytes(input[common.HashLength : common.HashLength+common.HashLength])
	return to,  nil
}


func getCoingeckoPrice(precompileAddr common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {

	log.Printf("createRebalance")

		theAddress, err := UnpackRebalanceInput(input)

		if err != nil {
			return nil, remainingGas, err
		}
		//stringAddress := fmt.Sprintf("%p", theAddress)
		//log.Printf(stringAddress)
		stringAddress := theAddress.String()
		//log.Printf(stringAddress)

		resp, err := http.Get("https://api.coingecko.com/api/v3/simple/token_price/avalanche?contract_addresses="+stringAddress+"&vs_currencies=usd")
		 if err != nil {
				log.Fatalln(err)
		 }
		 //We Read the response body on the line below.
		 body, err := ioutil.ReadAll(resp.Body)
		 if err != nil {
				log.Fatalln(err)
		 }
		 //Convert the body to type string
		 sb := string(body)
		 log.Printf(sb)

		 cgMap := make(map[string]interface{})
		 errTwo := json.Unmarshal(body, &cgMap)
		 if errTwo != nil {
				panic(err)
		 }

		 for _, record := range cgMap {
			//log.Printf(" [===>] Record: %s", record)
			if rec, ok := record.(map[string]interface{}); ok {
					for key, val := range rec {
						log.Printf(" [========>] %s = %f", key, val)
						value64 := val.(float64)
						intVersion := int64(value64)
						bigIntVersion := big.NewInt(intVersion)
						hashVersion := common.BigToHash(bigIntVersion)
						bytesVersion := common.Hash(hashVersion).Bytes()
						return bytesVersion, remainingGas, nil
					}
			} else {
					log.Printf("record not a map[string]interface{}: %v\n", record)
			}
		}

		//BigToHash sets byte representation of b to hash. If b is larger than len(h), b will be cropped from the left.
		intVersion := int64(0)
		bigIntVersion := big.NewInt(intVersion)
		hashVersion := common.BigToHash(bigIntVersion)
		bytesVersion := common.Hash(hashVersion).Bytes()
		return bytesVersion, remainingGas, nil
	}
}

func createRebalancePrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	getCGVar := newStatefulPrecompileFunction(getCoinGeckoSignature, getCoingeckoPrice(precompileAddr))
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{getCGVar})

	return contract
}

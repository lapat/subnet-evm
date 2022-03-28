package main

import (
	//"errors"
	//"fmt"
	//"math/big"
	"encoding/json"
	//"fmt"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ava-labs/subnet-evm/vmerrs"
	//"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
   "log"
   "net/http"
	// "reflect"

)



func main() {
		log.Printf("createRebalance")
		//to, amount, err := UnpackRebalanceInput(input)
		//_ = to
		//if err != nil {
		///	return nil, remainingGas, err
		//}

		resp, err := http.Get("https://api.coingecko.com/api/v3/simple/token_price/avalanche?contract_addresses=0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7&vs_currencies=usd")
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
						//f := int(value64)
						//amount :=int(val)
						//hashVersion := common.BigToHash(amount)
						//roleBytes := common.Hash(hashVersion).Bytes()
        	}
    	} else {
        	log.Printf("record not a map[string]interface{}: %v\n", record)
    	}
		}



  }

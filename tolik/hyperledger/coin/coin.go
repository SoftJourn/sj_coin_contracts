/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * TODO description
 * author: aokhotnikov@softjourn.com
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	//"bytes"
	//"encoding/json"
    "encoding/hex"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define coin identifier
const coinName string = "sjcoin"

// Color is not defined
var tokenColor int = 0

// Agent (or user) submitting the transaction is not defined
var msgSender []byte = nil

// Minter is not defined
var minter []byte = nil

// Define logger
var logger = shim.NewLogger(coinName)

// Define the Smart Contract structure
type CoinSmartContract struct {
}

// https://github.com/hyperledger/fabric/blob/master/proposals/r1/Custom-Events-High-level-specification.md

/*
 * The Init method is called when the Smart Contract "coin" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *CoinSmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("****************** Init ******************")
    _, args := APIstub.GetFunctionAndParameters()
    var coinKey, indexName string // coin key
    var err error

    if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	tokenColor, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for token color")
	}
	logger.Infof("   Coin token color: %d", tokenColor)

    msgSender, err = s.getAccountBytes(args[1])
    if err != nil {
    	logger.Error(err)
        return shim.Error(err.Error())
    }
    /* Does not make sense as the init is done by peerorg1Admin
	msgSender, err = s.hashCreator(APIstub)
	if err != nil {
  	    logger.Error(err)
        return shim.Error("Failed to hash msg sender")
    }
    */
    logger.Infof("   Coin minter hash: %s", fmt.Sprintf("%x", msgSender))

	//  ==== Index the coin to enable color-based range queries, e.g. check colored coins exist ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~name~color.
	//  This will enable very efficient state range queries based on composite keys matching indexName~name~color~*
	indexName = "name~color~"
	coinKey, err = APIstub.CreateCompositeKey(indexName, []string{coinName, "~", strconv.Itoa(tokenColor)})
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Infof("     Coin index key: %s", coinKey)

    // Check if we have color
    minter, err = APIstub.GetState(coinKey)
    if minter != nil {
        logger.Infof("Color coin minter is %x", minter)
        return shim.Error(fmt.Sprintf("Color %d already exists", tokenColor))
    }
    //minter = msgSender

    /************************* Write the state to the ledger *************************/

	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
    err = APIstub.PutState(coinKey, msgSender)
    if err != nil {
        return shim.Error(err.Error())
    }
    err = APIstub.PutState("tokenColor", []byte(strconv.Itoa(tokenColor)))
    if err != nil {
        return shim.Error(err.Error())
    }

	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *CoinSmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("****************** Invoke ******************")
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	var coinKey, indexName string // coin key
    var err error

    // Check color & load minter
    colorAsBytes, err := APIstub.GetState("tokenColor")
    if err != nil {
        logger.Info("Cannot get coin color for contract")
        return shim.Error(err.Error())
    }

	tokenColor, err = strconv.Atoi(string(colorAsBytes))
	if err != nil {
		return shim.Error("Expecting integer value for token color")
	}
	logger.Infof("   Coin token color: %d", tokenColor)
	indexName = "name~color~"
	coinKey, err = APIstub.CreateCompositeKey(indexName, []string{coinName, "~", strconv.Itoa(tokenColor)})
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Infof("     Coin index key: %s", coinKey)	
	minter, err = APIstub.GetState(coinKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    if minter == nil {
    	return shim.Error("Coin is not initialized. Missing minter hash");
    }
    logger.Infof("   Coin minter hash: %x", minter)
    // Get hash of the agent (or user) submitting the transaction
	msgSender, err = s.hashCreator(APIstub)
	if err != nil {
  	    logger.Error(err)
        return shim.Error("Failed to hash creator")
    }
    logger.Infof("Message sender hash: %x", msgSender)

	// Route to the appropriate handler function to interact with the ledger appropriately
    if function == "getColor" {
		return s.getColor(APIstub, args)
    } else if function == "mint" {
		return s.mint(APIstub, args)
	} else if function == "transfer" {
		return s.transfer(APIstub, args)
	} else if function == "transferFrom" {
		return s.transferFrom(APIstub, args)
	} else if function == "balanceOf" {
		return s.balanceOf(APIstub, args)
	} else if function == "approve" {
		return s.approve(APIstub, args)
	} else if function == "approveAndCall" {
		return s.approveAndCall(APIstub, args)
	} else if function == "allowance" {
		return s.allowance(APIstub, args)
	} else if function == "distribute" {
		return s.distribute(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name")
}

func (s *CoinSmartContract) getColor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    return shim.Success([]byte(strconv.Itoa(tokenColor)))
}
func (s *CoinSmartContract) mint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin mint ###########")
    if reflect.DeepEqual(msgSender, minter) {
    	if len(args)!= 1 {
    		return shim.Error("Incorrect number of arguments. Expecting 1")
    	}
    	Aval, err := strconv.Atoi(args[0])
	    if err != nil {
		    return shim.Error("Expecting integer value for mint amount")
	    }
	    if Aval > 0 {
	    	logger.Infof("Request to mint %d coins", Aval)
	    	return s.addBalance(APIstub, Aval, minter)
	    }
        return shim.Error("Invalid amount to mint coins")	    
    }
    return shim.Error(fmt.Sprintf("You are not allowed to mint this coin: %x - %x", msgSender, minter))
}

func (s *CoinSmartContract) transfer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin transfer ###########")
    if len(args) != 2 {
    	return shim.Error("Incorrect number of arguments. Expecting 2")
    }
    amount, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for transfer amount")
	}
    balance, err := s.getBalance(APIstub, msgSender)
    if err != nil {
    	logger.Info(err.Error())
    }
    if balance < amount {
    	return shim.Error("Insufficient coins");
    }
    user, err := hex.DecodeString(args[1])
    if err != nil {
    	return shim.Error(err.Error())
    }
    result := s.addBalance(APIstub, -amount, msgSender)
    if !reflect.DeepEqual(result, shim.Success([]byte{0x00})) {
        return shim.Error("Failed to transfer coins from sender");
    }
    result = s.addBalance(APIstub, amount, user)
    if !reflect.DeepEqual(result, shim.Success([]byte{0x00})) {
    	result = s.addBalance(APIstub, amount, msgSender)
        return shim.Error("Failed to transfer coins to receiver");
    }
    // TODO: is there a way to send an event?
    // TODO: record transaction - Tactical design: Use TransactionResult to store events
    // https://github.com/hyperledger/fabric/blob/master/proposals/r1/Custom-Events-High-level-specification.md
    return shim.Success([]byte{0x00})
}

func (t *CoinSmartContract) transferFrom(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin transferFrom ###########")
    return shim.Success([]byte{0x00})
}

func (s *CoinSmartContract) balanceOf(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin balanceOf ###########")
    if len(args) != 1 {
    	return shim.Error("Incorrect number of arguments. Expecting 1")
    }
	user, err := s.getAccountBytes(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
    balance, err := s.getBalance(APIstub, user)
    if err != nil {
    	logger.Info(err.Error())
    }
    logger.Infof("Balance[%x] = %d", user, balance)
    return shim.Success([]byte(fmt.Sprintf("%d", balance)))
}

func (s *CoinSmartContract) approve(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin approve ###########")
    var indexName, approveKey string
    if len(args) != 2 {
    	return shim.Error("Incorrect number of arguments. Expecting 2")
    }
    Aval, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for approve amount")
	}
 	user, err := s.getAccountBytes(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
 	indexName = "name~sender~spender~"
	approveKey, err = APIstub.CreateCompositeKey(indexName, []string{hex.EncodeToString(msgSender), "~", hex.EncodeToString(user)})
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Infof("   Approve index key: %s", approveKey)
    err = APIstub.PutState(approveKey, []byte(strconv.Itoa(Aval)))
    if err != nil {
        return shim.Error(err.Error())
    }
    // TODO: Is there a way to send event?
    // Approval(msg.sender, spender, amount);
    return shim.Success([]byte{0x00})
}

func (s *CoinSmartContract) approveAndCall(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin approveAndCall ###########")
    return shim.Success([]byte{0x00})
}

func (s *CoinSmartContract) allowance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin allowance ###########")
    var indexName, approveKey string
    if len(args) != 1 {
    	return shim.Error("Incorrect number of arguments. Expecting 1")
    }
 	user, err := s.getAccountBytes(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Infof("  Approve sender key: %x", msgSender)
	logger.Infof(" Approve spender key: %x", user)
 	indexName = "name~sender~spender~"
	approveKey, err = APIstub.CreateCompositeKey(indexName, []string{hex.EncodeToString(msgSender), "~", hex.EncodeToString(user)})
	if err != nil {
		logger.Info("Failed to create composite key")
		return shim.Error(err.Error())
	}
	logger.Infof("   Approve index key: %s", approveKey) 
	approveAsBytes, err := APIstub.GetState(approveKey)
    if err != nil {
        return shim.Error(err.Error())
    }  
    // TODO: Do we want to convert?
    return shim.Success(approveAsBytes)
}

func (t *CoinSmartContract) distribute(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    logger.Info("########### coin distribute ###########")
    return shim.Success([]byte{0x00})
}

func (s *CoinSmartContract) getAccountBytes(accountText string) ([]byte, error) {
	if accountText == "" || len([]rune(accountText)) != 64 {
		return nil, fmt.Errorf("Expecting 64 character string value for user identity")
	}
	logger.Infof("Account string: %s", accountText)
	accountAsBytes, err := hex.DecodeString(accountText)
	if err != nil {
		logger.Infof("Decoding failed!");
		return nil, err
	}
	return accountAsBytes, nil
}

func (s *CoinSmartContract) getBalance(APIstub shim.ChaincodeStubInterface, user []byte) (int, error) {
    logger.Info("########### coin getBalance ###########")

    indexName := "name~user~"
    balanceKey, err := APIstub.CreateCompositeKey(indexName, []string{coinName, "~", fmt.Sprintf("%x", user)})
	if err != nil {
		return 0, err
	}
	balanceAsBytes, err := APIstub.GetState(balanceKey)
    if err != nil {
        return 0, err
    }
    if balanceAsBytes == nil {
    	return 0, fmt.Errorf("This user has no balance")
    }	
    balance, _ := strconv.Atoi(string(balanceAsBytes))
    return balance, nil
}

func (s *CoinSmartContract) addBalance(APIstub shim.ChaincodeStubInterface, amount int, user []byte) sc.Response {
    logger.Info("########### addBalance ###########")
    balance, err := s.getBalance(APIstub, user)
    if err != nil {
    	logger.Info(err.Error())
    }
    logger.Infof("Current balance is %d coins", balance)
	if (balance+amount) >= 0 {
	    balance += amount
        indexName := "name~user~"
        balanceKey, err := APIstub.CreateCompositeKey(indexName, []string{coinName, "~", fmt.Sprintf("%x", user)})
        logger.Infof("Balance access key is %s", balanceKey)
	    if err != nil {
		    return shim.Error(err.Error())
	    }
        err = APIstub.PutState(balanceKey, []byte(strconv.Itoa(balance)))
        if err != nil {
            return shim.Error(err.Error())
        }
        return shim.Success([]byte{0x00}) //(strconv.Itoa(balance)))
	}
    return shim.Error("Invalid balance change")
}

func (s *CoinSmartContract) hashCreator(APIstub shim.ChaincodeStubInterface) ([]byte, error) {
    //logger.Info("########### coin hashCreator ###########")
    Creatorbytes, err := APIstub.GetCreator()
    //logger.Info(string(Creatorbytes))
    if err != nil {
	    return nil, fmt.Errorf("Failed to get creator")
    }
    if Creatorbytes == nil {
	    return nil, fmt.Errorf("Creator is not found")
    }
    //logger.Info("********************************************************************************")
    //logger.Infof("%x", Creatorbytes)
    Creatorhash, err := factory.GetDefault().Hash(Creatorbytes, &bccsp.SHA256Opts{})
    if err != nil {
	    return nil, fmt.Errorf(fmt.Sprintf("Failed computing SHA256 on [%x]", Creatorbytes))
    }
    //logger.Info("********************************************************************************")
    //logger.Infof("%x", Creatorhash)
    //logger.Info("********************************************************************************")
    return Creatorhash, nil
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(CoinSmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

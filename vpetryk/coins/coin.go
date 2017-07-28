package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"reflect"
	"encoding/json"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp"
	"strconv"
	"github.com/hyperledger/fabric/common/util"
)

var logger = shim.NewLogger("coin_cc")

type CoinChain struct {
}

var minter string = "minter"
var tokenColor string = "tokenColor"
var balances string = "balances"
var allowed string = "allowed"

type CurrentUser struct{
	HashValue []byte
	StringValue string
	BytesValue []byte
}

func (t *CoinChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {

	logger.Info("_____ coin_cc Init _____")

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expected 2")
	}

	color := args[1]
	err := stub.PutState(tokenColor, []byte(color))
	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Info("minter ", args[0])

	err = stub.PutState(minter, []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

	currentUser := t.getCurrentUser(stub)
	logger.Info("Current User ", currentUser.StringValue)

	balancesMap := map[string]uint{args[0]:100, currentUser.StringValue:100}
	t.saveMap(stub, balances, balancesMap);

	allowedMap := map[string]uint{}
	t.saveMap(stub, allowed, allowedMap)

	return shim.Success([]byte(color))
}

func (t *CoinChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	currentUser := t.getCurrentUser(stub)
	logger.Info("sender (current user)", currentUser.StringValue)

	if function == "getColor" {
		return t.getColor(stub, args)
	} else if function == "setColor" {
		return t.setColor(stub, args)
	} else if function == "mint"{
		return t.mint(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "transferFrom" {
		return t.transferFrom(stub, args)
	} else if function == "balanceOf" {
		return t.balanceOf(stub, args)
	} else if function == "approve" {
		return t.approve(stub, args)
	} else if function == "approveAndCall" {
		return t.approveAndCall(stub, args)
	} else if function == "allowance" {
		return t.allowance(stub, args)
	} else if function == "distribute" {
		return t.distribute(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")
}

//func (t *CoinChain) Query(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
//	fmt.Println("query is running " + function)
//
//	if function == "getColor" {
//		return t.getColor(stub, args)
//	}
//	fmt.Println("query did not find func: " + function)
//	return shim.Error("Received unknown function query: " + function)
//}

func (t *CoinChain) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var receiver string = args[0]

	amount := t.parseAmountUint(args[1])
	if amount == 0 {
		return shim.Error("Incorrect amount")
	}

	logger.Info("receiver ", receiver)
	logger.Info("amount ", amount)

	currentUser := t.getCurrentUser(stub)
	balancesMap := t.getMap(stub, balances)

	if balancesMap[currentUser.StringValue] < amount {
		return shim.Error("Not enough coins")
	}

	balancesMap[currentUser.StringValue] -= amount
	balancesMap[receiver] += amount

	t.saveMap(stub, balances, balancesMap)
	//TODO transfer event

	return shim.Success(nil)
}

func (t *CoinChain) transferFrom(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var from string = args[0]
	var to string = args[1]
	amount := t.parseAmountUint(args[2])
	if (amount == 0) {
		return shim.Error("Incorrect amount")
	}

	balancesMap := t.getMap(stub, balances)
	allowedMap := t.getMap(stub, allowed)

	checkAllowedKey, err := stub.CreateCompositeKey(allowed, []string{from, to} )
	if err != nil {
		return shim.Error(err.Error())
	}

	if balancesMap[from] >= amount && allowedMap[checkAllowedKey] >= amount && amount > 0 {
		balancesMap[from] -= amount
		balancesMap[to] += amount
		allowedMap[checkAllowedKey] -= amount
		t.saveMap(stub, balances, balancesMap)
		t.saveMap(stub, allowed, allowedMap)
		//TODO transfer event
		return shim.Success(nil)
	}
	return shim.Error("trnasfer Failed")
}

func (t *CoinChain) approve(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	spender := args[0]
	amount := t.parseAmountUint(args[1])
	if (amount == 0) {
		return shim.Error("Incorrect amount")
	}

	sender := t.getCurrentUser(stub)

	allowedKey, err :=stub.CreateCompositeKey("allowed", []string{sender.StringValue, spender} )
	if err != nil {
		return shim.Error(err.Error())
	}
	allowedMap := t.getMap(stub, allowed)
	allowedMap[allowedKey] += amount
	t.saveMap(stub, allowed, allowedMap)

	//TODO event
	//	Approval(msg.sender, spender, amount);

	return shim.Success(nil)
}

func (t *CoinChain) approveAndCall(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	spender := args[0]
	amount := t.parseAmountUint(args[1])
	if (amount == 0) {
		return shim.Error("Incorrect amount")
	}

	sender := t.getCurrentUser(stub)

	allowedKey, err :=stub.CreateCompositeKey("allowed", []string{sender.StringValue, spender} )
	if err != nil {
		return shim.Error(err.Error())
	}
	allowedMap := t.getMap(stub, allowed)
	allowedMap[allowedKey] += amount
	t.saveMap(stub, allowed, allowedMap)

	//TODO event
	//	Approval(msg.sender, spender, amount);

	function := "receiveApproval"
	queryArgs := util.ToChaincodeArgs(function, "params")
	t.callChaincode(stub, "foundation", queryArgs)
	return shim.Success(nil)
}

func (t *CoinChain) callChaincode(stub shim.ChaincodeStubInterface, chaincodeName string, queryArgs [][]byte) pb.Response {

	channel := "mychannel" //TODO

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channel)
	if response.Status == shim.OK {
		return shim.Success(nil)
	} else {
		return shim.Error("Failed to call chaincode")
	}
}

func (t *CoinChain) setColor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	minterValue, err := stub.GetState(minter)
	if err != nil {
		return shim.Error(err.Error())
	}

	creator, err := t.getCreatorHash(stub);

	if reflect.DeepEqual(creator, minterValue) {
		return shim.Error("User has no permissions")
	}

	color := args[0]
	err = stub.PutState(tokenColor, []byte(color))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *CoinChain) getColor(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	tokenColorValue, err := stub.GetState(tokenColor)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + tokenColor + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(tokenColorValue)
}

func (t *CoinChain) getMap(stub shim.ChaincodeStubInterface, mapName string) map[string]uint {

	logger.Info("------ getMap called")
	mapBytes, err := stub.GetState(mapName)
	if err != nil {
		return nil
	}

	var mapObject map[string]uint
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil
	}
	logger.Info("received map", mapObject)
	return mapObject
}

func (t *CoinChain) saveMap(stub shim.ChaincodeStubInterface, mapName string, mapObject map[string]uint) pb.Response {
	logger.Info("------ saveBalancesMap called")
	balancesMapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(mapName, balancesMapBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("saved ", mapObject)
	return shim.Success(nil)
}

func (t *CoinChain) mint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	minterValue, err := stub.GetState(minter)
	if err != nil {
		return shim.Error(err.Error())
	}

	var minterString string
	json.Unmarshal(minterValue, &minterString)

	currentUser := t.getCurrentUser(stub)

	if currentUser.StringValue != minterString {
		return shim.Error("No permissions")
	}

	balancesMap := t.getMap(stub, balances)

	amount := t.parseAmountUint(args[0])
	if (amount == 0) {
		return shim.Error("Incorrect amount")
	}

	balancesMap[currentUser.StringValue] += amount
	t.saveMap(stub, balances, balancesMap)

	return shim.Success(nil)
}

func (t *CoinChain) distribute(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting at least 3")
	}

	amount := t.parseAmountUint(args[len(args)-1])
	if (amount == 0) {
		return shim.Error("Incorrect amount")
	}
	accounts := args[:len(args)-1]
	logger.Info("accounts: ", accounts)
	logger.Info("amount ", amount)

	sender := t.getCurrentUser(stub)
	balancesMap := t.getMap(stub, balances)

	if balancesMap[sender.StringValue] < amount {
		return shim.Error("Not enough coins")
	}

	mean := amount/uint(len(accounts))
	logger.Warning("mean ", mean)

	var i uint = 0
	logger.Warning("uint(len(accounts)) ", uint(len(accounts)))
	for i < uint(len(accounts)) {
		logger.Warning("i ", i)
		balancesMap[sender.StringValue] -= mean
		logger.Warning("balancesMap[sender ", balancesMap[sender.StringValue])
		logger.Warning("accounts[i] ", accounts[i])
		balancesMap[string(accounts[i])] += mean
		logger.Warning("balancesMap[string(accounts[i])] ", balancesMap[string(accounts[i])])
		i += 1
	}
	t.saveMap(stub, balances,balancesMap)
	return shim.Success(nil)
}

//##### GET ######//

func (t *CoinChain) getCreatorHash(stub shim.ChaincodeStubInterface) ([]byte, error) {
	creatorHash, err := t.hashCreator(stub)
	if err != nil {
		logger.Error(err)
		return []byte{},  err
	}
	return  creatorHash, err
}

func (t *CoinChain) hashCreator(stub shim.ChaincodeStubInterface) ([]byte, error) {
	//logger.Info("########### Coin hashCreator ###########")
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return nil, fmt.Errorf("Failed to get creator")
	}
	if creatorBytes == nil {
		return nil, fmt.Errorf("Creator is not found")
	}
	creatorHash, err := factory.GetDefault().Hash(creatorBytes, &bccsp.SHA256Opts{})
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Failed computing SHA256 on [% x]", creatorBytes))
	}
	return creatorHash, nil
}

func (t *CoinChain) getCurrentUser(stub shim.ChaincodeStubInterface) *CurrentUser {
	creatorHash, _ := t.getCreatorHash(stub)
	creatorStr := fmt.Sprintf("%x", creatorHash)
	return &CurrentUser{creatorHash, creatorStr, []byte(creatorStr)}
}

func (t *CoinChain) balanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	balancesMap := t.getMap(stub, balances)
	return shim.Success([]byte(fmt.Sprintf("%d", balancesMap[args[0]])))
}

func (t *CoinChain) allowance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	sender := args[0]
	spender := args[1]
	allowedMap := t.getMap(stub, allowed)

	allowedKey, err :=stub.CreateCompositeKey(allowed, []string{sender, spender} )
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(fmt.Sprintf("%d", allowedMap[allowedKey])))
}

func (t *CoinChain) parseAmountUint(amount string) uint {
	amount32, err := strconv.ParseUint(amount, 10, 32)
	if err != nil {
		return 0
	}
	return uint(amount32)
}

func main() {
	err := shim.Start(new(CoinChain))
	if err != nil {
		logger.Errorf("Error starting Cinchain: %s", err)
	}
}
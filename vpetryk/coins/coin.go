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

const coinName string = "sjcoin"

var logger = shim.NewLogger(coinName)

type CoinChain struct {
}

var minter string = "minter"
var minterAccount string
var balances string = "balances"
var allowed string = "allowed"

var tokenColor string = "tokenColor"

var channel string = "mychannel"

type CurrentUser struct{
	HashValue []byte
	StringValue string
	BytesValue []byte
}

func (t *CoinChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {

	logger.Info("_____ coin_cc Init _____")

	_, args := stub.GetFunctionAndParameters()

	logger.Info("args:", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expected 2")
	}

	color := args[1]

	tokenColorKey, err := stub.CreateCompositeKey(coinName, []string{color})
	if err != nil {
		return shim.Error(err.Error())
	}

	isKey, err := stub.GetState(tokenColor)
	logger.Info("isKey ", isKey)

	err = stub.PutState(tokenColor, []byte(tokenColorKey))
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("tokenColorKey ", tokenColorKey)

	logger.Info("minterAccount ", args[0])

	minterAccount = args[0]

	minterBytes, err := json.Marshal(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(minter, minterBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	currentUser := t.getCurrentUser(stub)
	logger.Info("Current User ", currentUser.StringValue)

	balancesMap := t.getMap(stub, balances)
	logger.Info("currentBalancesMap ", balancesMap)

	if len(balancesMap) == 0 {
		balancesMap = map[string]uint{args[0]:0, currentUser.StringValue:0}
		t.saveMap(stub, balances, balancesMap);
	}

	allowedMap := t.getMap(stub, allowed)
	logger.Info("currentallowedMap ", allowedMap)
	if len(allowedMap) == 0 {
		allowedMap = map[string]uint{}
		t.saveMap(stub, allowed, allowedMap)
	}

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
	} else if function == "convert" {
		return t.convert(stub, args)
	} else if function == "withdraw" {
		return t.withdraw(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")
}

func (t *CoinChain) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var receiver string = args[0]
	logger.Info("receiver ", receiver)

	logger.Info("args[1] ", args[1])
	amount := t.parseAmountUint(args[1])
	logger.Info("amount ", amount)


	if amount == 0 {
		return shim.Error("Incorrect amount")
	}

	currentUser := t.getCurrentUser(stub)
	logger.Info("Current user: ",currentUser.StringValue)
	balancesMap := t.getMap(stub, balances)

	if balancesMap[currentUser.StringValue] < amount {
		return shim.Error("Not enough coins")
	}

	balancesMap[currentUser.StringValue] -= amount
	balancesMap[receiver] += amount

	t.saveMap(stub, balances, balancesMap)
	//TODO transfer event

	return shim.Success([]byte(strconv.FormatUint(uint64(balancesMap[receiver]), 10)))
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

	tokenColorKey, err := stub.CreateCompositeKey(coinName, []string{color})

	err = stub.PutState(tokenColor, []byte(tokenColorKey))
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
	logger.Info("currentUser.StringValue", currentUser.StringValue)
	logger.Info("minterString", minterString)

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

	return shim.Success([]byte(strconv.FormatUint(uint64(balancesMap[currentUser.StringValue]), 10)))
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

func (t *CoinChain) convert(stub shim.ChaincodeStubInterface, args []string ) pb.Response {

	currency := args[0]
	logger.Info("currency : ", currency)
	account := args[1]
	logger.Info("account : ", account)
	amount := t.parseAmountUint(args[2])
	logger.Info("args[2] : ", args[2])
	logger.Info("amount : ", amount)

	if (amount == 0) {
		return shim.Error("Error. Amount must be > 0")
	}

	queryArgs := util.ToChaincodeArgs("withdraw", account, args[2])
	response := stub.InvokeChaincode(currency, queryArgs, channel)
	logger.Info("Withdraw Result status : ", response.Status)

	if (response.Status == shim.OK){
		balancesMap := t.getMap(stub, balances)
		logger.Info("balancesMap[minterAccount] : ",  balancesMap[minterAccount])
		logger.Info("amount : ", amount)
		if (balancesMap[minterAccount] < amount) {
			return shim.Error("Failed convert. Minter has no money")
		}
		balancesMap[account] += amount
		balancesMap[minterAccount] -= amount
		t.saveMap(stub, balances, balancesMap)
		logger.Info("Convert success : ", amount)
		return shim.Success(nil)
	} else {
		logger.Info("Convert failed : ", amount)
		return shim.Error(response.Message)
	}

	return shim.Success(nil)
}

func (t *CoinChain) withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	account := args[0]
	amount := t.parseAmountUint(args[1])

	balancesMap := t.getMap(stub, balances)

	if balancesMap[account] < amount {
		logger.Error("Withdaraw failed ", amount)
		return shim.Error("Failed withdraw: Not enough funds")
	} else {
		balancesMap[account] -= amount
		balancesMap[minterAccount] += amount
		t.saveMap(stub, balances, balancesMap)
		logger.Info("Withdraw success ", amount)
		return shim.Success(nil)
	}
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
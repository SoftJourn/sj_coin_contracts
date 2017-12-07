package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/common/util"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp"
)

/*
https://github.com/SoftJourn/sj_coin_fabric_POC/blob/master/SJCoins_fabric/fixtures/src/github.com/foundation/foundation.go
*/

var logger *shim.ChaincodeLogger

type FoundationChain struct {
}

type Detail struct {
	Amount uint
	Id uint
	Time time.Time
	Note string
}

/*
https://stackoverflow.com/questions/27676707/marshal-nested-structs-into-json
https://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
*/

type Foundation struct {
	/*Contract's founder address*/
	creatorAccount string

	/*Contract's admin address*/
	adminAccount string

	/*Amount of coins to collect*/
	fundingGoal uint

	/*Amount of coins which were collected before contract has been closed*/
	collectedAmount uint

	/*Amount of coins which were collected after contract has been closed*/
	contractRemains uint

	/*Token address into which should be exchanged all other tokens*/
	mainCurrency string

	/*Contract's deadline(timestamp)*/
	deadline time.Time

	/*Condition of contract closing*/
	closeOnGoalReached bool

	/*Array of currencies which are allowed for contract*/
	acceptCurrencies map[string]bool

	/*Map name of keys: currency + account address, values: amount of donations*/
	donations string

	fundingGoalReached bool

	/*Is contract closed*/
	isContractClosed bool

	/*donations returned */
	isDonationReturned bool

	channelName string

	foundationAccountType string
	userAccountType string

	foundationName string

	Id uint
}

func main() {
	err := shim.Start(new(FoundationChain))
	if err != nil {
		logger.Errorf("Error starting Foundation chaincode: %s", err)
	}
}

func (t *FoundationChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {
    var foundationsMap map[string]Foundation
	foundationsMap = make(map[string]Foundation)
	saveFoundations(stub, foundationsMap)
	return shim.Success(nil)
}

func (t *FoundationChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	
	if (len(args) < 1 ) {
		return shim.Error("Foundation Id is missing")
	}

	foundationsMap := getFoundations(stub)
	if function == "initFoundation" {
		return t.initFoundation(stub, args, foundationsMap)
	} else {	
		// check foundation exists 
		if val, ok := foundationsMap[args[0]]; ok {
			logger.Infof("Found foundation %s", val.foundationName)
		} else {
			return shim.Error("Foundation is not found")
		}
		if function == "receiveApproval" {
			return t.receiveApproval(stub, args, foundationsMap)
		} else if function == "donate" {
			return t.donate(stub, args, foundationsMap)
		} else if function == "close" {
			return t.closeFoundation(stub, args, foundationsMap)
		} else if function == "isWithdrawAllowed" {
			return t.isWithdrawAllowed(stub, args, foundationsMap)
		}
	}

	return shim.Error("Invalid invoke function name.")
}

func (t *FoundationChain) initFoundation(stub shim.ChaincodeStubInterface, args []string, foundationsMap map[string]Foundation) pb.Response {
	/* args
		foundation Name
		admin account
		foundation account
		Goal
		Deadline Minutes
		Close on reached goal
		Currency
		[n, ...] - accept currencies
	*/

	if (len(args) < 8 ) {
		return shim.Error("Incorrect number of arguments. Expecting at least 8")
	}

	foundationName := args[0]

	logger = shim.NewLogger(foundationName)
	logger.Infof("######### %v Init ########", foundationName)

	adminAccount := args[1]
	logger.Info("admin ", adminAccount)

	creatorAccount := args[2]
	logger.Info("creatorAccount ", creatorAccount)

	fundingGoalArg, err := strconv.ParseUint(args[3], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	fundingGoal := uint(fundingGoalArg)
	logger.Info("funding Goal ", fundingGoal)

	minutesInt, err := strconv.ParseInt(args[4], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	duration := time.Minute * time.Duration(minutesInt)

	currentTime := time.Now()
	deadline := currentTime.Add(duration)
	logger.Info("deadline ", deadline.Format(time.RFC3339))

	closeOnGoal, err := strconv.ParseBool(args[5])
	if err != nil {
		return shim.Error(err.Error())
	}

	closeOnGoalReached := closeOnGoal
	logger.Info("closeOnGoalReached ", closeOnGoalReached)

	mainCurrency := args[6]
	logger.Info("mainCurrency ", mainCurrency)

	currencies := args[7:]
	logger.Info("currencies ", currencies)

	acceptCurrencies := make(map[string]bool)
	for _, v := range currencies {
		acceptCurrencies[v] = true
	}
	logger.Info("acceptCurrencies ", acceptCurrencies)

    newFoundation := Foundation{
    	creatorAccount: creatorAccount,
    	adminAccount: adminAccount,
    	fundingGoal: fundingGoal,
        collectedAmount: 0,
        contractRemains: 0,
        mainCurrency: mainCurrency,
        deadline: deadline,
        closeOnGoalReached: closeOnGoalReached,
        acceptCurrencies: acceptCurrencies,
        donations: "donations",
        fundingGoalReached: false,
        isContractClosed: false,
        isDonationReturned: false,
        channelName: "mychannel",
        foundationAccountType: "foundation_",
        userAccountType: "user_",
        foundationName: foundationName,
    	Id: uint(len(foundationsMap) + 1)}

	foundationKey := fmt.Sprintf("%d", newFoundation.Id)
	// check exists
	if val, ok := foundationsMap[foundationKey]; ok {
		logger.Errorf("Foundation already exists: %s", val.foundationName)
		return shim.Error(fmt.Sprintf("Foundation already exists %s", foundationKey))
	}
    foundationsMap[foundationKey] = newFoundation
	saveFoundations(stub, foundationsMap)

	var detailsMap map[int]Detail
	detailsMap = make(map[int]Detail)
	saveDetails(stub, detailsMap, foundationKey)

	return shim.Success(nil)
}

func (t *FoundationChain) donate(stub shim.ChaincodeStubInterface, args []string, foundationsMap map[string]Foundation) pb.Response {

	id, args := args[0], args[1:]
	foundation := foundationsMap[id]

	/* args
		0 - currency name (contract name - coin)
		1 - amount
	*/
	if foundation.isContractClosed {
		return shim.Error("Foundation is closed.")
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	currency := args[0]
	logger.Info("chaincodeName ", currency)

	logger.Info("acceptCurrencies ", foundation.acceptCurrencies)
	if !foundation.acceptCurrencies[currency] {
		return shim.Error("Can not accept currency " + currency)
	}

	amount := t.parseAmountUint(args[1])
	logger.Info("amount ", amount)

	if amount == 0 {
		return shim.Error("Error. Ammount must be > 0")
	}

	queryArgs := util.ToChaincodeArgs("transfer", foundation.foundationAccountType, foundation.foundationName, args[1])
	response := stub.InvokeChaincode(currency, queryArgs, foundation.channelName)
	logger.Info("Result ", response.Status)

	if (response.Status == shim.OK) {

		currentUser := t.getCurrentUser(stub)
		logger.Info(currentUser)

		donationsKey := fmt.Sprintf("donations_%d", foundation.Id);
		donationsMap := getMap(stub, donationsKey)

		donationKey, err := stub.CreateCompositeKey(currency, []string{foundation.userAccountType, currentUser.StringValue})
		if err != nil {
			return shim.Error(err.Error())
		}

		donationsMap[donationKey] += amount
		saveMap(stub, donationsKey, donationsMap)

		foundation.collectedAmount += amount
		logger.Info("AmountRaised ", foundation.collectedAmount)
		checkGoalReached(foundation)
		logger.Info("fundingGoalReached ", foundation.fundingGoalReached)
		logger.Info("isContractClosed ", foundation.isContractClosed)

		return shim.Success([]byte(strconv.FormatUint(uint64(foundation.collectedAmount), 10)))
	} else {
		return shim.Error(response.Message)
	}
}

func (t *FoundationChain) closeFoundation(stub shim.ChaincodeStubInterface, args []string, foundationsMap map[string]Foundation) pb.Response {

	id, args := args[0], args[1:]
	foundation := foundationsMap[id]

	checkGoalReached(foundation)

	if foundation.isContractClosed {
		return shim.Error("Failed. Foundation is alredy closed.")
	}

	currentUser := t.getCurrentUser(stub)

	if (currentUser.StringValue != foundation.adminAccount) {
		return shim.Error( "Failed. Only admin can close foundation" )
	}

	if !foundation.fundingGoalReached {
		if !foundation.isDonationReturned {
			donationsKey := fmt.Sprintf("donations_%d", foundation.Id);
			donationsMap := getMap(stub, donationsKey)
			for k, v := range donationsMap {
				if v > 0 {
					currency, parts, err := stub.SplitCompositeKey(k)
					logger.Info("Key : ", k)
					logger.Info("currency: ", currency)
					logger.Info("parts: ", parts)
					logger.Info("amount value v: ", v)

					if err != nil {
						return shim.Error(err.Error())
					}

					queryArgs := util.ToChaincodeArgs("transfer", foundation.userAccountType, parts[0], strconv.FormatUint(uint64(v), 10))
					response := stub.InvokeChaincode(currency, queryArgs, foundation.channelName)
					logger.Info("Result ", response.Status)

					if (response.Status == shim.OK){

					} else {
						return shim.Error(response.Message)
					}
					//donationsMap[k] = 0;
				}
			}
			saveMap(stub, donationsKey, donationsMap)
			foundation.isDonationReturned = true
		}
	} else {
		foundation.contractRemains = foundation.collectedAmount
		logger.Info("contractRemains ", foundation.contractRemains)
	}
	foundation.isContractClosed = true
	return shim.Success([]byte(strconv.FormatUint(uint64(foundation.contractRemains), 10)))
}

func checkGoalReached(foundation Foundation) bool {

	if foundation.collectedAmount >= foundation.fundingGoal {
		foundation.fundingGoalReached = true
	}

	if foundation.closeOnGoalReached {
		if foundation.collectedAmount >= foundation.fundingGoal || time.Now().After(foundation.deadline) {
			foundation.contractRemains = foundation.collectedAmount
			foundation.isContractClosed = true
		}
	}
	return foundation.fundingGoalReached
}

func (t *FoundationChain) isWithdrawAllowed(stub shim.ChaincodeStubInterface, args []string, foundationsMap map[string]Foundation) pb.Response {
	logger.Info("    ---- invoked isWithdrawAllowed")

	id, args := args[0], args[1:]
	foundation := foundationsMap[id]

	amount := t.parseAmountUint(args[0])
	note := args[1]
	logger.Info("amount: ", amount)
	logger.Info("note: ", note)
	logger.Info("contractRemains: ", foundation.contractRemains)

	result := false
	currentUser := t.getCurrentUser(stub)

	if (currentUser.StringValue == foundation.adminAccount && foundation.isContractClosed && amount <= foundation.contractRemains) {
		foundation.contractRemains -= amount
		detailsKey := fmt.Sprintf("details_%s", foundation.Id)
		detailsMap, err := getDetails(stub, detailsKey)
		if err != nil {
			return shim.Error(err.Error())
		}

		newDetail := Detail{Time:time.Now(), Amount:amount, Note:note, Id: uint(len(detailsMap) + 1)}
		detailsMap[len(detailsMap) + 1] = newDetail
		saveDetails(stub, detailsMap, detailsKey)
		logger.Info("detailsMap: ", detailsMap)

		result = true
	}
	logger.Info("    ---- isWithdrawAllowed result", result)
	return shim.Success([]byte(strconv.FormatBool(result)))
}

func (t *FoundationChain) parseAmountUint(amount string) uint {
	amount32, err := strconv.ParseUint(amount, 10, 32)
	if err != nil {
		return 0
	}
	return uint(amount32)
}

func (t *FoundationChain) receiveApproval(stub shim.ChaincodeStubInterface, args []string, foundationsMap map[string]Foundation) pb.Response {
	id, args := args[0], args[1:]
	foundation := foundationsMap[id]
	logger.Info("receiveApproval not implemented for %s", foundation.foundationName)
	return shim.Success(nil)
}

func getFoundations(stub shim.ChaincodeStubInterface) map[string]Foundation {

	logger.Info("------ getFoundations called")
	mapBytes, err := stub.GetState("foundations")
	if err != nil {
		return nil
	}

	var mapObject map[string]Foundation
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil
	}
	logger.Info("received map", mapObject)
	return mapObject
}

func saveFoundations(stub shim.ChaincodeStubInterface, mapObject map[string]Foundation) pb.Response {
	logger.Info("------ saveMap called")
	mapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("foundations", mapBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("saved ", mapObject)
	return shim.Success(nil)
}

func getMap(stub shim.ChaincodeStubInterface, mapName string) map[string]uint {

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

func saveMap(stub shim.ChaincodeStubInterface, mapName string, mapObject map[string]uint) pb.Response {
	logger.Info("------ saveMap called")
	mapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(mapName, mapBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("saved ", mapObject)
	return shim.Success(nil)
}

var details string = "details"

func getDetails(stub shim.ChaincodeStubInterface, foundationId string) (map[int]Detail, error) {

	logger.Info("------ getDetails called")
	mapBytes, err := stub.GetState(fmt.Sprintf("details_%s", foundationId))
	logger.Info("mapBytes", mapBytes)
	if err != nil {
		return nil, err
	}

	var mapObject map[int]Detail
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil, err
	}
	logger.Info("received Details map", mapObject)
	return mapObject, nil
}

func saveDetails(stub shim.ChaincodeStubInterface, mapObject map[int]Detail, foundationId string) pb.Response {
	logger.Info("------ saveDetails called")

	mapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(fmt.Sprintf("details_%s", foundationId), mapBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("saved ", mapObject)
	return shim.Success(nil)
}


//##### GET ######//

type CurrentUser struct{
	HashValue []byte
	StringValue string
	BytesValue []byte
}

func (t *FoundationChain) getCreatorHash(stub shim.ChaincodeStubInterface) ([]byte, error) {
	creatorHash, err := t.hashCreator(stub)
	if err != nil {
		logger.Error(err)
		return []byte{},  err
	}
	return  creatorHash, err
}

func (t *FoundationChain) hashCreator(stub shim.ChaincodeStubInterface) ([]byte, error) {
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

func (t *FoundationChain) getCurrentUser(stub shim.ChaincodeStubInterface) *CurrentUser {
	creatorHash, _ := t.getCreatorHash(stub)
	creatorStr := fmt.Sprintf("%x", creatorHash)
	return &CurrentUser{creatorHash, creatorStr, []byte(creatorStr)}
}

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

var logger = shim.NewLogger("foundation_cc")

type FoundationChain struct {
}

type Detail struct {
	amount uint
	id uint
	time time.Time
	note string
}

/*Contract's founder address*/
var foundationAccount string

/*Contract's admin address*/
var adminAccount string

/*Amount of coins to collect*/
var fundingGoal uint

/*Amount of coins which were collected before contract has been closed*/
var collectedAmount uint = 0

/*Amount of coins which were collected after contract has been closed*/
var contractRemains uint = 0

/*Token address into which should be exchanged all other tokens*/
var mainCurrency string;

/*Contract's deadline(timestamp)*/
var deadline time.Time

/*Condition of contract closing*/
var closeOnGoalReached bool

/*Array of currencies which are allowed for contract*/
var acceptCurrencies map[string]bool

/*Map name of keys: currency + account address, values: amount of donations*/
var donations string = "donations"

var fundingGoalReached bool = false

/*Is contract closed*/
var isContractClosed bool = false

/*donations returned */
var isDonationReturned bool = false

/*donations converted */
var isDonationConverted bool = false

var channel string = "mychannel"

func main() {
	err := shim.Start(new(FoundationChain))
	if err != nil {
		logger.Errorf("Error starting Foundation chaincode: %s", err)
	}
}

func (t *FoundationChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {

	/* args
		0 - admin
		1 - foundation account
		2 - Goal
		3 - Deadline Minutes
		4 - Close on reached goal
		5 - Currency
		[6, ...] - accept currencies
	*/

	logger.Info("########### foundation_cc Init ###########")
	_, args := stub.GetFunctionAndParameters()

	if (len(args) < 7 ) {
		return shim.Error("Incorrect number of arguments. Expecting at least 7")
	}

	adminAccount = args[0]
	logger.Info("admin ", adminAccount)

	foundationAccount = args[1]
	logger.Info("foundation ", foundationAccount)

	fundingGoalArg, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	fundingGoal = uint(fundingGoalArg)
	logger.Info("fundingGoal ", fundingGoal)

	minutesInt, err := strconv.ParseInt(args[3], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	duration := time.Minute * time.Duration(minutesInt)

	currentTime := time.Now()
	deadline = currentTime.Add(duration)
	logger.Info("deadline ", deadline)

	closeOnGoal, err := strconv.ParseBool(args[4])
	if err != nil {
		return shim.Error(err.Error())
	}

	closeOnGoalReached = closeOnGoal
	logger.Info("closeOnGoalReached ", closeOnGoalReached)

	mainCurrency = args[5]
	logger.Info("mainCurrency ", mainCurrency)

	currencies := args[6:]
	logger.Info("currencies ", currencies)

	acceptCurrencies = make(map[string]bool)
	for _, v := range currencies {
		acceptCurrencies[v] = true
	}
	logger.Info("acceptCurrencies ", acceptCurrencies)

	donationsMap := getMap(stub, donations)
//	if len(donationsMap) == 0 {
		donationsMap = make(map[string]uint)
		saveMap(stub, donations, donationsMap)
//	}

	detailsMap, _ := getDetails(stub)
	detailsMap = make(map[string]string)
	saveDetails(stub, detailsMap)

	logger.Info("acceptCurrencies ", acceptCurrencies)

	return shim.Success(nil)
}

func (t *FoundationChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "receiveApproval" {
		return t.receiveApproval(stub, args)
	} else if function == "donate" {
		return t.donate(stub, args)
	} else if function == "close" {
		return t.close(stub, args)
	} else if function == "convertToMainCurrency" {
		return t.convertToMainCurrency(stub, args)
	} else if function == "sendToFoundation" {
		return t.sendToFoundation(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

func (t *FoundationChain) donate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/* args
		0 - currency name (contract name - coin)
		1 - amount
	*/
	if isContractClosed {
		return shim.Error("Foundation is closed.")
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	currency := args[0]
	logger.Info("chaincodeName ", currency)

	logger.Info("acceptCurrencies ", acceptCurrencies)
	if !acceptCurrencies[currency] {
		return shim.Error("Can not accept currency " + currency)
	}

	amount := t.parseAmountUint(args[1])
	logger.Info("amount ", amount)

	if amount == 0 {
		return shim.Error("Error. Ammount must be > 0")
	}

	queryArgs := util.ToChaincodeArgs("transfer", adminAccount, args[1])
	response := stub.InvokeChaincode(currency, queryArgs, channel)
	logger.Info("Result ", response.Status)

	if (response.Status == shim.OK){

		currentUser := t.getCurrentUser(stub)
		logger.Info(currentUser)

		donationsMap := getMap(stub, donations)
		donationKey, err := stub.CreateCompositeKey("_", []string{currency, currentUser.StringValue })
		if err != nil {
			return shim.Error(err.Error())
		}

		donationsMap[donationKey] += amount
		saveMap(stub, donations, donationsMap)

		collectedAmount += amount
		logger.Info("AmountRaised ", collectedAmount)
		checkGoalReached()
		logger.Info("fundingGoalReached ", fundingGoalReached)
		logger.Info("isContractClosed ", isContractClosed)

		return shim.Success([]byte(strconv.FormatUint(uint64(collectedAmount), 10)))
	} else {
		return shim.Error(response.Message)
	}
}

func (t *FoundationChain) close(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	checkGoalReached()

	if isContractClosed {
		return shim.Error("Failed. Foundation is alredy closed.")
	}

	currentUser := t.getCurrentUser(stub)

	if (currentUser.StringValue != adminAccount) {
		return shim.Error( "Failed. Only admin can close foundation" )
	}

	if !fundingGoalReached {
		if !isDonationReturned {
			donationsMap := getMap(stub, donations)
			for k, v := range donationsMap {
				if v > 0 {
					someString, parts, err := stub.SplitCompositeKey(k)
					logger.Info("Key : ", k)
					logger.Info("someString: ", someString)
					logger.Info("parts: ", parts)
					logger.Info("amount value v: ", v)

					if err != nil {
						return shim.Error(err.Error())
					}

					queryArgs := util.ToChaincodeArgs("transfer", parts[1], strconv.FormatUint(uint64(v), 10))
					response := stub.InvokeChaincode(parts[0], queryArgs, channel)
					logger.Info("Result ", response.Status)

					if (response.Status == shim.OK){

					} else {
						return shim.Error(response.Message)
					}
					//donationsMap[k] = 0;
				}
			}
			saveMap(stub, donations, donationsMap)
			isDonationReturned = true
		}
	} else {
		contractRemains = collectedAmount
	}
	isContractClosed = true
	return shim.Success(nil)
}

func checkGoalReached() bool {

	if collectedAmount >= fundingGoal {
		fundingGoalReached = true
	}

	if closeOnGoalReached {
		if collectedAmount >= fundingGoal || time.Now().After(deadline) {
			isContractClosed = true
		}
	}
	return fundingGoalReached
}


func (t *FoundationChain) convertToMainCurrency(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if (!isContractClosed) {
		return shim.Error("Failed. Please close donations before conversion.")
	}

	if isDonationConverted {
		return shim.Error("Failed. Donations are already converted.")
	}

	donationsMap := getMap(stub, donations)
	for k, v := range donationsMap {
		if v > 0 {
			someString, parts, err := stub.SplitCompositeKey(k)
			logger.Info("Key : ", k)
			logger.Info("someString: ", someString)
			logger.Info("parts: ", parts)
			logger.Info("--- amount value v: ", v)

			if err != nil {
				return shim.Error(err.Error())
			}

			if parts[0] != mainCurrency {

				queryArgs := util.ToChaincodeArgs("convert", parts[0], adminAccount, strconv.FormatUint(uint64(v), 10))
				response := stub.InvokeChaincode(mainCurrency, queryArgs, channel)
				logger.Info("Result ", response.Status)

				if (response.Status == shim.OK){

				} else {
					return shim.Error(response.Message)
				}
			}

		}
	}
	saveMap(stub, donations, donationsMap)
	isDonationConverted = true
	return shim.Success(nil)
}

func (t *FoundationChain) sendToFoundation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	currentUser := t.getCurrentUser(stub)
	if (currentUser.StringValue != adminAccount) {
		return shim.Error( "Failed. Only admin can spend funds" )
	}
	if !fundingGoalReached {
		return shim.Error("Failed. Funding goal in not reached.")
	}
	if !isContractClosed {
		return shim.Error("Failed. Please close donations.")
	}
	if len(acceptCurrencies) != 1 && !isDonationConverted {
		return shim.Error("Failed. Please convert funds to main currency.")
	}

	logger.Info("args[0] amount ", args[0])
	amount := t.parseAmountUint(args[0])
	logger.Info("amount ", amount)
	note := args[1]
	logger.Info("note ", note)

	if amount <= collectedAmount {
		queryArgs := util.ToChaincodeArgs("transfer", foundationAccount, args[0])
		response := stub.InvokeChaincode(mainCurrency, queryArgs, channel)
		logger.Info("Result ", response.Status)

		if response.Status == shim.OK {
			collectedAmount -= amount

			detailMap, err := getDetails(stub)
			if err != nil {
				return shim.Error(err.Error())
			}

			key, err := stub.CreateCompositeKey("_", []string{strconv.Itoa(len(detailMap) + 1), time.Now().String(), args[0] })
			if err != nil {
				return shim.Error(err.Error())
			}

			detailMap[key] = note
			saveDetails(stub, detailMap)

			return shim.Success([]byte(strconv.FormatUint(uint64(collectedAmount), 10)))
		} else {
			return shim.Error(response.Message)
		}

	} else {
		return shim.Error("Failed. Not enough funds.")
	}
}

func (t *FoundationChain) parseAmountUint(amount string) uint {
	amount32, err := strconv.ParseUint(amount, 10, 32)
	if err != nil {
		return 0
	}
	return uint(amount32)
}

func (t *FoundationChain) receiveApproval(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func getDetails(stub shim.ChaincodeStubInterface) (map[string]string, error) {

	logger.Info("------ getDetails called")
	mapBytes, err := stub.GetState(details)
	logger.Info("mapBytes", mapBytes)
	if err != nil {
		return nil, err
	}

	var mapObject map[string]string
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil, err
	}
	logger.Info("received Details map", mapObject)
	return mapObject, nil
}

func saveDetails(stub shim.ChaincodeStubInterface, mapObject map[string]string) pb.Response {
	logger.Info("------ saveDetails called")

	mapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(details, mapBytes)
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
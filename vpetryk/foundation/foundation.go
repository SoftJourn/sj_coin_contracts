package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/common/util"
	"encoding/json"
	"fmt"
	"strings"
	"encoding/pem"
	"crypto/x509"
)

var logger *shim.ChaincodeLogger

type FoundationChain struct {
}

type Detail struct {
	Amount uint
	Id uint
	Time time.Time
	Note string
}

type Foundation struct {
	Name string
	CreatorAccount string
	AdminAccount string
	FundingGoal uint
	CollectedAmount uint
	ContractRemains uint
	MainCurrency string
	Deadline time.Time
	CloseOnGoalReached bool
	AcceptCurrencies map[string]bool
	DonationsMap map[string]uint
	DetailsMap map[int]Detail
	FundingGoalReached bool
	IsContractClosed bool
	IsDonationReturned bool
}

var channelName string = "mychannel"
var foundationAccountType string = "foundation_"
var userAccountType string = "user_"

var foundationsKey string = "foundations"

func main() {
	err := shim.Start(new(FoundationChain))
	if err != nil {
		logger.Errorf("Error starting Foundation chaincode: %s", err)
	}
}

func (t *FoundationChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {

	//_, args := stub.GetFunctionAndParameters()
	logger = shim.NewLogger(foundationsKey)
	logger.Infof("######### %v Init ########", foundationsKey)

	mapBytes, err := stub.GetState(foundationsKey)
	if err != nil {
		logger.Info("Init error ", err)
	}
	logger.Infof("Init mapBytes %s: ", mapBytes)

	if len(mapBytes) == 0 {
		foundationsMap := make(map[string]Foundation)
		err = saveFoundations(stub, foundationsMap)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	return shim.Success(nil)
}

func (t *FoundationChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createFoundation" {
		return t.createFoundation(stub, args)
	} else if function == "receiveApproval" {
		return t.receiveApproval(stub, args)
	} else if function == "donate" {
		return t.donate(stub, args)
	} else if function == "close" {
		return t.closeFoundation(stub, args)
	} else if function == "isWithdrawAllowed" {
		return t.isWithdrawAllowed(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

func (t *FoundationChain) createFoundation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/* args
		0 - foundation Name
		1 - admin account
		2 -  foundation account
		3 - Goal
		4 - Deadline Minutes
		5 - Close on reached goal
		6 - Currency
		... n - accept currencies
	*/

	if len(args) < 8 {
		return shim.Error("Incorrect number of arguments. Expecting at least 8")
	}

	foundations, err := getFoundations(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	_, exist := foundations[args[0]]
	if exist {
		return shim.Error("Foundation already exists.")
	}

	foundation := Foundation{}
	foundation.Name = args[0]
	logger.Info("foundationName ", foundation.Name)

	foundation.AdminAccount = args[1]
	logger.Info("admin ", foundation.AdminAccount)

	foundation.CreatorAccount = args[2]
	logger.Info("creatorAccount ", foundation.CreatorAccount)

	fundingGoalArg, err := strconv.ParseUint(args[3], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	foundation.FundingGoal = uint(fundingGoalArg)
	logger.Info("funding Goal ", foundation.FundingGoal)

	minutesInt, err := strconv.ParseInt(args[4], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	duration := time.Minute * time.Duration(minutesInt)
	currentTime := time.Now()
	foundation.Deadline = currentTime.Add(duration)
	logger.Info("deadline ", foundation.Deadline.Format(time.RFC3339))

	closeOnGoal, err := strconv.ParseBool(args[5])
	if err != nil {
		return shim.Error(err.Error())
	}

	foundation.CloseOnGoalReached = closeOnGoal
	logger.Info("closeOnGoalReached ", foundation.CloseOnGoalReached)

	foundation.MainCurrency = args[6]
	logger.Info("mainCurrency ", foundation.MainCurrency)

	currencies := args[7:]
	logger.Info("currencies ", currencies)

	foundation.AcceptCurrencies = make(map[string]bool)
	for _, v := range currencies {
		foundation.AcceptCurrencies[v] = true
	}
	logger.Info("acceptCurrencies ", foundation.AcceptCurrencies)

	foundation.DonationsMap = make(map[string]uint)
	foundation.DetailsMap = make(map[int]Detail)
	foundations[foundation.Name] = foundation
	err = saveFoundations(stub, foundations)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *FoundationChain) donate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/* args
		0 - currency name (contract name - coin)
		1 - amount
		2 - foundation name
	*/

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	foundations, err := getFoundations(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	foundation, exist := foundations[args[2]]
	if !exist {
		return shim.Error("Foundation does not exist.")
	}

	if foundation.IsContractClosed {
		return shim.Error("Foundation is closed.")
	}

	currency := args[0]
	logger.Info("chaincodeName: ", currency)

	logger.Info("acceptCurrencies ", foundation.AcceptCurrencies)
	if !foundation.AcceptCurrencies[currency] {
		return shim.Error("Can not accept currency " + currency)
	}

	amount := t.parseAmountUint(args[1])
	logger.Info("amount ", amount)

	if amount == 0 {
		return shim.Error("Error. Amount must be > 0")
	}

	queryArgs := util.ToChaincodeArgs("transfer", foundationAccountType, foundation.Name, args[1])
	response := stub.InvokeChaincode(currency, queryArgs, channelName)
	logger.Info("Result ", response.Status)

	if response.Status == shim.OK {

		currentUserId, err := getCurrentUserId(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		donationKey, err := stub.CreateCompositeKey(currency, []string{userAccountType, currentUserId})
		if err != nil {
			return shim.Error(err.Error())
		}

		foundation.DonationsMap[donationKey] += amount
		foundation.CollectedAmount += amount
		logger.Info("foundation.CollectedAmount ", foundation.CollectedAmount)
		checkGoalReached(&foundation)

		logger.Info("donate checkGoalReached FundingGoalReached: ", foundation.FundingGoalReached)
		logger.Info("fundingGoalReached ", foundation.FundingGoalReached)
		logger.Info("isContractClosed ", foundation.IsContractClosed)

		foundations[foundation.Name] = foundation
		err = saveFoundations(stub, foundations)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success([]byte(strconv.FormatUint(uint64(foundation.CollectedAmount), 10)))
	} else {
		return shim.Error(response.Message)
	}
}

func (t *FoundationChain) closeFoundation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/* args
		0 - foundation name
	*/

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	foundations, err := getFoundations(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	foundation, ok := foundations[args[0]]
	if !ok {
		return shim.Error("Foundation does not exist.")
	}

	checkGoalReached(&foundation)

	if foundation.IsContractClosed {
		return shim.Error("Failed. Foundation is alredy closed.")
	}

	currentUserId, err := getCurrentUserId(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if currentUserId != foundation.AdminAccount {
		return shim.Error( "Failed. Only admin can close foundation" )
	}

	if !foundation.FundingGoalReached {
		if !foundation.IsDonationReturned {
			//donationsMap := getMap(stub, donations)
			for k, v := range foundation.DonationsMap {
				if v > 0 {
					currency, parts, err := stub.SplitCompositeKey(k)
					logger.Info("Key : ", k)
					logger.Info("currency: ", currency)
					logger.Info("parts: ", parts)
					logger.Info("amount value v: ", v)

					if err != nil {
						return shim.Error(err.Error())
					}

					queryArgs := util.ToChaincodeArgs("transfer", userAccountType, parts[1], strconv.FormatUint(uint64(v), 10))
					response := stub.InvokeChaincode(currency, queryArgs, channelName)
					logger.Info("Result ", response.Status)

					if response.Status == shim.OK {

					} else {
						return shim.Error(response.Message)
					}
					//foundation.DonationsMap[k] = 0;
				}
			}
			//saveMap(stub, donations, donationsMap)
			foundation.IsDonationReturned = true

			//TODO
			//foundations[foundation.Name] = foundation
			//saveFoundations(stub, foundations)
		}
	} else {
		foundation.ContractRemains = foundation.CollectedAmount
		logger.Info("contractRemains ", foundation.ContractRemains)
	}

	foundation.IsContractClosed = true
	foundations[foundation.Name] = foundation
	err = saveFoundations(stub, foundations)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success([]byte(strconv.FormatUint(uint64(foundation.ContractRemains), 10)))
}

func (t *FoundationChain) isWithdrawAllowed(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/* args
		0 - foundation name
		1 - amount
		2 - note
	*/

	logger.Info("    ---- invoked isWithdrawAllowed")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	foundations, err := getFoundations(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	foundation, ok := foundations[args[0]]
	if !ok {
		return shim.Error("Foundation does not exist.")
	}

	amount := t.parseAmountUint(args[1])
	note := args[2]
	logger.Info("amount: ", amount)
	logger.Info("note: ", note)
	logger.Info("contractRemains: ", foundation.ContractRemains)

	result := false

	currentUserId, err := getCurrentUserId(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if currentUserId == foundation.AdminAccount && foundation.IsContractClosed && amount <= foundation.ContractRemains {
		foundation.ContractRemains -= amount
		//detailsMap, err := getDetails(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		newDetail := Detail{Time:time.Now(), Amount:amount, Note:note, Id: uint(len(foundation.DetailsMap) + 1)}
		foundation.DetailsMap[len(foundation.DetailsMap) + 1] = newDetail
		//saveDetails(stub, detailsMap)
		logger.Info("detailsMap: ", foundation.DetailsMap)

		foundations[foundation.Name] = foundation
		err = saveFoundations(stub, foundations)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = true
	}
	logger.Info("    ---- isWithdrawAllowed result", result)
	return shim.Success([]byte(strconv.FormatBool(result)))
}

func checkGoalReached(foundation *Foundation) bool {

	if foundation.CollectedAmount >= foundation.FundingGoal {
		foundation.FundingGoalReached = true
	}

	if foundation.CloseOnGoalReached {
		if foundation.CollectedAmount >= foundation.FundingGoal || time.Now().After(foundation.Deadline) {
			foundation.ContractRemains = foundation.CollectedAmount
			foundation.IsContractClosed = true
		}
	}
	logger.Info("checkGoalReached FundingGoalReached: ", foundation.FundingGoalReached)
	return foundation.FundingGoalReached
}

func getCurrentUserId(stub shim.ChaincodeStubInterface) (string, error) {

	var userId string

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return userId, err
	}

	creatorString :=fmt.Sprintf("%s",creatorBytes)
	index := strings.Index(creatorString, "-----BEGIN CERTIFICATE-----")
	certificate := creatorString[index:]
	block, _ := pem.Decode([]byte(certificate))

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return userId, err
	}

	userId = cert.Subject.CommonName
	logger.Infof("---- UserID: %v ", userId)
	return userId, err
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

func getFoundations(stub shim.ChaincodeStubInterface) (map[string]Foundation, error) {

	logger.Info("------ getFoundations called")
	mapBytes, err := stub.GetState(foundationsKey)
	if err != nil {
		return nil, err
	}

	var mapObject map[string]Foundation
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil, err
	}
	logger.Info("received Foundations map %s", mapObject)
	return mapObject, nil
}

func saveFoundations(stub shim.ChaincodeStubInterface, mapObject map[string]Foundation) error {
	logger.Info("------ saveFoundations called")

	mapBytes, err := json.Marshal(mapObject)
	if err != nil {
		return err
	}
	err = stub.PutState(foundationsKey, mapBytes)
	if err != nil {
		return err
	}
	logger.Info("saved ", mapObject)
	return nil
}
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

var logger *shim.ChaincodeLogger

var orgId string
var usersKey = "users"

var users map[string]UserData

type UserData struct {
	Email string
	FirstName string
	LastName string
	PersonId string
	PersistentFaceId string
	PersonGroupId string
}

type UsersChain struct {
}

func (t *UsersChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {

	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected 1")
	}
	if len(args[0]) == 0 {
		return shim.Error("Incorrect Organization ID")
	}

	orgId = args[0]

	logger = shim.NewLogger(orgId)
	logger.Infof("orgId: %s", orgId)

	users = t.getUsers(stub)
	if users == nil {
		users = make(map[string]UserData)
	}
	t.saveUsers(stub, users)

	return shim.Success(nil)
}

func (t *UsersChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	logger.Info("invoke is running " + function)

	if function == "addUser" {
		return t.addUser(stub, args)
	} else if function == "getUserDataById" {
		return t.getUserDataById(stub, args)
	}
	logger.Info("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")

}

func (t *UsersChain) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Infof("args: %v", args)

	if len(args) != 1 || len(args[0]) == 0 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var userData UserData
	err := json.Unmarshal([]byte(args[0]), &userData)
	if err != nil {
		logger.Errorf("\nerr: %v\n", err)
		return shim.Error(err.Error())
	}
	logger.Infof("\nuserData: %v\n", userData)

	users[userData.Email] = userData
	t.saveUsers(stub, users)

	return shim.Success(nil)
}

func (t *UsersChain) getUserDataById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Infof("args: %v", args)

	if len(args) != 1 || len(args[0]) == 0 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if value, ok := users[args[0]]; ok {
		userBytes, err := json.Marshal(value)
		if err != nil {
			return shim.Error(err.Error())
		}
		logger.Infof("userBytes: %s", userBytes)
		return shim.Success(userBytes)
	} else {
		logger.Errorf("Error: %v", "User not found")
		return shim.Error("User not found")
	}
}

func (t *UsersChain) getUsers(stub shim.ChaincodeStubInterface) map[string]UserData {

	logger.Info("------ getMap called")
	mapBytes, err := stub.GetState(usersKey)
	if err != nil {
		return nil
	}

	var mapObject map[string]UserData
	err = json.Unmarshal(mapBytes, &mapObject)
	if err != nil {
		return nil
	}
	logger.Info("received map", mapObject)
	return mapObject
}

func (t *UsersChain) saveUsers(stub shim.ChaincodeStubInterface, mapObject map[string]UserData) pb.Response {
	logger.Info("------ saveUsers called")
	userData, err := json.Marshal(mapObject)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(usersKey, userData)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("saved ", mapObject)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(UsersChain))
	if err != nil {
		logger.Errorf("Error starting Cinchain: %s", err)
	}
}
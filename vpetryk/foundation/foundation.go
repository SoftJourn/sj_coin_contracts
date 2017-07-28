package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("foundation_cc")

type FoundationChain struct {
}

type Detail struct {
	amount uint
	id uint
	time uint
	note string
}

func main() {
	err := shim.Start(new(FoundationChain))
	if err != nil {
		logger.Errorf("Error starting Foundation chaincode: %s", err)
	}
}


func (t *FoundationChain) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### foundation_cc Init ###########")
	//_, args := stub.GetFunctionAndParameters()

	return shim.Success(nil)
}

func (t *FoundationChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "receiveApproval" {
		return t.receiveApproval(stub, args)
	//} else if function == "query" {
	//	return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

func (t *FoundationChain) receiveApproval(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

//function receiveApproval(address _from, uint _value, address _token, bytes _extraData) returns (bool) {
///* If goal is not reached then - donate */
//	if (!fundingGoalReached) {
//	/* If it is not too late */
//		if (now <= deadline) {
//			return donate(_from, _value, _token);
//		}
//		else {
//			return false;
//		}
//	}
//	/* If goal is reached then - exchange all collected tokens into one token*/
//	else if (fundingGoalReached && mainToken == _token) {
//		if (contractRemains == _value) {
//			return exchange(_from, _value, _token);
//		}
//		else {
//			return false;
//		}
//	}
//	else {
//		return false;
//	}
//}
## Balance transfer (identify caller eample)

A patch to Balance transfer sample Node.js app https://github.com/hyperledger/fabric-samples/tree/release/balance-transfer
to demonstrate identity caller

### Call method

Creator is the serialized identity structure:

```
message SerializedIdentity {
    // The identifier of the associated membership service provider
    string mspid = 1;

    // the Identity, serialized according to the rules of its MPS
    bytes id_bytes = 2;
}
```

To create a serialized identity, you can use:
```
sId := &msp.SerializedIdentity{Mspid: mspID, IdBytes: certPEM}
raw, err := proto.Marshal(sId)
```

Hasing method returns byte array with the hash & error:

```
    Creatorhash, err = t.hashCreator(stub)
    if err != nil {
    	logger.Error(err)
        return shim.Error("Failed to hash creator")
    }
    logger.Infof("Creator hash: %s", fmt.Sprintf("%x", Creatorhash))
```

### Method

https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/interfaces.go#L166

Gets the `SignatureHeader.Creator` (e.g. an identity) of the agent (or user) submitting the transaction 
& hashes it with SHA256

```
func (t *SimpleChaincode) hashCreator(stub shim.ChaincodeStubInterface) ([]byte, error) {
	logger.Info("########### example_cc0 hashCreator ###########")

	Creatorbytes, err := stub.GetCreator()
	if err != nil {
		return nil, fmt.Errorf("Failed to get creator")
	}
	if Creatorbytes == nil {
		return nil, fmt.Errorf("Creator is not found")
	}

	Creatorhash, err := factory.GetDefault().Hash(Creatorbytes, &bccsp.SHA256Opts{})
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Failed computing SHA256 on [% x]", Creatorbytes))
	}
    logger.Infof("Creator hash: %x", Creatorhash)
    return Creatorhash, nil
}
```

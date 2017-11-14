var contracts = require('@monax/legacy-contracts');

// URL to the rpc endpoint of the Burrow server.
var burrowURL = "http://localhost:1337/rpc";
// See the 'Private Keys and Signing' section below for more info on this.
var accountData = require('/Users/acidumirae/.monax/chains/simplechain/accounts.json');
// newContractManagerDev lets you use an accountData object (address & private key) directly, i.e. no key/signing daemon is needed. This should only be used while developing/testing.
var contractManager = contracts.newContractManagerDev(burrowURL, accountData.simplechain_full_000);
// Create a new pipe.
//var pipe = new contracts.pipes.DevPipe(burrow, accountData.simplechain_full_000);
// Create a new contracts object using that pipe.
//var contractManager = contracts.newContractManager(pipe);

var myAbi = [{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"retVal","type":"uint256"}],"payable":false,"type":"function"}];
var myCompiledCode = "606060405260a18060106000396000f360606040526000357c01000000000000000000000000000000000000000000000000000000009004806360fe47b11460435780636d4ce63c14605d57603f565b6002565b34600257605b60048080359060200190919050506082565b005b34600257606c60048050506090565b6040518082815260200191505060405180910390f35b806000600050819055505b50565b60006000600050549050609e565b9056";

// Create a factory for the contract with the JSON interface 'myAbi'.
var myContractFactory = contractManager.newContractFactory(myAbi);

// To create a new instance and simultaneously deploy a contract use `new`:
var myNewContract;
myContractFactory.new({data: myCompiledCode}, function(error, contract){
    if (error) {
        // Something.
        throw error;
    }
    myNewContract = contract;
    console.log(contract);
});


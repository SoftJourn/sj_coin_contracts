// requires
var prompt = require('prompt');
var erisC = require('@monax/legacy-contracts');

// NOTE. On Windows/OSX do not use localhost. find the
// url of your chain with:
// docker-machine ls
// and find the docker machine name you are using (usually default or eris).
// for example, if the URL returned by docker-machine is tcp://192.168.99.100:2376
// then your erisdbURL should be http://192.168.99.100:1337/rpc
var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var idisContractAddress = "52DE1C4180124EFF078A718AB35E1DD2F1A4DE69";
var idisAbi = [{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"retVal","type":"uint256"}],"payable":false,"type":"function"}];

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('/Users/acidumirae/.monax/chains/simplechain/accounts.json');
var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

// properly instantiate the contract objects using the abi and address
var idisContract = contractsManager.newContractFactory(idisAbi).at(idisContractAddress);

// display the current value of idi's contract by calling
// the `get` function of idi's contract
function getValue(callback) {
  idisContract.get(function(error, result){
    if (error) { throw error }
    console.log("Idi's number is:\t\t\t");
    console.log(result.toString());
    callback();
  });
}

// prompt the user to change the value of idi's contract
function changeValue() {
  prompt.message = "What number should Idi make it?";
  prompt.delimiter = "\t";
  prompt.start();
  prompt.get(['value'], function (error, result) {
    if (error) { throw error }
    console.log(result.value);
    setValue(result.value)
  });
}

// using eris-contracts call the `set` function of idi's
// contract using the value which was recieved from the
// changeValue prompt
function setValue(value) {
  idisContract.set(value, function(error, result){
    if (error) { throw error }
    getValue(function(){});
  });
}

// run
getValue(changeValue);

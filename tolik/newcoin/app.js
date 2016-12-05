// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

// NOTE. On Windows/OSX do not use localhost. find the
// url of your chain with:
// docker-machine ls
// and find the docker machine name you are using (usually default or eris).
// for example, if the URL returned by docker-machine is tcp://192.168.99.100:2376
// then your erisdbURL should be http://192.168.99.100:1337/rpc
var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var contractData = require('./epm.json');
var coinContractAddress = contractData["deployCoin"];
var coinAbi = JSON.parse(fs.readFileSync("./abi/" + coinContractAddress));

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('./accounts.json');
var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

// properly instantiate the contract objects using the abi and address
var coinContract = contractsManager.newContractFactory(coinAbi).at(coinContractAddress);

function mintCoin(owner, amount) {
  coinContract.mint(owner, amount, function(error, result){
    if (error) { throw error }
    queryBalanceCoin(owner,function(){});
  });
}

function sendCoin(receiver, amount) {
  coinContract.transfer(receiver, amount, function(error, result){
    if (error) { throw error }
    queryBalanceCoin(receiver,function(){});
  });
}

function queryBalanceCoin(addr, callback) {
  coinContract.balanceOf(addr, function(error, result){
    if (error) { throw error }
    console.log("Coin balance is:\t\t\t" + result.toNumber());
    callback();
  });
}

function sendCoins() {
  prompt.message = "What number of coins should we send?";
  prompt.delimiter = "\t";
  prompt.start();
  prompt.get(['value'], function (error, result) {
    if (error) { throw error }
    sendCoin(accountData.simplechain_root_000.address,result.value);
    queryBalanceCoin(accountData.simplechain_full_000.address,function(){});
  });
}

function mintCoins() {
  prompt.message = "What number of coins should we mint?";
  prompt.delimiter = "\t";
  prompt.start();
  prompt.get(['value'], function (error, result) {
    if (error) { throw error }
    mintCoin(accountData.simplechain_full_000.address,result.value);
  });
}

function getInfo(callback) {
  coinContract.getInfo(function(error, result){
    if (error) { throw error }
    console.log(result);
    //console.log(JSON.parse(result));
    //console.log(result[0]);
    //console.log(result[1]);
    //console.log(result[2]);
    callback();
  });
}

function setInfo(callback) {
  coinContract.setInfo('SJ Coin','sjc',7,function(error, result){
    if (error) { throw error }
    console.log(result);
    callback();
  });
}

function getTokenName(callback) {
  console.log(coinContract.tokenName);
  console.log(coinContract.tokenSymbol);
  console.log(coinContract.tokenColor);
}


function runIssuerWallet() {
  //mintCoins();
  //sendCoins();
  //setInfo(function(){getInfo(function(){})});
  getTokenName(function(){});
  //queryBalanceCoin(accountData.simplechain_full_000.address,function(){});
  //queryBalanceCoin(accountData.simplechain_root_000.address,function(){});
}

// run
runIssuerWallet();

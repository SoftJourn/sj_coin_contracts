// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var contractData = require('./epm.json');
var coinContractAddress = contractData["deployCoin"];
//var coinContractAddress = "1362F2BD1FDF54543E82807673CF285B84BE0C55";
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
    console.log(addr + " balance is:\t" + result.toNumber());
    callback();
  });
}

function seedMoney() {
  mintCoin(accountData.simplechain_full_000.address, 100000);
  sendCoin(accountData.simplechain_root_000.address, 50000);
  sendCoin(accountData.simplechain_root_001.address, 50000);
}

function querySeed() {
  queryBalanceCoin(accountData.simplechain_full_000.address,function(){});
  queryBalanceCoin(accountData.simplechain_root_000.address,function(){});
  queryBalanceCoin(accountData.simplechain_root_001.address,function(){});
}

function runSeed() {
  //seedMoney();
  querySeed();
}

// run
runSeed();

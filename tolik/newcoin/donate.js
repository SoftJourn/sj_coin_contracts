// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var contractData = require('./epm.json');
var coinContractAddress = contractData["deployCoin"];
var coinAbi = JSON.parse(fs.readFileSync("./abi/" + coinContractAddress));

var donateData = require('../crowdsale/epm.json');
var donateContractAddr = donateData["deployCrowdsale"];

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('./accounts.json');
var who = accountData.simplechain_root_000;
var contractsManager = erisC.newContractManagerDev(erisdbURL, who);

// properly instantiate the contract objects using the abi and address
var coinContract = contractsManager.newContractFactory(coinAbi).at(coinContractAddress);

function donateCoin(sender, spender, amount) {
  coinContract.approveAndCall(spender, amount, function(error, result){
    if (error) { throw error }
    queryBalanceCoin(sender,function(){});
    queryBalanceCoin(spender,function(){});
  });
}

function queryBalanceCoin(addr, callback) {
  coinContract.balanceOf(addr, function(error, result){
    if (error) { throw error }
    console.log(addr + " balance is:\t" + result.toNumber());
    callback();
  });
}

function approveAllowance(from, to, amount) {
  coinContract.approve(to, amount, function(error, result) {
    if (error) { throw error }
    checkAllowance(from, to);
  });
}

function checkAllowance(from, to) {
  coinContract.allowance(from, to, function(error, result){
    if (error) { throw error }
    console.log(result.toNumber());
  });
}

function runIssuerWallet() {
  //approveAllowance(who.address, donateContractAddr, 5000);
  checkAllowance(who.address, donateContractAddr);
  //donateCoin(who.address, donateContractAddr, 4000);
}

// run
runIssuerWallet();

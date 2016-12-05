// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var contractData = require('../crowdsale/epm.json');
var saleContractAddress = contractData["deployCrowdsale"];
var saleAbi = JSON.parse(fs.readFileSync("../crowdsale/abi/" + saleContractAddress));

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('./accounts.json');
var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

// properly instantiate the contract objects using the abi and address
var saleContract = contractsManager.newContractFactory(saleAbi).at(saleContractAddress);

function setCoins(addressOfTokensAccumulated) {
  saleContract.setTokens(addressOfTokensAccumulated, function(error, result){
    if (error) { throw error }
    console.log("tokens loaded: " + result);
  });
}

function getInfo() {
  saleContract.beneficiary(function(error, result){
    console.log("Beneficiary: \t\t"+result);
  });  
  saleContract.creator(function(error, result){
    console.log("Creator: \t\t"+result);
  });
  saleContract.fundingGoal(function(error, result){
    console.log("Funding goal: \t\t"+result);
  });
  saleContract.amountRaised(function(error, result){
    console.log("Amount raised: \t\t"+result);
  });
  saleContract.deadline(function(error, result){
    console.log("Deadline: \t\t"+result);
  });
  /*
  saleContract.addressOfTokensAccumulated(function(error, result){
    console.log("Address of tokens accumulated: "+result);
  });
*/
  saleContract.getTokensCount(function(error, result){
    console.log("Tokens accumulated: \t"+result);
    for (i=0; i<result.toNumber(); i++) {
      saleContract.tokensAccumulated(i, function(error, result){
        console.log("\ttoken => " + result);
        saleContract.tokenAmounts(result, function(error, result){
          console.log("\t\tamount => " + result);
        });
      });
    }
  });
  /* [_from][_token]
  saleContract.balanceOf(function(error, result){
    console.log("balanceOf: \t\t"+result);
  });*/
}

function runCrowdsale() {
  var coins = new Array("51B242DB1EF7DF6682A1409E09BFF8A3933E8214","115348972793B1950D6C0D8E836F2C3828FCA478");
  //setCoins(coins);
  getInfo();
}

// run
runCrowdsale();

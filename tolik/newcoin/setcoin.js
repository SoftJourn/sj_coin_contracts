// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var contractData = require('../crowdsale/epm.json');
//var saleContractAddress = contractData["deployCrowdsale"];
//var saleContractAddress = '3ECC2E6CE35DBDFF9A6936367CD4A885FAA8AB5D';
var saleContractAddress = '224867A1E43D2308A3C5AD0AC58BF0684C9DC059';
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
  saleContract.closeOnGoalReached(function(error, result){
    console.log("Close on goal reached: \t"+result);
  });
  /*
  saleContract.addressOfTokensAccumulated(function(error, result){
    console.log("Address of tokens accumulated: "+result);
  });
*/
  saleContract.getTokensCount(function(error, result){
    console.log("Tokens accumulated: \t"+result);
    if (typeof result != "undefined")
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

function checkGoalReached() {
  saleContract.checkGoalReached(function(error, result){
    console.log("checkGoalReached: \t\t"+result);
  }); 
}

function withdraw() {
  saleContract.safeWithdrawal(function(error, result){
    console.log("safeWithdrawal: \t\t"+result);
  }); 
}

function runCrowdsale() {
  var coins = new Array("1362F2BD1FDF54543E82807673CF285B84BE0C55","02C58F28348774E53ACC58015C900068B9D0AFB8");
  //setCoins(coins);
  getInfo();
  //checkGoalReached();
  //withdraw();
}

// run
runCrowdsale();

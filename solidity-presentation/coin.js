// requires
var fs = require ('fs');
var prompt = require('prompt');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

// get the abi and deployed data squared away
var coinContractAddress = "D7F96BA1A8FC462981E58B8BABA4411FAA897A73";
var coinAbi = JSON.parse(fs.readFileSync("abi/Coin.abi"));

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('/Users/acidumirae/.eris/chains/simplechain/accounts.json');
var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

// properly instantiate the contract objects using the abi and address
var coinContract = contractsManager.newContractFactory(coinAbi).at(coinContractAddress);

function setCoins(addressOfTokensAccumulated) {
  coinContract.setTokens(addressOfTokensAccumulated, function(error, result){
    if (error) { throw error }
    console.log("tokens loaded: " + result);
  });
}

function setIssuer() {
  coinContract.isIssuer(accountData.simplechain_root_000.address, function(error, result){
    if (error) { throw error }
    console.log("Is issuer? " + result);
    coinContract.addIssuer(accountData.simplechain_root_000.address, function(error, result){
      if (error) { throw error }
      console.log("Set issuer!");
      coinContract.isIssuer(accountData.simplechain_root_000.address, function(error, result){
        if (error) { throw error }
        console.log("Is issuer? " + result);
      });
    });
  });
}

function setBalance() {
  coinContract.balanceOf(accountData.simplechain_root_000.address, function(error, result){
    if (error) { throw error }
    console.log("Balance? " + result);
    coinContract.issue(10000, function(error, result){
      if (error) { throw error }
      console.log("Issued 10000 coins!");
      coinContract.transfer(accountData.simplechain_root_000.address, 5000, function(error, result){
        if (error) { throw error }
        console.log("Transfered 5000 coins!");
        coinContract.balanceOf(accountData.simplechain_root_000.address, function(error, result){
          if (error) { throw error }
          console.log("Balance? " + result);
        });
      });
    });
  });
}

function runCoin() {
  //setIssuer();
  setBalance();
}

// run
runCoin();

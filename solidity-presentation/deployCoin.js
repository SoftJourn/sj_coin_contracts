var fs = require('fs');
var erisC = require('eris-contracts');
var solc = require('solc');

var erisdbURL = "http://localhost:1337/rpc";

//source of the smart contract
var contractSource = fs.readFileSync("./Coin.sol", 'utf8');
const compiled = solc.compile(contractSource, 1).contracts.Coin;
const abi = JSON.parse(compiled.interface);
// result keys
// [ 'assembly',
//   'bytecode',
//   'functionHashes',
//   'gasEstimates',
//   'interface', <== this is the ABI
//   'opcodes',
//   'runtimeBytecode',
//   'solidity_interface' ]

var accountData = require('/Users/acidumirae/.eris/chains/simplechain/accounts.json');
var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

contractsManager.newContractFactory(abi).new({data: compiled.bytecode}, function(error, contract){
  if (error) {
    return console.log(error);
  }
  // Called twice. First with the transactionHash, then after the transaction has been mined.
  if (contract.address) {
    console.log('Deployed to: ' + contract.address);
    console.log('ABI: ' + compiled.interface);
  } else {
    console.log('Waiting for mining...');
  }    
});

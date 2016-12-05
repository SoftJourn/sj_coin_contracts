// requires
var fs = require ('fs');

var erisC = require('eris-contracts');
var solc = require('solc');

var erisdbURL = "http://localhost:1337/rpc";

// properly instantiate the contract objects manager using the erisdb URL
// and the account data (which is a temporary hack)
var accountData = require('./accounts.json');

//source of the smart contract
var contractSource = fs.readFileSync("./Coin.sol", 'utf8');
const compiled = solc.compile(contractSource, 1).contracts['Coin'];
const abi = JSON.parse(compiled.interface);

var pipe = new erisC.pipes.DevPipe(erisdbURL, accountData.simplechain_full_000);
var contractManager = erisC.newContractManager(pipe); 
//var contractManager = erisC.newContractManagerDev(erisdbURL, accountData.simplechain_full_000);

// Create a factory for the contract with the JSON interface 'abi'.
const contractFactory = contractManager.newContractFactory(abi);

//console.log(compiled);

contractFactory.new({data: compiled.bytecode}, function(error, contract){
	if(error) {throw error}
	console.log(contract);
});


function test() {
	//contractFactory.at(contract.address, function(error, contract){
	//});
}

// run
test();
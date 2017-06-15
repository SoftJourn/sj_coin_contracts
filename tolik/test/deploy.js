// requires
var fs = require ('fs');

// Get '@monax/legacy-contracts'. 
var contracts = require('@monax/legacy-contracts');
 
// Get '@monax/legacy-db' (the javascript API for Burrow) 
var burrowModule = require("@monax/legacy-db");
 
// Create a new instance of Burrow that uses the given URL. 
var burrow = burrowModule.createInstance("http://localhost:1337/rpc");
// The private key. 
var accountData = require('./accounts.json');
 
// Create a new pipe. 
var pipe = new contracts.pipes.DevPipe(burrow, accountData.simplechain_full_000);
// Create a new contracts object using that pipe. 
var contractManager = contracts.newContractManager(pipe);

//source of the smart contract
const myJsonAbi = JSON.parse(fs.readFileSync("./MyContract.abi", 'utf8'));
const myCode = fs.readFileSync("./MyContract.bin", 'utf8');

// Create a factory (or contract template) from 'myJsonAbi' 
var myContractFactory = contractManager.newContractFactory(myJsonAbi);

var myContract;
 
myContractFactory.new(myCode, function(error, contract){
    if(error) {throw error}
    myContract = contract;
    console.log(contract);
});

//console.log(compiled);
//console.log(abi);

function test() {
	//contractFactory.at(contract.address, function(error, contract){
	//});
}

// run
test();

/*
Contract {
  address: 'A4609A127F40EFA9656B20637B21F1E6A72B2F3E',
  abi:
   [ { constant: true,
       inputs: [Object],
       name: 'add',
       outputs: [Object],
       payable: false,
       type: 'function' } ],
  add:
   { [Function: bound ]
     request: [Function: bound ],
     call: [Function: bound ],
     sendTransaction: [Function: bound ],
     estimateGas: [Function: bound ],
     getData: [Function: bound ],
     'int256,int256': [Circular] } }
*/
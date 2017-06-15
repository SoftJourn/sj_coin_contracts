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
const myJsonAbi = JSON.parse(fs.readFileSync("./abi/IdisContractsFTW.abi", 'utf8'));
const myCode = fs.readFileSync("./abi/IdisContractsFTW.bin", 'utf8');

// Create a factory (or contract template) from 'myJsonAbi' 
var myContractFactory = contractManager.newContractFactory(myJsonAbi);

var myContract;
 
myContractFactory.new(myCode, function(error, contract){
    if(error) {throw error}
    myContract = contract;
    console.log(contract);
});

//console.log(myCode);
//console.log(myJsonAbi);

function test() {
	//contractFactory.at(contract.address, function(error, contract){
	//});
}

// run
test();

/*
Contract {
  address: '8745C5D9611F960251890A0C692C2E2E4E9A5D52',
  abi:
   [ { constant: false,
       inputs: [Object],
       name: 'set',
       outputs: [],
       payable: false,
       type: 'function' },
     { constant: true,
       inputs: [],
       name: 'get',
       outputs: [Object],
       payable: false,
       type: 'function' } ],
  set:
   { [Function: bound ]
     request: [Function: bound ],
     call: [Function: bound ],
     sendTransaction: [Function: bound ],
     estimateGas: [Function: bound ],
     getData: [Function: bound ],
     uint256: [Circular] },
  get:
   { [Function: bound ]
     request: [Function: bound ],
     call: [Function: bound ],
     sendTransaction: [Function: bound ],
     estimateGas: [Function: bound ],
     getData: [Function: bound ],
     '': [Circular] } }
*/
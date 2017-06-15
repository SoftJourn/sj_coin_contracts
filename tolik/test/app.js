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

var address = "A4609A127F40EFA9656B20637B21F1E6A72B2F3E";
var myContract = myContractFactory.at(address);

var res;
 
try{
    res = myContract.add(3, 2, function (error,result){
        if (error) throw error;
        console.log(result);
        console.log(res);
    });
} catch (error) {
    console.log(error);
}
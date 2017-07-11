// requires
var fs = require ('fs');

// Get '@monax/legacy-contracts'. 
var contracts = require('@monax/legacy-contracts');
 
// Get '@monax/legacy-db' (the javascript API for Burrow) 
var burrowModule = require("@monax/legacy-db");
 
// Create a new instance of Burrow that uses the given URL. 
var burrow = burrowModule.createInstance("http://192.168.33.10:1337/rpc");
// The private key. 
var accountData = require('./accounts.json');
 
// Create a new pipe. 
var pipe = new contracts.pipes.DevPipe(burrow, accountData.multichain_full_000);
// Create a new contracts object using that pipe. 
var contractManager = contracts.newContractManager(pipe);

//source of the smart contract
const myJsonAbi = JSON.parse(fs.readFileSync("./abi/GSFactory", 'utf8'));

// Create a factory (or contract template) from 'myJsonAbi' 
var myContractFactory = contractManager.newContractFactory(myJsonAbi);

// 'monax pkgs do' output file
var json = JSON.parse(fs.readFileSync("./jobs_output.json",'utf8'));
var address = json.deployGSFactory;
var myContract = myContractFactory.at(address);

var res;
 
try{
    myContract.getLast(function (error,result){
        if (error) throw error;
        console.log(result);
    });
} catch (error) {
    console.log(error);
}

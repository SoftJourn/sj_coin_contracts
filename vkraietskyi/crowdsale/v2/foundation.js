let erisC = require('eris-contracts');
let fs = require('fs');
let rx = require('rxjs');
let http = require('http');
let accountData = require('./accounts.json');

let erisCompilerHost = '172.17.0.4';
let erisCompilerPort = 9099;
let erisdbURL = "http://172.17.0.3:1337/rpc";
let coinContractPath = "vkraietskyi/crowdsale/v2/Coin.sol";
const coinContractName = "Coin";
let foundationContractPath = "vkraietskyi/crowdsale/v2/Foundation.sol";
const foundationContractName = "Foundation";
let contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.sjcoins_full_000);

// read files
let coinContract = fs.readFileSync(coinContractPath, 'utf8');
let foundationContract = fs.readFileSync(foundationContractPath, 'utf8');

// deploy contracts
deploy(coinContract, coinContractName, 1).subscribe(result => {
    console.log(result);
    let foundation = accountData.sjcoins_full_000.address;
    let totalToCollect = 100;
    let deadline = 100;
    let closeOnGoalReached = true;
    let finalToken = result['address'];
    let addressOfTokensAccumulated = [result['address']];
    deploy(foundationContract, foundationContractName, foundation, totalToCollect, deadline, closeOnGoalReached, finalToken, addressOfTokensAccumulated).subscribe(result => {
        console.log(result);
    });
});

function stringToByteArray(content) {
    let buf = new ArrayBuffer(content.length * 2); // 2 bytes for each char
    let bufView = new Uint16Array(buf);
    for (let i = 0, strLen = content.length; i < strLen; i++) {
        bufView[i] = content.charCodeAt(i);
    }
    return bufView;
}

function getObjectNames(code) {
    let names = new Array();
    let pattern = /(?:contract)\s(\w+)/g;
    let matches = code.match(pattern);
    if (matches == null) {
        throw new Error("Contract name was not found!");
    } else {
        matches.forEach(function (match, index, array) {
            names.push(match.split(" ")[1]);
        });
        return names;
    }
}

//function deploys contracts
function deploy(contract, contractName, ...parameters) {
    return rx.Observable.create(observer => {
        let compiledContract = rx.Observable.create(observer => {
            getCompiledContract(contract, contractName).subscribe(result => {
                observer.next(result);
            });
        });
        compiledContract.subscribe(result => {
            contractsManager.newContractFactory(JSON.parse(result['abi'])).new(
                parameters,
                {data: result['bytecode']},
                function (error, contract) {
                    if (error) {
                        console.log(error);
                        throw error;
                    }
                    result['address'] = contract['address'];
                    observer.next(result);
                });
        });
    });
}
function getCompiledContract(contract, contractName) {
    return rx.Observable.create(observer => {
        let compiledCoin = rx.Observable.create(observer => {
            compile(contract).subscribe(result => {
                observer.next(result);
            });
        });
        compiledCoin.subscribe(result => {
            result['objects'].filter(value => {
                if (value['objectname'] == contractName) {
                    observer.next(value);
                }
            });
        });
    });
}
//function compile contracts
function compile(contract) {
    return rx.Observable.create(observer => {
        let objectNames = getObjectNames(contract);
        let codeInBytes = stringToByteArray(contract);
        let binstr = Array.prototype.map.call(codeInBytes, function (ch) {
            return String.fromCharCode(ch);
        }).join('');
        let coinCodeBase64 = new Buffer(binstr).toString('base64');
        let compileRequestObject = {
            name: "",
            language: "sol",
            includes: {"object.sol": {objectNames: objectNames, script: coinCodeBase64}},
            libraries: "",
            optimize: true,
            replacement: {
                "object.sol": "object.sol"
            }
        };
        compileRequest(compileRequestObject).subscribe(result => {
            observer.next(result);
        });
    });
}

function compileRequest(data) {
    return rx.Observable.create(observer => {
        let options = {
            hostname: erisCompilerHost,
            port: erisCompilerPort,
            path: '',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': JSON.stringify(data).length
            }
        };
        let result = '';
        let request = http.request(options, (response) => {
            response.setEncoding('utf8');
            response.on('data', chunk => {
                result += chunk;
            });
            response.on('end', () => {
                observer.next(JSON.parse(result));
            });
        });
        request.write(JSON.stringify(data));
    });
}



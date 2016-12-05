var fs = require('fs');
var erisC = require('eris-contracts');

var erisdbURL = "http://localhost:1337/rpc";

var idisAbi = JSON.parse(fs.readFileSync("./abi/Crowdsale", 'utf8'));

var accountData = require('./accounts.json');

var contractsManager = erisC.newContractManagerDev(erisdbURL, accountData.test_full_000);

var compiled = "60606040526000600760006101000a81548160ff02191690837f01000000000000000000000000000000000000000000000000000000000000009081020402179055506000600760016101000a81548160ff02191690837f0100000000000000000000000000000000000000000000000000000000000000908102040217905550604051610e6c380380610e6c833981016040528080519060200190919080519060200190919080519060200190919080518201919060200150505b83600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c0100000000000000000000000090810204021790555082600160005081905550603c820242016003600050819055508060046000509080519060200190828054828255906000526020600020908101928215610194579160200282015b828111156101935782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c010000000000000000000000009081020402179055509160200191906001019061013f565b5b5090506101db91906101a1565b808211156101d757600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016101a1565b5090565b50505b50505050610c7c806101f06000396000f360606040523615610095576000357c01000000000000000000000000000000000000000000000000000000009004806301cb3b201461009a57806329dcb0cf146100ae57806338af3eed146100d65780637a3a0e84146101145780637b3e5e7b1461013c5780638f4ffcb114610164578063e4e0ef35146101da578063f7888aec14610221578063fd6b7ef81461025b57610095565b610002565b34610002576100ac600480505061026f565b005b34610002576100c0600480505061038e565b6040518082815260200191505060405180910390f35b34610002576100e86004805050610397565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b346100025761012660048050506103bd565b6040518082815260200191505060405180910390f35b346100025761014e60048050506103c6565b6040518082815260200191505060405180910390f35b34610002576101d86004808035906020019091908035906020019091908035906020019091908035906020019082018035906020019191908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509090919050506103cf565b005b34610002576101f5600480803590602001909190505061070a565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3461000257610245600480803590602001909190803590602001909190505061074c565b6040518082815260200191505060405180910390f35b346100025761026d6004805050610777565b005b6003600050544210151561038b5760016000505460026000505410151561034b576001600760006101000a81548160ff02191690837f01000000000000000000000000000000000000000000000000000000000000009081020402179055507fec3f991caf7857d61663fd1bba1739e04abd4781238508cde554bb849d790c85600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260005054604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15b6001600760016101000a81548160ff02191690837f01000000000000000000000000000000000000000000000000000000000000009081020402179055505b5b5b565b60036000505481565b600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60016000505481565b60026000505481565b6000600060006000600760019054906101000a900460ff16156103f157610002565b60009350600192505b60046000508054905084101561049157600460005084815481101561000257906000526020600020900160005b9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff16141561048457600092508250610491565b83806001019450506103fa565b821561049c57610002565b3091508573ffffffffffffffffffffffffffffffffffffffff166323b872dd89848a600060405160200152604051847c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b156100025760325a03f1156100025750505060405180519060200150156106ff5786905080600560005060008a73ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008873ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005081905550600660005080548060010182818154818355818115116106155781836000526020600020918201910161061491906105f6565b8082111561061057600081815060009055506001016105f6565b5090565b5b5050509190906000526020600020900160005b8a909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c010000000000000000000000009081020402179055505080600260008282825054019250508190555061068061026f565b7f1fe43a085f5507370c77d08b78179e99076ad281426b5189351310be1b72625c8887836001604051808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390a15b5b5050505050505050565b600460005081815481101561000257906000526020600020900160005b9150909054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6005600050602052816000526040600020600050602052806000526040600020600091509150505481565b6000600060006000600060006000600060036000505442101515610c7157600760009054906101000a900460ff1615156109e857600097505b6006600050805490508810156109e357600660005088815481101561000257906000526020600020900160005b9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169650600560005060008873ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000509550600094505b6004600050805490508510156109de57600460005085815481101561000257906000526020600020900160005b9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1693508560008573ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005054925060008311156109d9578373ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8885600060405160200152604051837c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b156100025760325a03f1156100025750505060405180519060200150507f1fe43a085f5507370c77d08b78179e99076ad281426b5189351310be1b72625c8785856000604051808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390a15b610834565b6107b0565b610c6f565b600760009054906101000a900460ff168015610a5157503373ffffffffffffffffffffffffffffffffffffffff16600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b15610c6e5760009150600190505b600460005080549050821015610c6d57600460005082815481101561000257906000526020600020900160005b9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260005054600060405160200152604051837c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b156100025760325a03f115610002575050506040518051906020015015610c5b577f1fe43a085f5507370c77d08b78179e99076ad281426b5189351310be1b72625c600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600460005084815481101561000257906000526020600020900160005b9054906101000a900473ffffffffffffffffffffffffffffffffffffffff166002600050546000604051808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152602001821515815260200194505050505060405180910390a1610c60565b610002565b8180600101925050610a5f565b5b5b5b5b5b505050505050505056";

var myNewContract;


contractsManager.newContractFactory(idisAbi).new({data: compiled}, function(error, contract){
    if (error) {
        // Something.
        throw error;
    }
    myNewContract = contract;
    console.log(myNewContract);
});

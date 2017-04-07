import {Print} from "./print";
let erisContracts = require('eris-contracts');
let accountData = require('./accounts.json');
let contractData = require('./contracts.json');
let assert = require('assert');

let erisdbURL = "http://172.17.0.4:1337/rpc";

let contractsManager = erisContracts.newContractManagerDev(erisdbURL, accountData.sjcoins_full_000);

function newCoin(contractData: any, color: number) {
    return new Promise<any>(resolve => {
        contractsManager.newContractFactory(contractData.abi).new(
            color,
            {data: contractData.bytecode},
            function (error, contract) {
                if (error) {
                    console.log(error);
                    throw error;
                }
                resolve(contract);
            });
    });
}

function newFoundation(contractData: any, foundation: string, fundingGoal: number, deadline: number, closeOnGoalReached: boolean, mainToken: string, foundationContract: Array<string>) {
    return new Promise<any>(resolve => {
        contractsManager.newContractFactory(contractData.abi).new(
            foundation, fundingGoal, deadline, closeOnGoalReached, mainToken, foundationContract,
            {data: contractData.bytecode},
            function (error, contract) {
                if (error) {
                    console.log(error);
                    throw error;
                }
                resolve(contract);
            });
    });
}

function mint(contract: any, address: string, amount: number) {
    return new Promise<void>(resolve => {
        contract.mint(address, amount, function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            resolve();
        })
    });
}

function approveAndCall(contract: any, address: string, amount: number) {
    return new Promise<void>(resolve => {
        contract.approveAndCall(address, amount, function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            resolve(result);
        })
    });
}

function close(foundationContract: any) {
    return new Promise<void>(resolve => {
        foundationContract.close(function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            resolve(result);
        })
    });
}
function withdraw(foundationContract: any, amount: number, id: number, note: string) {
    return new Promise<void>(resolve => {
        foundationContract.withdraw(amount, id, note, function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            resolve(result);
        })
    });
}

function sleep(seconds) {
    return new Promise(resolve => setTimeout(resolve, seconds * 1000));
}

async function normalFlow() {
    console.log("NORMAL FLOW TEST");
    // deploy coins contract
    let coinContracts = [];
    for (let i = 0; i < 2; i++) {
        let contract = await newCoin(contractData[0], i);
        coinContracts.push(contract);
    }
    await Print.print("Coins contracts have been deployed successful");
    // preparing parameters for foundation contract
    let coinAddresses = coinContracts.map(value => {
        return value.address
    });
    let foundation = accountData.sjcoins_full_000.address;
    let totalToCollect = 200;
    let lifeTime = 1;
    let closeOnGoalReached = true;
    let finalToken = coinAddresses[0];
    // deploy foundation contract
    let foundationContract = await newFoundation(contractData[1], foundation, totalToCollect, lifeTime, closeOnGoalReached, finalToken, coinAddresses);
    await Print.printVSSeparator("Foundation contract have been deployed successful");
    // print contract variables to be sure that all was deployed fine
    await Print.printValueOf(foundationContract, "foundation", "Foundation account address:", "string");
    await Print.printValueOf(foundationContract, "fundingGoal", "Foundation coins to collect:", "number");
    await Print.printValueOf(foundationContract, "deadline", "Foundation contract's deadline:", "date");
    await Print.printValueOf(foundationContract, "closeOnGoalReached", "Foundation close on goal reached condition", "bool");
    await Print.printValueOf(foundationContract, "mainToken", "Foundation main coin", "string");
    await Print.printValueOf(foundationContract, "getTokens", "Foundation allowed coins", "string");
    // mint some coins
    await mint(coinContracts[0], foundation, 200);
    await mint(coinContracts[1], foundation, 200);
    // print balances
    await Print.printBalance(coinContracts[0], foundation, "Foundation account balance:");
    await Print.printBalance(coinContracts[1], foundation, "Foundation account balance:");
    // donate some coins
    let firstDonate = await approveAndCall(coinContracts[0], foundationContract.address, 200);
    assert.equal(true, firstDonate);
    let secondDonate = await approveAndCall(coinContracts[1], foundationContract.address, 200);
    assert.equal(true, secondDonate);
    await Print.printVSSeparator("Donate to foundation using first coin result:" + firstDonate);
    await Print.printVSSeparator("Donate to foundation using second coin result:" + secondDonate);
    // print balances to be sure that coins were transferred fine
    await Print.printBalance(coinContracts[0], foundation, "Foundation account balance:");
    await Print.printBalance(coinContracts[1], foundation, "Foundation account balance:");
    // print amount of collected coins
    await Print.printValueOf(foundationContract, "amountRaised", "Foundation coins collected:", "number");
    // print contract balances to be sure that donates were transferred fine
    await Print.printBalance(coinContracts[0], foundationContract.address, "Foundation contract balance:");
    await Print.printBalance(coinContracts[1], foundationContract.address, "Foundation contract balance:");
    // closing contract
    let closeResult = await close(foundationContract);
    assert.equal(true, closeResult);
    await Print.printVSSeparator("Close foundation contract result:" + closeResult);
    // print all amount of coins that were collected before contract was closed
    await Print.printValueOf(foundationContract, "contractRemains", "Foundation coins to exchange:", "number");
    // mint more coins to perform exchange
    await mint(coinContracts[0], foundation, 400);
    // check the amount of minted coins
    await Print.printBalance(coinContracts[0], foundation, "Foundation account balance: ");
    // exchange coins
    let exchangeResult = await approveAndCall(coinContracts[0], foundationContract.address, 400);
    assert.equal(true, exchangeResult);
    await Print.printVSSeparator("Exchange all coins result: " + exchangeResult);
    // check balances after exchange
    await Print.printBalance(coinContracts[0], foundation, "Foundation account balance: ");
    await Print.printBalance(coinContracts[1], foundation, "Foundation account balance: ");
    await Print.printBalance(coinContracts[0], foundationContract.address, "Foundation contract balance:");
    // withdraw coins
    let withdrawResult = await withdraw(foundationContract, 400, 1, "Brought 400 coins");
    assert.equal(true, withdrawResult);
    await Print.printVSSeparator("Withdraw foundation contract result: " + withdrawResult);
    // check contract balance to be sure that withdrawal was performed fine
    await Print.printBalance(coinContracts[0], foundationContract.address, "Foundation contract balance:");
    await Print.printBalance(coinContracts[0], foundation, "Foundation account balance: ");
    await Print.printBalance(coinContracts[1], foundation, "Foundation account balance: ");
    // get data about withdrawal operations
    await Print.printItemsOf(foundationContract, "getContractFulfilmentRecord", "getContractFulfilmentRecordLength", "Fulfilment records");
}

async function donateAfterDeadline() {
    console.log("DONATE AFTER DEADLINE TEST");
    // deploy coins contract
    let coinContract = await newCoin(contractData[0], 1);
    Print.print("Coins contracts have been deployed successful");
    // preparing parameters for foundation contract
    let foundation = accountData.sjcoins_full_000.address;
    let totalToCollect = 200;
    let lifeTime = 1; // 1 minute
    let closeOnGoalReached = true;
    let finalToken = coinContract.address;
    // deploy foundation contract
    let foundationContract = await newFoundation(contractData[1], foundation, totalToCollect, lifeTime, closeOnGoalReached, finalToken, [coinContract.address]);
    await Print.printVSSeparator("Foundation contract have been deployed successful");
    // print contract variables to be sure that all was deployed fine
    await Print.printValueOf(foundationContract, "foundation", "Foundation account address:", "string");
    await Print.printValueOf(foundationContract, "fundingGoal", "Foundation coins to collect:", "number");
    await Print.printValueOf(foundationContract, "deadline", "Foundation contract's deadline:", "date");
    await Print.printValueOf(foundationContract, "closeOnGoalReached", "Foundation close on goal reached condition", "bool");
    await Print.printValueOf(foundationContract, "mainToken", "Foundation main coin", "string");
    await Print.printValueOf(foundationContract, "getTokens", "Foundation allowed coins", "string");
    // mint some coins
    await mint(coinContract, foundation, 200);
    // print balances
    await Print.printBalance(coinContract, foundation, "Foundation account balance:");
    await Print.printVSSeparator("Waiting for deadline");
    await sleep(60);
    let lateDonate = await approveAndCall(coinContract, foundationContract.address, 100);
    assert.equal(false, lateDonate);
    await Print.printVSSeparator("Donate to foundation using second2 coin result: " + lateDonate);
}

async function withdrawalBeforeContractWasClosed() {
    console.log("WITHDRAWAL BEFORE CONTRACT WAS CLOSED TEST");
    // deploy coins contract
    let coinContract = await newCoin(contractData[0], 1);
    Print.print("Coins contracts have been deployed successful");
    // preparing parameters for foundation contract
    let foundation = accountData.sjcoins_full_000.address;
    let totalToCollect = 200;
    let lifeTime = 1; // 1 minute
    let closeOnGoalReached = true;
    let finalToken = coinContract.address;
    // deploy foundation contract
    let foundationContract = await newFoundation(contractData[1], foundation, totalToCollect, lifeTime, closeOnGoalReached, finalToken, [coinContract.address]);
    await Print.printVSSeparator("Foundation contract have been deployed successful");
    // print contract variables to be sure that all was deployed fine
    await Print.printValueOf(foundationContract, "foundation", "Foundation account address:", "string");
    await Print.printValueOf(foundationContract, "fundingGoal", "Foundation coins to collect:", "number");
    await Print.printValueOf(foundationContract, "deadline", "Foundation contract's deadline:", "date");
    await Print.printValueOf(foundationContract, "closeOnGoalReached", "Foundation close on goal reached condition", "bool");
    await Print.printValueOf(foundationContract, "mainToken", "Foundation main coin", "string");
    await Print.printValueOf(foundationContract, "getTokens", "Foundation allowed coins", "string");
    // print balances
    await Print.printVSSeparator("Withdraw before contract was closed");
    let withdrawResult = await withdraw(foundationContract, 400, 1, "Brought 400 coins");
    assert.equal(false, withdrawResult);
}

function main() {
    // normalFlow();
    // donateAfterDeadline();
    withdrawalBeforeContractWasClosed();
}

main();
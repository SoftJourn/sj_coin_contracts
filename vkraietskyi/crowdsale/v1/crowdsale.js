let erisC = require('eris-contracts');

let accountData = require('./accounts.json');
let erisdbURL = "http://46.101.203.71:1337/rpc";

// create contract manager
let fullContractManager = erisC.newContractManagerDev(erisdbURL, accountData.sj_coins_full_000);
let meContractManager = erisC.newContractManagerDev(erisdbURL, accountData.me);
let crowdAbi = [{"constant":false,"inputs":[],"name":"checkGoalReached","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[],"name":"creator","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[],"name":"deadline","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"beneficiary","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[],"name":"getTokensCount","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"_token","type":"address"}],"name":"isTokenAccumulated","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[],"name":"fundingGoal","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"amountRaised","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"},{"name":"_token","type":"address"},{"name":"_extraData","type":"bytes"}],"name":"receiveApproval","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"tokenAmounts","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"getTokens","outputs":[{"name":"","type":"address[]"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"donators","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"getNow","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"getDonators","outputs":[{"name":"","type":"address[]"}],"type":"function"},{"constant":true,"inputs":[],"name":"closeOnGoalReached","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"tokensAccumulated","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"inputs":[{"name":"ifSuccessfulSendTo","type":"address"},{"name":"fundingGoalInTokens","type":"uint256"},{"name":"durationInMinutes","type":"uint256"},{"name":"onGoalReached","type":"bool"},{"name":"addressOfTokensAccumulated","type":"address[]"}],"type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"beneficiary","type":"address"},{"indexed":false,"name":"amountRaised","type":"uint256"}],"name":"GoalReached","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"backer","type":"address"},{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"isContribution","type":"bool"}],"name":"FundTransfer","type":"event"}];
let coinAbi = [{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"accounts","type":"address[]"},{"name":"amount","type":"uint256"}],"name":"distribute","outputs":[{"name":"success","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"_tokenColor","type":"uint8"}],"name":"setColor","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"tokenColor","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"owner","type":"address"},{"name":"amount","type":"uint256"}],"name":"mint","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"getColor","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"type":"function"},{"inputs":[{"name":"_tokenColor","type":"uint8"}],"type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}];

let crowd = "B16A0C8A1F463204B48DA96D3358FFCDA74BFF64";
let coin = "3DEFC7DDD478F01D23E3C210E9B47D8C9DBB5392";

let donators;
let tokens;

let fullCoinContract = fullContractManager.newContractFactory(coinAbi).at(coin);
let fullNewCrowd = fullContractManager.newContractFactory(crowdAbi).at(crowd);
let meCoinContract = meContractManager.newContractFactory(coinAbi).at(coin);
let meNewCrowd = meContractManager.newContractFactory(crowdAbi).at(crowd);

function main() {
    fullFundingGoal(function (error, result) {
        console.log("Full FundingGoal Error", error);
        console.log("Full FundingGoal Result", parseInt(result));
        console.log("-----------------------------------");
    });
    fullDeadline(function (error, result) {
        console.log("Full deadline Error", error);
        console.log("Full deadline timestamp Result", parseInt(result));
        console.log("Full deadline date Result", new Date(parseInt(result)*1000));
        console.log("-----------------------------------");
    });
    fullGetNow(function (error, result) {
        console.log("Full now Error", error);
        console.log("Full now timestamp Result", parseInt(result));
        console.log("Full now date Result", new Date(parseInt(result)*1000));
        console.log("-----------------------------------");
    });
    fullDonate(crowd, 1, function (error, result) {
        console.log("Full account donate Error", error);
        console.log("Full account donate Result", result);
        console.log("-----------------------------------");
    });
    meDonate(crowd, 1, function (error, result) {
        console.log("Me account donate Error", error);
        console.log("Me account donate Result", result);
        console.log("-----------------------------------");
    });
    getTokens(function (error, result) {
        console.log("Tokens Error", error);
        console.log("Tokens Result", result);
        if (result) {
            tokens = result
        }
        console.log("-----------------------------------");
    });
    getDonators(function (error, result) {
        console.log("Donators Error", error);
        console.log("Donators Result", result);
        if (result) {
            donators = result
        }
        console.log("-----------------------------------");
        tokens.forEach(function (token, ti, arr) {
            donators.forEach(function (donator, di, arr) {
                getDonates(donator, token, function (error, result) {
                    console.log("Donator ", donator);
                    console.log("Token ", token);
                    console.log("Donate Error", error);
                    console.log("Donate Result", parseInt(result));
                    console.log("-----------------------------------");
                });
            });
        });
        meWithdraw(function (error, result) {
            console.log("Me Withdraw Error", error);
            console.log("Me Withdraw Result", parseInt(result));
            console.log("-----------------------------------");
            fullWithdraw(function (error, result) {
                console.log("Full Withdraw Error", error);
                console.log("Full Withdraw Result", parseInt(result));
                console.log("-----------------------------------");
            });

            fullCloseOnGoalReached(function (error, result) {
                console.log("Full СloseOnGoalReached Error", error);
                console.log("Full СloseOnGoalReached Result", result);
                console.log("-----------------------------------");
            });
        });
    });
}

function fullDonate(crowd, amount, callback) {
    fullCoinContract.approveAndCall(crowd, amount, function (error, result) {
        callback(error, result);
    });
}

function meDonate(crowd, amount, callback) {
    meCoinContract.approveAndCall(crowd, amount, function (error, result) {
        callback(error, result);
    });
}

function getDonators(callback) {
    fullNewCrowd.getDonators(function (error, result) {
        callback(error, result);
    });
}
function getTokens(callback) {
    fullNewCrowd.getTokens(function (error, result) {
        callback(error, result);
    });
}

function getDonates(token, donator, callback) {
    fullNewCrowd.balanceOf(token, donator, function (error, result) {
        callback(error, result);
    });
}

function meWithdraw(callback) {
    meNewCrowd.safeWithdrawal(function (error, result) {
        callback(error, result);
    });
}

function fullWithdraw(callback) {
    fullNewCrowd.safeWithdrawal(function (error, result) {
        callback(error, result);
    });
}
function fullDeadline(callback) {
    fullNewCrowd.deadline(function (error, result) {
        callback(error, result);
    });
}
function fullCloseOnGoalReached(callback) {
    fullNewCrowd.closeOnGoalReached(function (error, result) {
        callback(error, result);
    });
}

function fullFundingGoal(callback) {
    fullNewCrowd.fundingGoal(function (error, result) {
        callback(error, result);
    });
}

function fullGetNow(callback) {
    fullNewCrowd.getNow(function (error, result) {
        callback(error, result);
    });
}

main();
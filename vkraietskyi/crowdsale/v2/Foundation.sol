contract Token {
    event Transfer(address from, address to, uint value);

    event Approval(address from, address to, uint value);

    function transferFrom(address from, address to, uint amount) returns (bool success);

    function transfer(address receiver, uint amount) returns (bool);
}


contract Foundation {

    struct Detail {
    uint amount;
    uint id;
    uint time;
    string note;
    }

/*Contract's founder address*/
    address public foundation;

/*Contract's creator address*/
    address public creator;

/*Amount of coins to collect*/
    uint public fundingGoal;
/*Amount of coins which were collected before contract has been closed*/
    uint public amountRaised;

/*Amount of coins which were collected after contract has been closed*/
    uint public contractRemains = uint(0x0);

/*Details about contract performance*/
    Detail[] public contractFulfilmentRecord;

/*Token address into which should be exchanged all other tokens*/
    address public mainToken;

/*Contract's deadline(timestamp)*/
    uint public deadline;

/*Condition of contract closing*/
    bool public closeOnGoalReached;

/*Array of token addresses which are allowed for contract*/
    address[] public tokensAccumulated;

/*Array of account addresses of donators*/
    address[] public donators;

/*Status of contract*/
    bool fundingGoalReached = false;

/*Is contract closed*/
    bool crowdsaleClosed = false;

/*Map that contains amount of coins which depend on token address and donator address*/
    mapping (address => mapping (address => uint)) public balanceOf;

/*Map that contains amount of coins which depend on token address*/
    mapping (address => uint) public tokenAmounts;

    event GoalReached(address foundation, uint amountRaised);

    event FundTransfer(address backer, address token, uint amount, bool isContribution);

/* data structure to hold information about campaign contributors */
/*  at initialization, setup the owner */
    function Foundation(
    address ifSuccessfulSendTo,
    uint fundingGoalInTokens,
    uint durationInMinutes,
    bool onGoalReached,
    address finalToken,
    address[] addressOfTokensAccumulated) {
        creator = msg.sender;
        foundation = ifSuccessfulSendTo;
        fundingGoal = fundingGoalInTokens;
        deadline = now + durationInMinutes * 1 minutes;
        closeOnGoalReached = onGoalReached;
        mainToken = finalToken;
        setTokens(addressOfTokensAccumulated);
    }

/* You must run this once only with the same tokens or else :) */
    function setTokens(address[] addressOfTokensAccumulated) internal returns (uint) {
        if (msg.sender != foundation && msg.sender != creator) return 0;
    //tokensAccumulated = new address[](0);
        uint keyIndex = 0;
        while (keyIndex < addressOfTokensAccumulated.length) {
            address token = addressOfTokensAccumulated[keyIndex];
            tokensAccumulated.push(token);
            tokenAmounts[token] = 0;
            keyIndex++;
        }
        return keyIndex;
    }
/*---------------------------------------------Public methods---------------------------------------------------------*/

/*Method returns count of tokens that are allowed in contract*/
    function getTokensCount() constant returns (uint) {
        return tokensAccumulated.length;
    }

/*Method shows is token allowed to be used in contract*/
    function isTokenAccumulated(address _token) constant returns (bool) {
        if (crowdsaleClosed) return false;
        uint keyIndex = 0;
        while (keyIndex < tokensAccumulated.length) {
            if (_token == tokensAccumulated[keyIndex]) {
                return true;
            }
            keyIndex++;
        }
        return false;
    }

/* The function that gets donators addresses */
    function getDonators() returns (address[]){
        return donators;
    }

/* The function that gets token addresses */
    function getTokens() returns (address[]){
        return tokensAccumulated;
    }

/* The function without name is the default function that is called whenever anyone sends funds to a contract */
    function receiveApproval(address _from, uint _value, address _token, bytes _extraData) returns (bool) {
    /* If goal is not reached then - donate */
        if (!fundingGoalReached) {
        /* If it is not too late */
            if (now <= deadline) {
                return donate(_from, _value, _token);
            }
            else {
                return false;
            }
        }
        /* If goal is reached then - exchange all collected tokens into one token*/
        else if (fundingGoalReached && mainToken == _token) {
            if (contractRemains == _value) {
                return exchange(_from, _value, _token);
            }
            else {
                return false;
            }
        }
        else {
            return false;
        }
    }

/* checks if the goal or time limit has been reached and ends the campaign */
    function checkGoalReached() returns (bool) {
        if (now >= deadline || closeOnGoalReached) {
            if (amountRaised >= fundingGoal) {
                fundingGoalReached = true;
                GoalReached(foundation, amountRaised);
            }
            if (now >= deadline) {
                crowdsaleClosed = true;
            }
            if (fundingGoalReached) {
                crowdsaleClosed = true;
            }
        }
        return crowdsaleClosed;
    }

/*Method closes contract depends on conditions*/
    function close() returns (bool) {
    /* Do not allow to withdraw anything util crowdsale is closed */
        if (!checkGoalReached()) return false;
        if (msg.sender != creator) return false;
        uint keyIndex;
        uint donatorsIndex;
        address donator;
        address token;
        uint amount;
    /* if funding goal is not reached, then donator can withdraw its donation after deadline ;) */
        if (!fundingGoalReached) {
            donatorsIndex = 0;
        /* run through donators*/
            while (donatorsIndex < donators.length) {
                keyIndex = 0;
                donator = donators[donatorsIndex];
            /* run through tokens*/
                while (keyIndex < tokensAccumulated.length) {
                    token = tokensAccumulated[keyIndex];
                /* if token address exists and it is posible to get balance of donator by this token*/
                    if (token != address(0x0) && balanceOf[donator][token] != uint(0x0)) {
                        amount = balanceOf[donator][token];
                        if (amount > 0) {
                            if (Token(token).transfer(donator, amount)) {
                                FundTransfer(donator, token, amount, false);
                            }
                            else {
                                balanceOf[donator][token] = amount;
                            }
                        }
                    }
                    keyIndex++;
                }
                donatorsIndex++;
            }
            return true;
        }
    /* if funding goal is reached - save amount of collected coins */
        if (fundingGoalReached) {
            contractRemains = amountRaised;
            return true;
        }
        return false;
    }

/*Method withdraws coins back to foundation and records this withdrawal*/
    function withdraw(uint amount, uint id, string note) returns (bool) {
        if (amount <= contractRemains && !(amount > contractRemains) && contractRemains != uint(0x0)) {
            if (Token(mainToken).transfer(foundation, amount)) {
                contractRemains -= amount;
                contractFulfilmentRecord.push(Detail(amount, id, now, note));
                return true;
            }
            else {
                return false;
            }
        }
        else {
            return false;
        }
    }
/*Method returns count of withdrawal records*/
    function getContractFulfilmentRecordLength() constant returns (uint){
        return contractFulfilmentRecord.length;
    }

/*Method returns withdrawal record*/
    function getContractFulfilmentRecord(uint index) constant returns (uint, uint, uint, string){
        var data = contractFulfilmentRecord[index];
        return (data.amount, data.id, data.time, data.note);
    }

/*-----------------------------------------Private methods------------------------------------------------------------*/
/*Method transfers coins from account to contract*/
    function donate(address _from, uint _value, address _token) private returns (bool){
        if (!isTokenAccumulated(_token)) return false;
        address _to = this;
        if (!Token(_token).transferFrom(_from, _to, _value)) {
            return false;
        }
        if (balanceOf[_from][_token] == uint(0x0)) {
            balanceOf[_from][_token] = _value;
            donators.push(_from);
        }
        else {
            balanceOf[_from][_token] += _value;
        }
        tokenAmounts[_token] += _value;
        amountRaised += _value;
        FundTransfer(_from, _token, _value, true);
        return true;
    }

/*Method takes amount of mainTokens and transfer back amount of all other tokens*/
    function exchange(address _from, uint _value, address _token) private returns (bool) {
        uint keyIndex;
        address token;
        uint amount;
        keyIndex = 0;
        address _to = this;
    /* Transfer coins to this contract*/
        if (!Token(_token).transferFrom(_from, _to, _value)) {
            return false;
        }
    /* Transfer all kind of coins to foundation*/
        while (keyIndex < tokensAccumulated.length) {
            if (tokensAccumulated[keyIndex] != address(0x0)) {
                token = tokensAccumulated[keyIndex];
                amount = tokenAmounts[token];
                if (Token(token).transfer(foundation, amount)) {
                    FundTransfer(foundation, tokensAccumulated[keyIndex], amount, false);
                    delete tokenAmounts[token];
                }
                else {
                    throw;
                // Hopefully throw will roll back anything we sent to the blockchain so far
                }
            }
            keyIndex++;
        }
        return true;
    }
}

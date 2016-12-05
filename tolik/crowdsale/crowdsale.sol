contract Crowdsale {
  address public beneficiary;
  address public creator;
  uint public fundingGoal; uint public amountRaised; uint public deadline;
  address[] public tokensAccumulated;
  mapping(address => mapping(address => uint)) public balanceOf;
  mapping(address => uint) public tokenAmounts;
  bool fundingGoalReached = false;
  event GoalReached(address beneficiary, uint amountRaised);
  event FundTransfer(address backer, address token, uint amount, bool isContribution);
  bool crowdsaleClosed = false;
  /* data structure to hold information about campaign contributors */
  /*  at initialization, setup the owner */
  function Crowdsale(
      address ifSuccessfulSendTo,
      uint fundingGoalInTokens,
      uint durationInMinutes
  ) {
    creator = msg.sender;
    beneficiary = ifSuccessfulSendTo;
    fundingGoal = fundingGoalInTokens;
    deadline = now + durationInMinutes * 1 minutes;
  }
  function setTokens(address[] addressOfTokensAccumulated) returns (uint) {
    if (msg.sender != beneficiary && msg.sender != creator) throw;
    uint keyIndex = 0;
    while (keyIndex < addressOfTokensAccumulated.length) {
      address token = addressOfTokensAccumulated[keyIndex];
      tokensAccumulated.push(token);
      tokenAmounts[token] = 0;
      keyIndex++;
    }
    return keyIndex;
  }
  function getTokensCount() returns (uint) {
    return tokensAccumulated.length;
  }
  /* The function without name is the default function that is called whenever anyone sends funds to a contract */
  function receiveApproval(address _from, uint _value, address _token, bytes _extraData) {
    if (crowdsaleClosed) throw;
    uint keyIndex = 0;
    bool notFound = true;
    while (keyIndex < tokensAccumulated.length) {
      if (_token == tokensAccumulated[keyIndex]) {
        notFound = false;
        break;
      }
      keyIndex++;
    }
    if (notFound) throw;
    address _to = this;
    if(!_token.call(bytes4(bytes32(sha3("transferFrom(address,address,uint)"))), _from, _to, _value)) {
      throw;
    }
    if (balanceOf[_from][_token] == uint(0x0)) {
      balanceOf[_from][_token] = _value;
    } else {
      balanceOf[_from][_token] += _value;
    }
    tokenAmounts[_token] += _value;
    amountRaised += _value;
    FundTransfer(_from, _token, _value, true);
  }
  /* checks if the goal or time limit has been reached and ends the campaign */
  function checkGoalReached() {
    if (now >= deadline) {
      if (amountRaised >= fundingGoal){
          fundingGoalReached = true;
          GoalReached(beneficiary, amountRaised);
      }
      crowdsaleClosed = true;
    }
  }
  function safeWithdrawal() {
    if (now < deadline) throw;
    uint keyIndex;
    address coin;
    uint amount;
    /* if funding goal is not reached, then donator can withdraw its donation after deadline ;) */
    if (!fundingGoalReached) {
      keyIndex = 0;
      while (keyIndex < tokensAccumulated.length) {
        coin = tokensAccumulated[keyIndex];
        if (coin != address(0x0) && balanceOf[msg.sender][coin] != uint(0x0)) {
          amount = balanceOf[msg.sender][coin];
          balanceOf[msg.sender][coin] = 0;
          if (amount > 0) {
            if (coin.call(bytes4(bytes32(sha3("transfer(address,uint)"))), msg.sender, amount)) {
              FundTransfer(msg.sender, coin, amount, false);
            } else {
              balanceOf[msg.sender][coin] = amount;
            }
          }            
        }
        keyIndex++;
      }
    }
    /* if funding goal is reached then beneficiary can withdraw everything */
    if (fundingGoalReached && beneficiary == msg.sender) {
      bool notFound = true;
      keyIndex = 0;
      while (keyIndex < tokensAccumulated.length) {
        if (tokensAccumulated[keyIndex] != address(0x0)) {
          coin = tokensAccumulated[keyIndex];
          amount = tokenAmounts[coin];
          if (coin.call(bytes4(bytes32(sha3("transfer(address,uint)"))), beneficiary, amount)) {
            FundTransfer(beneficiary, tokensAccumulated[keyIndex], amountRaised, false);
            delete tokenAmounts[coin];
          } else {
            //If we fail to send the funds to beneficiary, unlock funders balance
            //fundingGoalReached = false;
            /* WUT??? */
          }            
        }
        keyIndex++;
      }
    }
  }
}
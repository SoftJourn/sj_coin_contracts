pragma solidity ^0.4.4;
contract Token {
  event Transfer(address from, address to, uint value);
  event Approval(address from, address to, uint value);
  function transferFrom(address from, address to, uint amount) returns (bool success);
  function transfer(address receiver, uint amount) returns (bool);
}
contract Crowdsale {
  address public beneficiary;
  address public creator;
  uint public fundingGoal; uint public amountRaised; uint public deadline; bool public closeOnGoalReached;
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
      uint durationInMinutes,
      bool onGoalReached
  ) {
    creator = msg.sender;
    beneficiary = ifSuccessfulSendTo;
    fundingGoal = fundingGoalInTokens;
    deadline = now + durationInMinutes * 1 minutes;
    closeOnGoalReached = onGoalReached;
  }
  /* You must run this once only with the same tokens or else :) */
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
  function isTokenAccumulated(address _token) returns (bool) {
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
  /* The function without name is the default function that is called whenever anyone sends funds to a contract */
  function receiveApproval(address _from, uint _value, address _token, bytes _extraData) returns (bool)
 {
    if (!isTokenAccumulated(_token)) return false;
    address _to = this;
    if(!Token(_token).transferFrom(_from, _to, _value)) {
      return false;
    }
    if (balanceOf[_from][_token] == uint(0x0)) {
      balanceOf[_from][_token] = _value;
    } else {
      balanceOf[_from][_token] += _value;
    }
    tokenAmounts[_token] += _value;
    amountRaised += _value;
    FundTransfer(_from, _token, _value, true);
    return true;
  }
  /* checks if the goal or time limit has been reached and ends the campaign */
  function checkGoalReached() returns (bool) {
    if (now >= deadline || closeOnGoalReached) {
      if (amountRaised >= fundingGoal){
          fundingGoalReached = true;
          GoalReached(beneficiary, amountRaised);
      }
      if (now >= deadline) {
        crowdsaleClosed = true;
      }
    }
    return crowdsaleClosed;
  }
  function safeWithdrawal() returns (bool) {
    /* Do not allow to withdraw anything util crowdsale is closed */
    if (!checkGoalReached()) return false;
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
            if (Token(coin).transfer(msg.sender, amount)) {
              FundTransfer(msg.sender, coin, amount, false);
            } else {
              balanceOf[msg.sender][coin] = amount;
            }
          }            
        }
        keyIndex++;
      }
      return true;
    }
    /* if funding goal is reached then beneficiary can withdraw everything */
    if (fundingGoalReached && beneficiary == msg.sender) {
      bool notFound = true;
      keyIndex = 0;
      while (keyIndex < tokensAccumulated.length) {
        if (tokensAccumulated[keyIndex] != address(0x0)) {
          coin = tokensAccumulated[keyIndex];
          amount = tokenAmounts[coin];
          if (Token(coin).transfer(beneficiary, amount)) {
            FundTransfer(beneficiary, tokensAccumulated[keyIndex], amountRaised, false);
            delete tokenAmounts[coin];
          } else {
            throw; // Hopefully throw will roll back anything we sent to the blockchain so far
          }            
        }
        keyIndex++;
      }
      return true;
    }
    return false;
  }
}
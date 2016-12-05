pragma solidity ^0.4.4;

import "Token.sol";
import "Recipient.sol";

contract Crowdsale is Recipient {

  address public beneficiary;

  uint public fundingGoal; 
  
  uint public amountRaised; 
  
  uint public deadline;

  address[] public tokensAccumulated;

  mapping(address => mapping(address => uint256)) public balanceOf;

  address[] accounts;

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
      address[] addressOfTokensAccumulated
  ) {
      beneficiary = ifSuccessfulSendTo;
      fundingGoal = fundingGoalInTokens;
      deadline = now + durationInMinutes * 1 minutes;
      tokensAccumulated = addressOfTokensAccumulated;
  }

  /* The function without name is the default function that is called whenever anyone sends funds to a contract */
  function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) {
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
      if(Token(_token).transferFrom(_from, _to, _value)) {
        uint amount = _value;
        balanceOf[_from][_token] = amount;
        accounts.push(_from);
        amountRaised += amount;
        checkGoalReached();
        FundTransfer(_from, _token, amount, true);
      }
  }

  modifier afterDeadline() { if (now >= deadline) _; }

  /* checks if the goal or time limit has been reached and ends the campaign */
  function checkGoalReached() afterDeadline {
      if (amountRaised >= fundingGoal){
          fundingGoalReached = true;
          GoalReached(beneficiary, amountRaised);
      }
      crowdsaleClosed = true;
  }

  function safeWithdrawal() afterDeadline {
      if (!fundingGoalReached) {
          uint i = 0;
          while (i < accounts.length) {
              address acc = accounts[i];
              mapping(address => uint256) accFunds = balanceOf[acc];
              uint j = 0;
              while (j < tokensAccumulated.length) {
                  address token = tokensAccumulated[j];
                  uint256 amount = accFunds[token];
                  if (amount > 0) {
                      Token(token).transfer(acc, amount);
                      FundTransfer(acc, token, amount, false);
                  }
              }
          }
      } else if (fundingGoalReached && beneficiary == msg.sender) {
          uint keyIndex = 0;
          bool notFound = true;
          while (keyIndex < tokensAccumulated.length) {
              if (Token(tokensAccumulated[keyIndex]).transfer(beneficiary, amountRaised)) {
                  FundTransfer(beneficiary, tokensAccumulated[keyIndex], amountRaised, false);
              } else {
                  //If we fail to send the funds to beneficiary, unlock funders balance
                  throw; //to rollback previous sending                  
              }            
              keyIndex++;
          }
     
      }
  }
}

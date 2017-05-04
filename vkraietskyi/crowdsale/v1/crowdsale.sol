contract Token {
  event Transfer(address from, address to, uint value);
  event Approval(address from, address to, uint value);
  function transferFrom(address from, address to, uint amount) returns (bool success);
  function transfer(address receiver, uint amount) returns (bool);
}

contract Crowdsale {

  address public beneficiary;
  address public creator;
  uint public fundingGoal;
  uint public amountRaised;
  uint public deadline;
  bool public closeOnGoalReached;
  address[] public tokensAccumulated;
  address[] public donators;
  bool fundingGoalReached = false;
  bool crowdsaleClosed = false;

  mapping(address => mapping(address => uint)) public balanceOf;
  mapping(address => uint) public tokenAmounts;
  event GoalReached(address beneficiary, uint amountRaised);
  event FundTransfer(address backer, address token, uint amount, bool isContribution);

  /* data structure to hold information about campaign contributors */
  /*  at initialization, setup the owner */
  function Crowdsale(
      address ifSuccessfulSendTo,
      uint fundingGoalInTokens,
      uint durationInMinutes,
      bool onGoalReached,
      address[] addressOfTokensAccumulated
  ) {
    creator = msg.sender;
    beneficiary = ifSuccessfulSendTo;
    fundingGoal = fundingGoalInTokens;
    deadline = now + durationInMinutes * 1 minutes;
    closeOnGoalReached = onGoalReached;
    setTokens(addressOfTokensAccumulated);
  }

  /* You must run this once only with the same tokens or else :) */
  function setTokens(address[] addressOfTokensAccumulated) internal returns (uint) {
    if (msg.sender != beneficiary && msg.sender != creator) return 0;
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

  function getTokensCount() constant returns (uint) {
    return tokensAccumulated.length;
  }

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
  function receiveApproval(address _from, uint _value, address _token, bytes _extraData) returns (bool)
 {
    if (!isTokenAccumulated(_token)) return false;
    address _to = this;
    if(!Token(_token).transferFrom(_from, _to, _value)) {
      return false;
    }
    if (balanceOf[_from][_token] == uint(0x0)) {
      balanceOf[_from][_token] = _value;
      donators.push(_from);
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
      if (fundingGoalReached) {
        crowdsaleClosed = true;
      }
    }
    return crowdsaleClosed;
  }

  function safeWithdrawal() returns (bool) {
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
      while(donatorsIndex < donators.length){
        keyIndex = 0;
        donator = donators[donatorsIndex];
        /* run through tokens*/
        while(keyIndex < tokensAccumulated.length){
          token = tokensAccumulated[keyIndex];
          /* if token address exists and it is pisible to get balance of donator by this token*/
          if (token != address(0x0) && balanceOf[donator][token] != uint(0x0)){
            amount = balanceOf[donator][token];
            if (amount > 0) {
              if (Token(token).transfer(donator, amount)) {
                FundTransfer(donator, token, amount, false);
              } else {
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
    /* if funding goal is reached then beneficiary can withdraw everything */
    if (fundingGoalReached) {
      bool notFound = true;
      keyIndex = 0;
      while (keyIndex < tokensAccumulated.length) {
        if (tokensAccumulated[keyIndex] != address(0x0)) {
          token = tokensAccumulated[keyIndex];
          amount = tokenAmounts[token];
          if (Token(token).transfer(beneficiary, amount)) {
            FundTransfer(beneficiary, tokensAccumulated[keyIndex], amountRaised, false);
            delete tokenAmounts[token];
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

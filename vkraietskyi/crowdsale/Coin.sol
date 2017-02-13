contract Recipient {
  event ReceivedApproval(address from, uint value, address tokenContract, bytes extraData);
  function receiveApproval(address from, uint value, address token, bytes extraData) returns (bool);
}
contract Coin {
  address minter;
  uint8 public tokenColor;
  mapping (address => uint) balances;
  mapping (address => mapping (address => uint)) allowed;
  event Transfer(address from, address to, uint value);
  event Approval(address from, address to, uint value);
  function Coin(uint8 _tokenColor) {
      minter = msg.sender;
      tokenColor = _tokenColor;
  }
  function setColor(uint8 _tokenColor) {
      if (msg.sender != minter) throw;
      tokenColor = _tokenColor;
  }
  function getColor() returns (uint8) {
      return tokenColor;
  }
  function mint(address owner, uint amount) {
      if (msg.sender != minter) return;
      balances[owner] += amount;
  }
  function transfer(address receiver, uint amount) returns (bool) {
      if (balances[msg.sender] < amount) return false;
      balances[msg.sender] -= amount;
      balances[receiver] += amount;
      Transfer(msg.sender, receiver, amount);
      return true;
  }
  function transferFrom(address from, address to, uint amount) returns (bool success) {
      if (balances[from] >= amount && allowed[from][msg.sender] >= amount && amount > 0) {
          balances[to] += amount;
          balances[from] -= amount;
          allowed[from][msg.sender] -= amount;
          Transfer(from, to, amount);
          return true;
      }
      return false;
  }
  function balanceOf(address owner) constant returns (uint balance) {
      return balances[owner];
  }
  function approve(address spender, uint amount) {
      allowed[msg.sender][spender] = amount;
      Approval(msg.sender, spender, amount);
      return;
  }
  function approveAndCall(address spender, uint amount) returns (bool success) {
      allowed[msg.sender][spender] = amount;
      Approval(msg.sender, spender, amount);
      return Recipient(spender).receiveApproval(msg.sender, amount, this, "\x00");
  }
  function allowance(address owner, address spender) constant returns (uint remaining) {
      return allowed[owner][spender];
  }
  function distribute(address[] accounts, uint amount) returns (bool success) {
      if (balances[msg.sender] < amount) return false;
      uint mean = amount / accounts.length;
      uint i = 0;
      while (i < accounts.length) {
          balances[msg.sender] -= mean;
          balances[accounts[i]] += mean;
          i = i + 1;
      }
      return true;
  }
}

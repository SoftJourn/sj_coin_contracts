/*
Implements ERC 20 Token standard: https://github.com/ethereum/EIPs/issues/20
.*/

import "Token.sol";
import "Recipient.sol";

contract CoinsToken is Token {

    address minter;
    
    mapping (address => uint256) balances;
    mapping (address => mapping (address => uint256)) allowed;

    function CoinsToken() {
        minter = msg.sender;
    }

    function mint(address owner, uint amount) {
        if (msg.sender != minter) throw;
        balances[owner] += amount;
    }

    function transfer(address _to, uint256 _value) returns (bool success) {
        if (balances[msg.sender] >= _value && _value > 0) {
            balances[msg.sender] -= _value;
            balances[_to] += _value;
            Transfer(msg.sender, _to, _value);
            return true;
        } else { return false; }
    }

    function transferFrom(address _from, address _to, uint256 _value) returns (bool success) {
        if (balances[_from] >= _value && allowed[_from][msg.sender] >= _value && _value > 0) {
            balances[_to] += _value;
            balances[_from] -= _value;
            allowed[_from][msg.sender] -= _value;
            Transfer(_from, _to, _value);
            return true;
        } else { return false; }
    }

    function balanceOf(address _owner) constant returns (uint256 balance) {
        return balances[_owner];
    }

    function approve(address _spender, uint256 _value) {
        allowed[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        return;
    }

    function approveAndCall(address _spender, uint256 _value) {
        allowed[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
    	Recipient(_spender).receiveApproval(msg.sender, _value, this, "\x00");
        return;
    }

    function allowance(address _owner, address _spender) constant returns (uint256 remaining) {
      return allowed[_owner][_spender];
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


/*
    Recepient implementation to withdraw coins like a cash 
*/

import "Recipient.sol";
import "Token.sol";

contract OfflineRecipient is Recipient {

    mapping (bytes32 => mapping (address => uint256)) public cheques;

    function receiveApproval(address from, uint256 value, address tokenContract, bytes extraData) {
         bytes32 data = withdraw(tokenContract, from, value);
    }

    function deposite(bytes32 chequeId, address tokenContract) returns (bool ret){
        uint amount = cheques[chequeId][tokenContract];
        if (amount == 0) {
            return false;
        } else {
	        if (tokenContract.call(bytes4(bytes32(sha3("transfer(address,uint256)"))), msg.sender, amount)) {
            	cheques[chequeId][tokenContract] = 0;
    	    	return true;
	        } else {
	    	    return false;
	        }
        }
    }

    function withdraw(address tokenContract, uint256 amount) returns (bytes32 cheque){
      	return withdraw(tokenContract, msg.sender, amount);
    } 

    function withdraw(address tokenContract, address from, uint256 amount) private returns (bytes32 cheque){
    	uint256 allowed = Token(tokenContract).allowance(from,this);
        if (allowed < amount) {
           throw;
        } else {
	        Token(tokenContract).transferFrom(from, this, amount);
            bytes32 chequeId = sha3(block.gaslimit, block.number, block.timestamp, from, amount);           
            cheques[chequeId][tokenContract] = amount;
            return chequeId;
        }
    }

}


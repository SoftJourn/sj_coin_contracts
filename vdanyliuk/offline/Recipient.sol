/*
    Standard recepient interface
*/

contract Recipient {

  event ReceivedApproval(address from, uint256 value, address tokenContract, bytes extraData);

  function receiveApproval(address from, uint256 value, address token, bytes extraData);

}

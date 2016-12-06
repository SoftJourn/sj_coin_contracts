import "VoteAuthenticator.sol";

contract BallotBox {
    VoteAuthenticator auth;
    address owner;
    
    function BallotBox(VoteAuthenticator _auth) {
        owner = msg.sender;
        setVoteAuthenticator(_auth);
    }
    
    function setVoteAuthenticator(VoteAuthenticator _auth) {
        if (msg.sender == owner)
            auth = _auth;
    }
    
    function submitVote(uint id) {
        if (auth.allow(msg.sender)) {
            // do things
        }
    }
}
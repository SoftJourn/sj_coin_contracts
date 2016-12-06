contract Coin {
    /* Public variables of the token */
    string public name = 'Cool Coin Token';
    string public symbol = 'CCT';
    uint8 public decimals = 2;
    uint public totalSupply = 0;
    address public authority;
    
    /* This creates an array with all balances */
    mapping (address => uint) public balanceOf;
    mapping (address => uint) public issuedBy;
    mapping (address => bool) public isIssuer;

    /* This generates a public event on the blockchain that will notify clients */
    event Transfer(address from, address to, uint value);
    event TokenIssued(address by);

    /* Constructor: initializes contract with initial supply tokens to the creator of the contract */
    function Coin() {
        authority = msg.sender;
        isIssuer[msg.sender] = true;
    }
    
    function addIssuer(address _newIssuer) public {
        if (msg.sender != authority) {
            throw;
        }
        isIssuer[_newIssuer] = true;
    }
    
    // TODO: removeIssuer()
    
    /* Issue CCT */
    function issue(uint _value) public {
        if (isIssuer[msg.sender] != true) {
            throw;
        }
        issuedBy[msg.sender] += _value;
        balanceOf[msg.sender] += _value;
        totalSupply += _value;
        TokenIssued(msg.sender);
    }

    /** Remove CCT from blockchain */
    function redeem(uint _value) public {
        if (isIssuer[msg.sender] != true) {
            throw;
        }
        if (balanceOf[msg.sender] < _value || issuedBy[msg.sender] < _value) {
            throw;                                  
        }

        issuedBy[msg.sender] -= _value;
        balanceOf[msg.sender] -= _value;
        totalSupply -= _value;
        TokenIssued(msg.sender);
    }

    /* Send coins */
    function transfer(address _to, uint _value) public {
        if (balanceOf[msg.sender] < _value) {
            throw;                                   // Check if the sender has enough
        }
        balanceOf[msg.sender] -= _value;             // Subtract from the sender
        balanceOf[_to] += _value;                    // Add the same to the recipient
        Transfer(msg.sender, _to, _value);           // Notify anyone listening that this transfer took place
    }

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () {
        throw;     // Prevents accidental sending of ether
    }
}
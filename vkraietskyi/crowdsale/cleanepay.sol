import * as User from "users.sol";
import * as Treasury from "treasury.sol";
import * as Categories from "categories.sol";

contract User {

    address creator;

    address treasury;
    address registry;
    address categories;z

    string firstname;
    string lastname;
    string email;
    string accountnumber;           // 'account number of user in format: CM341234123412xxxxxxxxxxx12, e.g. CM3412341234127678689693312',
    string status;                  // varchar(45) DEFAULT 'pending' COMMENT 'status fields: approved, rejected, pending',
    string cpid;                    // varchar(45) DEFAULT NULL COMMENT 'account number (cleanepay id), e.g.: 76786896933',
    bool admin;                     // tinyint(1) NOT NULL DEFAULT '0' COMMENT 'true (1) means admin rights, false (0) means no admin rights',
    string  currency;               // varchar(45) NOT NULL DEFAULT 'Euro' COMMENT 'Can be "Euro" or "CFA-Franc", nothing else, please exactly like this (capital and small letters)',
    uint256 gender;                 // int(11) DEFAULT NULL COMMENT '1 Male / 2 Female',
    string date_of_birth;           // date DEFAULT NULL,
    string place_of_birth;          // varchar(255) DEFAULT NULL,
    string passport_number;         // varchar(255) DEFAULT NULL,
    string _address;                // varchar(255) DEFAULT NULL,

    mapping(string => uint256) budget_by_category;
    string[] categoriesArray;

    mapping (bytes32 => uint256) cheques;

    function User(address _treasury, address _registry, address _categories, string _firstname, string _lastname, string _email, string _accountnumber, string _cpid,  bool _admin, string  _currency, uint256 _gender, string _date_of_birth, string _place_of_birth, string _passport_number, string _address_) {
        creator = msg.sender;

        treasury = _treasury;
        registry = _registry;
        categories = categories;

        firstname = _firstname;
        lastname = _lastname;
        email = _email;
        accountnumber = _accountnumber;
        status = "pending";
        cpid = _cpid;
        admin = _admin;
        currency = _currency;
        gender = _gender;
        date_of_birth = _date_of_birth;
        place_of_birth = _place_of_birth;
        passport_number = _passport_number;
        _address = _address_;

        Users(registry).add(_email, _accountnumber, _cpid, creator);
    }

    function approve() {
	    if (User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            status = "approved";
        } else {
            throw;
        }
    }

    function reject() {
	    if (User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            status = "rejected";
        } else {
            throw;
        }
    }

    function isAdmin() constant returns (bool isAdmin){
        return admin;
    }

    function getMoney() constant returns(uint256 money) {
	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            return Treasury(treasury).getAmount(this);
        } else {
            throw;
        }
    }

    function transfer(address to, uint256 count, uint256 category) {
	    if (creator == msg.sender && equal(status,"approved")) {
            Treasury(treasury).transfer(to, count, category);
        } else {
            throw;
        }
    }

    function approveMoneyRequest(bytes32 hash, uint256 category) {
        if (creator == msg.sender && equal(status,"approved")) {
            Treasury(treasury).approve(hash, category);
        } else {
            throw;
        }
    }

    function setBudget(string category, uint256 budget) returns (bool res) {
	    if ((creator == msg.sender && equal(status,"approved")) || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            if (Categories(categories).exists(category)) {
                budget_by_category[category] = budget;
                categoriesArray.push(category);
                return true;
            } else {
                return false;
            }
	    } else {
	        throw;
	    }
    }

    function getBudget() constant returns (uint256 _budget) {
	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            uint256 result = 0;
            uint i = 0;
            while (i < categoriesArray.length) {
            	result = result + budget_by_category[categoriesArray[i]];
            }
            return result;
        } else {
            throw;
        }
    }

    function getBudget(string category) constant returns (uint256 _budget) {
    	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
                return budget_by_category[category];
            } else {
                throw;
            }
        }

    function isSet(string category) constant private returns (bool res) {
        uint i = 0;
        while (i < categoriesArray.length) {
            if (equal(categoriesArray[i], category)) {
                return true;
            }
        }
        return false;
    }

    function account_1() constant returns (string _firstname, string _lastname, string _email, string _accountnumber, string _status, uint256 _amount) {
	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            uint256 amount = Treasury(treasury).getAmount(this);
            return (firstname, lastname, email, accountnumber, status, amount);
        } else {
            throw;
        }
    }

    function account_2() constant returns (string _cpid,  bool _admin, uint256 budget, string  _currency, uint256 _gender) {
	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            uint256 _budget = getBudget();
            return (cpid, admin, _budget, currency, gender);
        } else {
            throw;
        }
    }

    function account_3() constant returns (string _date_of_birth, string _place_of_birth, string _passport_number, string _address_) {
	    if (creator == msg.sender || User(Users(registry).findByCreator(msg.sender)).isAdmin()) {
            return (date_of_birth, place_of_birth, passport_number, _address);
        } else {
            throw;
        }
    }

    function compare(string _a, string _b) constant private returns (int) {
        bytes memory a = bytes(_a);
        bytes memory b = bytes(_b);
        uint minLength = a.length;
        if (b.length < minLength) minLength = b.length;
        for (uint i = 0; i < minLength; i ++)
            if (a[i] < b[i])
                return -1;
            else if (a[i] > b[i])
                return 1;
        if (a.length < b.length)
            return -1;
        else if (a.length > b.length)
            return 1;
        else
            return 0;
    }
    /// dev Compares two strings and returns true if they are equal.
    function equal(string _a, string _b) constant private returns (bool) {
        return compare(_a, _b) == 0;
    }
}
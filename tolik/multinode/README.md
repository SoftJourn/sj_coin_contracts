
# Hyperledger Burrow (Monax) multinode tutorial

This a basic multinode setup for 4 blocjchain nodes.

See https://wiki.softjourn.if.ua/technologyguilds:cryptofinance:hyperledger:burrow:multinode
based on https://monax.io/docs/chain-deploying/

---

## Vagrant

Setup virtual machines to run nodes:

  mkdir -p ~/Monax/Vagrant && cd ~/Monax/Vagrant
  mkdir CL{0,1,2,3}
  curl -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/Vagrantfile
  vagrant up

Check for password, something like:

  cat ~/.vagrant.d/boxes/ubuntu-VAGRANTSLASH-zesty64/20170412.1.0/virtualbox/Vagrantfile|grep password|cut -d '"' -f 2
  221c36362c947c7882bd3db1

==== Create & copy files ====

Create chain on CL0:
  monax chains make multichain --unsafe --account-types=Full:1,Validator:3 --seeds-ip=192.168.33.10:46656,192.168.33.11:46656,192.168.33.12:46656,192.168.33.13:46656

Copy to CL1 (on CL1):

  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_000 .

Copy to CL2 (on CL2):

  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_001 .

Copy to CL3 (on CL3):

  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_002 .

==== Start ====

On CL0, run:

  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_full_000 --logrotate

On CL1, run:

  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_000 --logrotate

On CL2, run:

  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_001 --logrotate

On CL3, run:

  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_002 --logrotate

==== Test contract ====

https://monax.io/docs/deploying-advanced-smart-contracts-to-a-chain/

On CL0 (get address & deploy):

  mkdir ~/.monax/apps/GSFactory && ~/.monax/apps/GSFactory
  culr -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/GSFactory/GSFactory.sol
  curl -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/GSFactory/epm.yaml
  cat ~/.monax/chains/multichain/accounts.json |grep "address"|head -n 1|cut -d '"' -f 4
  50297D60ADFDE7E04C4AA454A059366051EB86A8
  MONAX_PULL_APPROVE=true monax pkgs do --chain "multichain" --address "50297D60ADFDE7E04C4AA454A059366051EB86A8" --set "setStorageBase=5"

---

## The API & testing out authentication
We've created a quick little "API server" on [Google's Firebase Platform](https://firebase.google.com/). You can get your own API up and running within minutes too:

1. Signup for a [Firebase account](https://firebase.google.com/)
1. Create a new project - eg. "React Native Starter App"
1. Turn on email/password __Authentication__
1. Enable the __Database__ feature, and import the `firebase-sample-data.json` file found in this repo
1. Get the Firebase project's API credentials, copy `/.env.sample` to `/.env` and fill in the respective variables (eg. `APIKEY=d8f72k10s39djk29js`). You can get your projects details from Firebase, by clicking on the cog icon, next to overview > 'Add Firebase to your web app'.
1. Add the following __rules__ to the Database

```json
{
  "rules": {
    ".read": false,
    ".write": false,

    "meals": {
      ".read": true
    },

    "recipes": {
      ".read": true,
    	".indexOn": ["category"]
    },

    "users": {
      "$uid": {
        ".read": "auth != null && auth.uid == $uid",
        ".write": "auth != null && auth.uid == $uid",

        "firstName": { ".validate": "newData.isString() && newData.val().length > 0" },
        "lastName": { ".validate": "newData.isString() && newData.val().length > 0" },
        "lastLoggedIn": { ".validate": "newData.val() <= now" },
        "signedUp": { ".validate": "newData.val() <= now" },
        "role": {
          ".validate": "(root.child('users/'+auth.uid+'/role').val() === 'admin' && newData.val() === 'admin') || newData.val() === 'user'"
        }
      }
    },

    "favourites": {
    	"$uid": {
      	".read": "auth != null && auth.uid == $uid",
      	".write": "auth != null && auth.uid == $uid"
    	}
  	}
  }
}
```

Want to experiment even more with Firebase? Check out the [Firebase Cloud Functions](/docs/README.md)

---

## Understanding the File Structure

- `/android` - The native Android stuff
- `/ios` - The native iOS stuff
- `/src` - Contains the full React Native App codebase
  - `/components` - 'Dumb-components' / presentational. [Read More &rarr;](/src/components/README.md)
  - `/constants` - App-wide variables and config
  - `/containers` - 'Smart-components' / the business logic. [Read More &rarr;](/src/containers/README.md)
  - `/images` - Self explanatory right?
  - `/lib` - Utils, custom libraries, functions
  - `/navigation`- Routes - wire up the router with any & all screens. [Read More &rarr;](/src/navigation/README.md)
  - `/redux` - Redux Reducers & Actions grouped by type. [Read More &rarr;](/src/redux/README.md)
  - `/theme` - Theme specific styles and variables

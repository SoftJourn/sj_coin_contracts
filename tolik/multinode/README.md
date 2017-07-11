
# Hyperledger Burrow (Monax) multinode tutorial

This a basic multinode setup for 4 blocjchain nodes.

See [Softjourn's cryptofinance wiki](https://wiki.softjourn.if.ua/technologyguilds:cryptofinance:hyperledger:burrow:multinode)
based on https://monax.io/docs/chain-deploying/

---

## Vagrant

Vagrantfile will create 4 blockchain nodes & bootstrap monax:

1. `CL0` - 192.168.33.10
1. `CL1` - 192.168.33.11
1. `CL2` - 192.168.33.12
1. `CL3` - 192.168.33.13

- Setup virtual machines to run nodes:

```shell
  mkdir -p ~/Monax/Vagrant && cd ~/Monax/Vagrant
  mkdir CL{0,1,2,3}
  curl -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/Vagrantfile
  vagrant up
```
- Check for password, something like:

```shell
  cat ~/.vagrant.d/boxes/ubuntu-VAGRANTSLASH-zesty64/20170412.1.0/virtualbox/Vagrantfile|grep password|cut -d '"' -f 2
  221c36362c947c7882bd3db1
```

---

### Create & copy files

- Create chain on CL0:

```shell
  monax chains make multichain --unsafe --account-types=Full:1,Validator:3 --seeds-ip=192.168.33.10:46656,192.168.33.11:46656,192.168.33.12:46656,192.168.33.13:46656
```

- Copy to CL1 (on CL1):

```shell
  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_000 .
```

- Copy to CL2 (on CL2):

```shell
  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_001 .
```

- Copy to CL3 (on CL3):

```shell
  mkdir ~/.monax/chains/multichain && cd ~/.monax/chains/multichain
  scp -r ubuntu@192.168.33.10:~/.monax/chains/multichain/multichain_validator_002 .
```

---

### Start

- On CL0, run:

```shell
  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_full_000 --logrotate
```

- On CL1, run:

```shell
  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_000 --logrotate
```

- On CL2, run:

```shell
  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_001 --logrotate
```

- On CL3, run:

```shell
  MONAX_PULL_APPROVE=true monax chains start multichain --init-dir ~/.monax/chains/multichain/multichain_validator_002 --logrotate
```

---

### Test contract

- https://monax.io/docs/deploying-advanced-smart-contracts-to-a-chain/

- On CL0 (get address & deploy):

```shell
  mkdir ~/.monax/apps/GSFactory && ~/.monax/apps/GSFactory
  culr -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/GSFactory/GSFactory.sol
  curl -O https://raw.githubusercontent.com/SoftJourn/sj_coin_contracts/master/tolik/multinode/GSFactory/epm.yaml
  cat ~/.monax/chains/multichain/accounts.json |grep "address"|head -n 1|cut -d '"' -f 4
  50297D60ADFDE7E04C4AA454A059366051EB86A8
  MONAX_PULL_APPROVE=true monax pkgs do --chain "multichain" --address "50297D60ADFDE7E04C4AA454A059366051EB86A8" --set "setStorageBase=5"
```

### Node.js app

- `GSFactory/app.js` - sets/gets the number with GSContract
- `GSFactory/get.js` - gets the number with GSContract
- `GSFactory/get_addr.js` - gets GSContract address from GSFactory

---

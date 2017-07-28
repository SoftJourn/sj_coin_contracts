#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
        echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
        echo
        exit 1
fi
starttime=$(date +%s)
if [ -z "$1" ]; then
        echo "No token color argument supplied"
        echo
        exit 1
fi
if [ -z "$2" ]; then
        echo "No minter argument supplied"
        echo
        exit 1
fi
echo "POST Login Request on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=org1')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo

echo "POST Invoke request mint"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "peers": ["localhost:7051", "localhost:7056"],
        "fcn":"mint",
        "args":["1000000"]
}'
echo
echo

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/coin$1?peer=peer1&fcn=getColor&args=%5B%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/coin$1?peer=peer1&fcn=balanceOf&args=%5B%22$2%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

TESTLIST='"39d3f8d0655bd215e958890ea91b41330895a719cdefa9a690e0f59b2b4e9b14","4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5","553767d1cc2d34b5ac2ec1e45e1b8d51dc78ea091f5a16487ade1dd7e34c4097"'
echo "POST Invoke request for distribute"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "peers": ["localhost:7051", "localhost:7056"],
        "fcn":"distribute",
        "args":["1000",'$TESTLIST']
}'
echo
echo

echo "***** checking balances *****"
echo
TESTLIST=$TESTLIST',"'$2'"'
for i in `echo $TESTLIST|tr ',' ' '|sed 's/"//g'`; do
  echo "GET query chaincode on peer1 of Org1: balnceOf [$i]"
  echo
  curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/coin$1?peer=peer1&fcn=balanceOf&args=%5B%22$i%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
  echo
  echo
done

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

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"getColor",
	"args":[]
}'
echo
echo

echo "POST Invoke request"
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

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"balanceOf",
	"args":["'$2'"]
}'
echo
echo

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"transfer",
	"args":["1000","4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5"]
}'
echo
echo

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"balanceOf",
	"args":["4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5"]
}'
echo
echo

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"balanceOf",
	"args":["'$2'"]
}'
echo
echo

echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"approve",
	"args":["1000","4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5"]
}'
echo
echo


echo "POST Invoke request"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/coin$1 \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051", "localhost:7056"],
	"fcn":"allowance",
	"args":["4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5"]
}'
echo
echo


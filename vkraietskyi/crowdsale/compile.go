package main

import "fmt"
import (
	client "github.com/eris-ltd/eris-compilers/network"
)

func main() {
	url := "http://172.17.0.4:9099/compile"
	filename := "crowdsale.sol"
	optimize := true
	librariesString := "maLibrariez:0x1234567890"

	output, err := client.BeginCompile(url, filename, optimize, librariesString)

	contractName := output.Objects[0].Objectname // contract C would give you C here
	binary := output.Objects[0].Bytecode // gives you the binary
	abi := output.Objects[0].ABI // gives you the ABI
	fmt.Println("Error: ", err)
	fmt.Println("Contract name: ", contractName)
	fmt.Println("Contract bytecode: ", binary)
	fmt.Println("Contract interface: ", abi)
}
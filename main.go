package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"linumlabs/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL             = "https://goerli.drpc.org/"
	contractHexAddress = "0x3D216932E996c025E1d417c0396b1105a68963c6"
)

func main() {

	input := os.Args[1:]

	instance := connectToGoerli()

	switch input[0] {
	case "tokenOwner":
		if len(input) < 2 {
			inputHelper()
			log.Fatal("invalid number of arguments")
		}
		getTokenOwner(instance, input[1])
	case "mint":
		if len(input) < 3 {
			inputHelper()
			log.Fatal("invalid number of arguments")
		}
		mintOperation(instance, input[1], input[2])
	default:
		log.Println("invalid option passed")
		inputHelper()

	}

}

func connectToGoerli() *token.Token {
	client, err := ethclient.Dial(rpcURL)

	if err != nil {
		log.Fatal("Oops! There was a problem connecting to Goerli network ", err)
	}

	tokenAddress := common.HexToAddress(contractHexAddress)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	return instance
}

func getTokenOwner(instance *token.Token, tokenID string) {
	address, err := instance.OwnerOf(&bind.CallOpts{}, getBigIntFromHex(tokenID))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The owner of the token is", address.Hex())
}

func mintOperation(instance *token.Token, toURI, tokenURI string) {
	toAddress := common.HexToAddress(toURI)
	_, err := instance.Mint(&bind.TransactOpts{}, toAddress, tokenURI)

	if err != nil {
		log.Fatal("Oops! There was a problem ", err)
	}

	fmt.Println("Token successfully minted")
}

func inputHelper() {
	fmt.Println("Usage:")
	fmt.Println("   ./linumwallet tokenOwner <token ID>")
	fmt.Println("or")
	fmt.Println("   ./linumwallet mint <to address> <token ID>")
}

func getBigIntFromHex(original string) *big.Int {
	id := new(big.Int)
	id, ok := id.SetString(original, 16) // TODO check if its base 16 the ID
	if !ok {
		log.Fatal("invalid token ID")
	}
	return id
}

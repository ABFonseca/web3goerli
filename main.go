package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"linumlabs/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL             = "https://rpc.ankr.com/eth_goerli"
	contractHexAddress = "0x3D216932E996c025E1d417c0396b1105a68963c6"
)

// for purpose of testing the private key is here,
// in real code the private key would never be in clear text and never in the codebase, it would come from env variable or some secure service
var privateKey = "4718c7bbdde6c15114ff5e2ee1e9b2f25eaf8f1209c88a9632caec2aa5074819"

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
	address, err := instance.OwnerOf(&bind.CallOpts{}, getBigIntFromDecString(tokenID))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The owner of the token is", address.Hex())
}

func mintOperation(instance *token.Token, toURI, tokenURI string) {
	// Load the private key
	pvtKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(5))
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(toURI)

	tx, err := instance.Mint(transactOpts, toAddress, tokenURI)
	if err != nil {
		log.Fatal("error trying to mint ", err)
	}

	fmt.Println("Token successfully minted. Transaction HASH: ", tx.Hash())

}

func inputHelper() {
	fmt.Println("Usage:")
	fmt.Println("   ./linumwallet tokenOwner <token ID>")
	fmt.Println("or")
	fmt.Println("   ./linumwallet mint <to address> <token ID>")
}

func getBigIntFromDecString(original string) *big.Int {
	id := new(big.Int)
	id, ok := id.SetString(original, 10)
	if !ok {
		log.Fatal("invalid token ID")
	}
	return id
}

package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
)

func main() {

	clientNFT, error := nftables.New()
	defer clientNFT.CloseLasting()
	if error != nil {
		fmt.Println("Error Initializing nftables", error)
		os.Exit(1)
	} else {
		fmt.Println("nftables initialized")

	}

	chains, error := clientNFT. ListChains()
	if error != nil {
		fmt.Println("Error Getting Chains", error)
		os.Exit(1)

	for index, value := range chains {
		fmt.Printf("Chain: %v\n", value)
	}

}

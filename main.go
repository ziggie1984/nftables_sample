package main

import (
	"fmt"

	"github.com/google/nftables"
)

func main() {

	clientNFT, error := nftables.New()
	if error != nil {
		fmt.Println("Error Initializing nftables", error)
	} else {
		fmt.Println("nftables initialized")

	}

}

package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
)

func main() {

	option := nftables.AsLasting()

	clientNFT, error := nftables.New(option)

	// defer clientNFT.CloseLasting()

	if error != nil {
		fmt.Println("Error Initializing nftables", error)
		os.Exit(1)
	} else {
		fmt.Println("nftables initialized")
		fmt.Printf("Connection %v, %v\n", clientNFT.NetNS, clientNFT.TestDial)

	}

	wgTable := nftables.Table{
		Name:   "wg0",
		Family: nftables.TableFamilyINet,
	}

	prerouting := nftables.Chain{
		Name:     "FUCKYOUCHAIN",
		Table:    &wgTable,
		Hooknum:  nftables.ChainHookPrerouting,
		Priority: nftables.ChainPriorityNATDest,
		Type:     nftables.ChainTypeNAT,
	}
	fmt.Printf("Chain Self-Created: %v\n", prerouting)
	clientNFT.AddTable(&wgTable)
	clientNFT.AddChain(&prerouting)

	clientNFT.Flush()

	tables, error := clientNFT.ListTables()
	if error != nil {
		fmt.Println("Error Getting Chains", error)
		os.Exit(1)
	}

	for _, value := range tables {
		fmt.Printf("Table: %v\n", value)
	}

	chains, error := clientNFT.ListChains()
	if error != nil {
		fmt.Println("Error Getting Chains", error)
		os.Exit(1)
	}

	for _, value := range chains {
		fmt.Printf("Chain: %v\n", value)
	}

}

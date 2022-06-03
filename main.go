package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
)

func main() {

	clientNFT, error := nftables.New()

	defer func() {
		clientNFT.Flush()
		clientNFT.CloseLasting()
	}()

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
	table := clientNFT.AddTable(&wgTable)
	fmt.Printf("Table wg0: %v\n", table.Name)

	policy := nftables.ChainPolicyAccept
	prerouting := nftables.Chain{
		Name:     "FUCKYOUCHAIN",
		Table:    &wgTable,
		Hooknum:  nftables.ChainHookPrerouting,
		Priority: nftables.ChainPriorityNATDest,
		Type:     nftables.ChainTypeNAT,
		Policy:   &policy,
	}
	fmt.Printf("Chain Self-Created: %v\n", prerouting)

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

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

	chains, error := clientNFT.ListChains()
	if error != nil {
		fmt.Println("Error Getting Chains", error)
		os.Exit(1)
	}

	for _, value := range chains {
		fmt.Printf("Chain: %v\n", value)
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

	clientNFT.AddChain(&prerouting)

	clientNFT.Flush()

}

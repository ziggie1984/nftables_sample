package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func main() {

	option := nftables.AsLasting()

	clientNFT, error := nftables.New(option)

	defer clientNFT.CloseLasting()

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

	rule := nftables.Rule{
		Table: &wgTable,
		Chain: &prerouting,
		// The list of possible flags are specified by nftnl_rule_attr, see
		// https://git.netfilter.org/libnftnl/tree/include/libnftnl/rule.h#n21
		// Current nftables go implementation supports only
		// NFTNL_RULE_POSITION flag for setting rule at position 0
		Exprs: []expr.Any{},
	}
	clientNFT.AddRule(&rule)

	clientNFT.Flush()

}

package main

import (
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func main() {

	// // option := nftables.AsLasting()

	// // clientNFT, error := nftables.New(option)

	// clientNFT := &nftables.Conn{}

	// defer clientNFT.CloseLasting()

	// // if error != nil {
	// // 	fmt.Println("Error Initializing nftables", error)
	// // 	os.Exit(1)
	// // } else {
	// // 	fmt.Println("nftables initialized")
	// // 	fmt.Printf("Connection %v, %v\n", clientNFT.NetNS, clientNFT.TestDial)

	// // }

	// wgTable := &nftables.Table{
	// 	Name:   "testtable",
	// 	Family: nftables.TableFamilyINet,
	// }
	// clientNFT.AddTable(wgTable)

	// prerouting := &nftables.Chain{
	// 	Name:     "testchain",
	// 	Table:    wgTable,
	// 	Hooknum:  nftables.ChainHookPrerouting,
	// 	Priority: nftables.ChainPriorityNATDest,
	// 	Type:     nftables.ChainTypeNAT,
	// }
	// fmt.Printf("Chain Self-Created: %v\n", prerouting)
	// prerouting = clientNFT.AddChain(prerouting)

	// tables, error := clientNFT.ListTables()
	// if error != nil {
	// 	fmt.Println("Error Getting Chains", error)
	// 	os.Exit(1)
	// }

	// for _, value := range tables {
	// 	fmt.Printf("Table: %v\n", value)
	// }

	// chains, error := clientNFT.ListChains()
	// if error != nil {
	// 	fmt.Println("Error Getting Chains", error)
	// 	os.Exit(1)
	// }

	// for _, value := range chains {
	// 	fmt.Printf("Chain: %v\n", value)
	// }

	// // rule := nftables.Rule{
	// // 	Table: &wgTable,
	// // 	Chain: &prerouting,
	// // 	// The list of possible flags are specified by nftnl_rule_attr, see
	// // 	// https://git.netfilter.org/libnftnl/tree/include/libnftnl/rule.h#n21
	// // 	// Current nftables go implementation supports only
	// // 	// NFTNL_RULE_POSITION flag for setting rule at position 0
	// // 	Exprs: []expr.Any{},
	// // }
	// // clientNFT.AddRule(&rule)

	// set := &nftables.Set{
	// 	Name:    "whitelist",
	// 	Table:   wgTable,
	// 	KeyType: nftables.TypeIPAddr, // our keys are IPv4 addresses
	// }

	// clientNFT.AddRule(&nftables.Rule{
	// 	Table: wgTable,
	// 	Chain: prerouting,
	// 	Exprs: []expr.Any{
	// 		// [ payload load 4b @ network header + 16 => reg 1 ]
	// 		&expr.Payload{
	// 			DestRegister: 1,
	// 			Base:         expr.PayloadBaseNetworkHeader,
	// 			Offset:       16,
	// 			Len:          4,
	// 		},
	// 		// [ lookup reg 1 set whitelist ]
	// 		&expr.Lookup{
	// 			SourceRegister: 1,
	// 			SetName:        set.Name,
	// 			SetID:          set.ID,
	// 		},
	// 		//[ immediate reg 0 drop ]
	// 		&expr.Verdict{
	// 			Kind: expr.VerdictDrop,
	// 		},
	// 	},
	// })

	// clientNFT.Flush()

	// }
	c := &nftables.Conn{}

	// Basic boilerplate; create a table & chain.
	table := &nftables.Table{
		Family: nftables.TableFamilyIPv4,
		Name:   "ip_filter",
	}
	table = c.AddTable(table)

	myChain := c.AddChain(&nftables.Chain{
		Name:     "filter_chain",
		Table:    table,
		Type:     nftables.ChainTypeFilter,
		Hooknum:  nftables.ChainHookInput,
		Priority: nftables.ChainPriorityFilter,
	})

	set := &nftables.Set{
		Name:    "whitelist",
		Table:   table,
		KeyType: nftables.TypeIPAddr, // our keys are IPv4 addresses
	}

	// Create the set with a bunch of initial values.
	if err := c.AddSet(set, []nftables.SetElement{
		{Key: net.ParseIP("8.8.8.8")},
	}); err != nil {
		// handle error
	}

	c.AddRule(&nftables.Rule{
		Table: table,
		Chain: myChain,
		Exprs: []expr.Any{
			// [ payload load 4b @ network header + 16 => reg 1 ]
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       16,
				Len:          4,
			},
			// [ lookup reg 1 set whitelist ]
			&expr.Lookup{
				SourceRegister: 1,
				SetName:        set.Name,
				SetID:          set.ID,
			},
			//[ immediate reg 0 drop ]
			&expr.Verdict{
				Kind: expr.VerdictDrop,
			},
		},
	})
	if err := c.Flush(); err != nil {
		// handle error
	}
}

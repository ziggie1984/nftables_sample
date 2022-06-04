package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
)

func main() {

	option := nftables.AsLasting()

	clientNFT, error := nftables.New(option)

	// clientNFT := &nftables.Conn{}

	defer clientNFT.CloseLasting()

	if error != nil {
		fmt.Println("Error Initializing nftables", error)
		os.Exit(1)
	} else {
		fmt.Println("nftables initialized")
		fmt.Printf("Connection %v, %v\n", clientNFT.NetNS, clientNFT.TestDial)

	}

	wgTable := &nftables.Table{
		Name:   "wg0",
		Family: nftables.TableFamilyINet,
	}
	clientNFT.AddTable(wgTable)
	clientNFT.Flush()

	prerouting := clientNFT.AddChain(&nftables.Chain{
		Name:     "base-chain",
		Table:    wgTable,
		Type:     nftables.ChainTypeNAT,
		Hooknum:  nftables.ChainHookInput,
		Priority: nftables.ChainPriorityNATDest,
	})

	fmt.Printf("Chain Self-Created: %v\n", prerouting)
	prerouting = clientNFT.AddChain(prerouting)

	clientNFT.Flush()

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

	// clientNFT.Flush()
	// c := &nftables.Conn{}
	// defer c.CloseLasting()

	// filter := c.AddTable(&nftables.Table{
	// 	Family: nftables.TableFamilyIPv4,
	// 	Name:   "filter",
	// })

	// prerouting := c.AddChain(&nftables.Chain{
	// 	Name:     "base-chain",
	// 	Table:    filter,
	// 	Type:     nftables.ChainTypeFilter,
	// 	Hooknum:  nftables.ChainHookPrerouting,
	// 	Priority: nftables.ChainPriorityFilter,
	// })

	// c.AddRule(&nftables.Rule{
	// 	Table: filter,
	// 	Chain: prerouting,
	// 	Exprs: []expr.Any{
	// 		&expr.Verdict{
	// 			// [ immediate reg 0 drop ]
	// 			Kind: expr.VerdictAccept,
	// 		},
	// 	},
	// })

	// c.AddRule(&nftables.Rule{
	// 	Table: filter,
	// 	Chain: prerouting,
	// 	Exprs: []expr.Any{
	// 		&expr.Verdict{
	// 			// [ immediate reg 0 drop ]
	// 			Kind: expr.VerdictAccept,
	// 		},
	// 	},
	// })

	// c.InsertRule(&nftables.Rule{
	// 	Table: filter,
	// 	Chain: prerouting,
	// 	Exprs: []expr.Any{
	// 		&expr.Verdict{
	// 			// [ immediate reg 0 accept ]
	// 			Kind: expr.VerdictAccept,
	// 		},
	// 	},
	// })

	// c.InsertRule(&nftables.Rule{
	// 	Table: filter,
	// 	Chain: prerouting,
	// 	Exprs: []expr.Any{
	// 		&expr.Verdict{
	// 			// [ immediate reg 0 queue ]
	// 			Kind: expr.VerdictAccept,
	// 		},
	// 	},
	// })

	// if err := clientNFT.Flush(); err != nil {
	// 	fmt.Println(err)
	// }

	// rules, _ := c.GetRules(filter, prerouting)

	// for _, r := range rules {
	// 	rr, _ := r.Exprs[0].(*expr.Verdict)

	// 	fmt.Println(rr)
	// }

}

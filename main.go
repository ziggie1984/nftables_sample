package main

import (
	"fmt"
	"os"

	"github.com/google/nftables"
)

const networkIf = "tunnelsats"

func main() {
	settingUpFirewall()

}

func settingUpFirewall() {

	//Setup WG Table and Add Rule to Forward Chain

	option := nftables.AsLasting()

	nftClient, error := nftables.New(option)

	defer nftClient.CloseLasting()

	if error != nil {
		fmt.Println("Error Initializing nftables", error)
		os.Exit(1)
	} else {
		fmt.Println("nftables connection established")
		fmt.Printf("Connection %v\n", nftClient)

	}

	wgTable := &nftables.Table{
		Name:   "tunnelsats",
		Family: nftables.TableFamilyIPv4,
	}
	nftClient.AddTable(wgTable)
	fmt.Println("Creating Table: ", wgTable.Name, wgTable.Family)
	nftClient.Flush()

	//Adding Set
	datatypes := []nftables.SetDatatype{nftables.TypeIPAddr, nftables.TypeInetService}
	concat, err := nftables.ConcatSetType(datatypes...)
	if err != nil {
		fmt.Printf("Error %v\n", err)

	}

	portFw := &nftables.Set{
		Name:          "DNAT_LNPorts",
		Table:         wgTable,
		IsMap:         true,
		Concatenation: true,
		KeyType:       nftables.TypeInetService,
		DataType:      concat,
	}
	nftClient.AddSet(portFw, []nftables.SetElement{})
	nftClient.Flush()

	portFw_1 := &nftables.Set{
		Name:    "DNAT_LNPorts_Set",
		Table:   wgTable,
		KeyType: nftables.TypeInetService,
	}
	error = nftClient.AddSet(portFw_1, []nftables.SetElement{})
	if error != nil {
		// handle error
		fmt.Println(error)

	}
	nftClient.Flush()

	//Add A Sample Element

	// element := []byte("1.1.1.1 . 8080")

	error = nftClient.SetAddElements(portFw_1, []nftables.SetElement{{Key: []byte{1, 1}}})
	nftClient.Flush()

	if error != nil {
		fmt.Println(error)
	}
	error = nftClient.SetAddElements(portFw, []nftables.SetElement{{
		Key: []byte{1, 1},
		Val: []byte{1, 1, 1, 1, 1, 1},
	}})

	if error != nil {
		fmt.Println(error)
	}

	if err := nftClient.Flush(); err != nil {
		fmt.Println(err)
	}

	// prerouting := nftClient.AddChain(&nftables.Chain{
	// 	Name:     "base-chain",
	// 	Table:    wgTable,
	// 	Type:     nftables.ChainTypeRoute,
	// 	Hooknum:  nftables.ChainHookOutput,
	// 	Priority: nftables.ChainPriorityNATDest,
	// })

	// fmt.Printf("Chain Self-Created: %v\n", prerouting)
	// prerouting = nftClient.AddChain(prerouting)

	// nftClient.Flush()

	// wgTableInet := &nftables.Table{
	// 	Name:   "wg0_inet",
	// 	Family: nftables.TableFamilyINet,
	// }

	// nftClient.AddTable(wgTableInet)
	// nftClient.Flush()

	// preroutingInet := nftClient.AddChain(&nftables.Chain{
	// 	Name:     "base-chain",
	// 	Table:    wgTableInet,
	// 	Type:     nftables.ChainTypeRoute,
	// 	Hooknum:  nftables.ChainHookOutput,
	// 	Priority: nftables.ChainPriorityNATDest,
	// })

	// nftClient.AddChain(preroutingInet)

	// tables, error := nftClient.ListTables()
	// if error != nil {
	// 	fmt.Println("Error Getting Chains", error)
	// 	os.Exit(1)
	// }

	// for _, value := range tables {
	// 	fmt.Printf("Table: %v\n", value)
	// }

	// chains, error := nftClient.ListChains()
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
	// // nftClient.AddRule(&rule)

	// nftClient.Flush()
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

	// if err := nftClient.Flush(); err != nil {
	// 	fmt.Println(err)
	// }

	// rules, _ := c.GetRules(filter, prerouting)

	// for _, r := range rules {
	// 	rr, _ := r.Exprs[0].(*expr.Verdict)

	// 	fmt.Println(rr)
	// }

}

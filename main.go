package main

import (
	"fmt"

	"github.com/google/nftables"
)

func main() {

	option := nftables.AsLasting()

	nftClient, _ := nftables.New(option)

	defer nftClient.CloseLasting()

	wgTable := &nftables.Table{
		Name:   "wg0",
		Family: nftables.TableFamilyIPv4,
	}

	nftClient.AddTable(wgTable)
	nftClient.Flush()

	prerouting := nftClient.AddChain(&nftables.Chain{
		Name:     "base-chain",
		Table:    wgTable,
		Type:     nftables.ChainTypeRoute,
		Hooknum:  nftables.ChainHookOutput,
		Priority: nftables.ChainPriorityNATDest,
	})

	nftClient.AddChain(prerouting)

	nftClient.Flush()

	wgTableInet := &nftables.Table{
		Name:   "wg0_inet",
		Family: nftables.TableFamilyINet,
	}

	nftClient.AddTable(wgTableInet)
	nftClient.Flush()

	preroutingInet := nftClient.AddChain(&nftables.Chain{
		Name:     "base-chain",
		Table:    wgTableInet,
		Type:     nftables.ChainTypeRoute,
		Hooknum:  nftables.ChainHookOutput,
		Priority: nftables.ChainPriorityNATDest,
	})

	nftClient.AddChain(preroutingInet)

	if err := nftClient.Flush(); err != nil {
		fmt.Println(err)
	}

}

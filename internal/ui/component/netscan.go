package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

// external network object data implementation.
type network interface {
	GetSsid() string
	GetFreq() string
	GetLevel() string
}

func NetScan(synk *binding.Synk) ComponentInterface {
	return new("netScan", func() ComponentInterface {
		self := &netScan{
			Synk:  synk,
			table: style.ContentTable().SetSelectable(true, false),
			action: func(n network) {
				synk.TriggerModal(n.GetSsid())
			},
		}
		self.table.SetSelectedFunc(func(r, _ int) {
			self.action(self.networks[r-1])
		})
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type netScan struct {
	table    *tview.Table
	action   func(n network)
	networks []network
	*binding.Synk
}

func (n *netScan) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.table, 0, 1, true)
}

func (n *netScan) delay() int8 {
	return 5
}

func (n *netScan) draw() {
	n.table.Clear()
	n.table.SetCell(0, 0, style.CellHeader("ssid"))
	n.table.SetCell(0, 1, style.CellHeader("freq"))
	n.table.SetCell(0, 2, style.CellHeader("level"))

	for row, network := range n.networks {
		n.table.SetCell(row+1, 0, style.CellContent(network.GetSsid()))
		n.table.SetCell(row+1, 1, style.CellContent(network.GetFreq()))
		n.table.SetCell(row+1, 2, style.CellContent(network.GetLevel()))
	}
}

func (n *netScan) reload(ctx context.Context) {
	result, _ := net.NewNetworkManager().Scan(ctx)

	networks := []network{}
	for _, item := range result {
		networks = append(networks, item)
	}
	n.networks = networks
	n.draw()
	n.PutLog("netscan")
}

func (n *netScan) triggerAppDraw() {
	n.TriggerAppDraw()
}

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
			table: style.NewTable().WithTitle("networks").WithCount(0).AsContent(),
			action: func(n network) {
				synk.TriggerModal(n.GetSsid())
			},
		}
		self.table.Object.SetSelectedFunc(func(r, _ int) {
			self.action(self.networks[r-1])
		})
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type netScan struct {
	table    *style.Table
	action   func(n network)
	networks []network
	*binding.Synk
}

func (n *netScan) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.table.Object, 0, 1, true)
}

func (n *netScan) delay() int8 {
	return 5
}

func (n *netScan) draw() {
	n.table.Object.Clear()
	n.table.AddCellHeader(0, 0, "ssid")
	n.table.AddCellHeader(0, 1, "freq")
	n.table.AddCellHeader(0, 2, "level")

	for row, network := range n.networks {
		n.table.AddCellContent(row+1, 0, network.GetSsid())
		n.table.AddCellContent(row+1, 1, network.GetFreq())
		n.table.AddCellContent(row+1, 2, network.GetLevel())
	}
	n.table.WithCount(len(n.networks))
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

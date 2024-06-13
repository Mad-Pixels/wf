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

type netScan struct {
	*binding.Synk

	table    *style.Table
	action   func(n network)
	networks []network
}

func (n *netScan) delay() int8 {
	return 5
}

func (n *netScan) triggerAppDraw() {
	n.TriggerAppDraw()
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
	defer n.draw()
	n.networks = []network{}

	result, err := net.NewNetworkManager().Scan(ctx)
	if err != nil {
		n.PutLog(err.Error())
		return
	}
	for _, item := range result {
		n.networks = append(n.networks, item)
	}
}

func (n *netScan) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.table.Object, 0, 1, true)
}

func NetScan(synk *binding.Synk) ComponentInterface {
	return new("netscan", func() ComponentInterface {
		self := &netScan{
			Synk: synk,
			table: style.NewTable().
				WithTitle("networks").
				WithCount(0).
				AsContent(),
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

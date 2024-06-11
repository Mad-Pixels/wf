package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/rivo/tview"
)

// external network object data implementation.
type network interface {
	GetSsid() string
	GetFreq() string
	GetLevel() string
}

// external frame object with modal window.
type modal interface {
	ShowModal(text string, buttons []string, doneFunc func(buttonIndex int))
}

func NetScan(f modal, ch chan struct{}) ComponentInterface {
	return new("netScan", func() ComponentInterface {
		self := &netScan{
			table: tview.NewTable().SetSelectable(true, false),
			action: func(n network) {
				f.ShowModal(
					n.GetSsid(),
					[]string{
						"connect",
						"cancel",
					},
					func(buttonIndex int) {

					})
			},
			ch: ch,
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

	ch chan struct{}
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
	n.table.SetCell(0, 0, tview.NewTableCell("ssid").SetSelectable(false))
	n.table.SetCell(0, 1, tview.NewTableCell("freq").SetSelectable(false))
	n.table.SetCell(0, 2, tview.NewTableCell("level").SetSelectable(false))

	for row, network := range n.networks {
		n.table.SetCell(row+1, 0, tview.NewTableCell(network.GetSsid()))
		n.table.SetCell(row+1, 1, tview.NewTableCell(network.GetFreq()))
		n.table.SetCell(row+1, 2, tview.NewTableCell(network.GetLevel()))
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
	n.ch <- struct{}{}
}

package component

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/Mad-Pixels/wf/internal/ui/modal"
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
	LoggerInterface
	RenderInterface

	table    *style.Table
	networks []network
}

func (n *netScan) delay() int8 {
	return 5
}

func (n *netScan) renderComponent() {
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
	defer n.renderComponent()
	n.networks = []network{}

	result, err := net.NewNetwork().Scan(ctx)
	if err != nil {
		n.WriteMsg(err.Error())
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

func NetScan(render RenderInterface, logger LoggerInterface, view ViewInterface) ComponentInterface {
	return new("netscan", func() ComponentInterface {
		self := &netScan{
			LoggerInterface: logger,
			RenderInterface: render,

			table: style.NewTable().
				WithTitle("networks").
				WithCount(0).
				AsContent(),
		}
		self.table.Object.SetSelectedFunc(func(r, _ int) {
			view.Open(
				modal.NewWiFiConn(
					func() string {
						selectedRow, _ := self.table.Object.GetSelection()
						return self.networks[selectedRow-1].GetSsid()
					}(),
					func(ssid string) {
						err := net.NewNetwork().Conn(context.Background(), ssid, "Tislam55")
						if err != nil {
							self.WriteMsg(err.Error())
						}
						self.WriteMsg(fmt.Sprintf("exec proc for %s", ssid))
					},
				),
			)
		})
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

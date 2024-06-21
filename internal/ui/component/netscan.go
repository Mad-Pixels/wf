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
	GetBssid() string
	GetSsid() string
	GetMode() string
	GetChannel() string
	GetRate() string
	GetSignal() string
	GetBars() string
	GetSecurity() string

	// GetFreq() string
	// GetLevel() string
	// GetMac() string
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
	n.table.AddCellHeader(0, 0, "bssid")
	n.table.AddCellHeader(0, 1, "ssid")
	n.table.AddCellHeader(0, 2, "mode")
	n.table.AddCellHeader(0, 3, "channel")
	n.table.AddCellHeader(0, 4, "rate")
	n.table.AddCellHeader(0, 5, "signal")
	n.table.AddCellHeader(0, 6, "bars")
	n.table.AddCellHeader(0, 7, "security")

	for row, network := range n.networks {
		n.table.AddCellContent(row+1, 0, network.GetBssid())
		n.table.AddCellContent(row+1, 1, network.GetSsid())
		n.table.AddCellContent(row+1, 2, network.GetMode())
		n.table.AddCellContent(row+1, 3, network.GetChannel())
		n.table.AddCellContent(row+1, 4, network.GetRate())
		n.table.AddCellContent(row+1, 5, network.GetSignal())
		n.table.AddCellContent(row+1, 6, network.GetBars())
		n.table.AddCellContent(row+1, 7, network.GetSecurity())
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

func NetScan[R RenderInterface, L LoggerInterface, V ViewInterface](render R, logger L, view V) ComponentInterface {
	return new("netscan", func() ComponentInterface {
		self := &netScan{
			LoggerInterface: logger,
			RenderInterface: render,

			table: style.NewTable().
				WithTitle("networks").
				WithFixedHeader().
				WithExpansion().
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
						err := net.NewNetwork().Conn(context.Background(), ssid, "qwerty")
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

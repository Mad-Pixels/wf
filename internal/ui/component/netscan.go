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
	GetQuality() uint8
	GetFreq() uint32
	GetMaxBitrate() uint32
	GetMacAddr() string
	GetSecType() string
	// GetBssid() string
	// GetSsid() string
	// GetMode() string
	// GetChannel() string
	// GetRate() string
	// GetSignal() string
	// GetBars() string
	// GetSecurity() string

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
	n.table.AddCellHeader(0, 0, "ssid")
	n.table.AddCellHeader(0, 1, "quality")
	n.table.AddCellHeader(0, 2, "freq")
	n.table.AddCellHeader(0, 3, "maxbitrate")
	n.table.AddCellHeader(0, 4, "mac")
	n.table.AddCellHeader(0, 5, "sectype")

	for row, network := range n.networks {
		n.table.AddCellContent(row+1, 0, network.GetSsid())
		n.table.AddCellContent(row+1, 1, fmt.Sprintf("%d", network.GetQuality()))
		n.table.AddCellContent(row+1, 2, fmt.Sprintf("%d", network.GetFreq()))
		n.table.AddCellContent(row+1, 3, fmt.Sprintf("%d", network.GetMaxBitrate()))
		n.table.AddCellContent(row+1, 4, network.GetMacAddr())
		n.table.AddCellContent(row+1, 5, network.GetSecType())
	}
	n.table.WithCount(len(n.networks))
}

func (n *netScan) reload(ctx context.Context) {
	defer n.renderComponent()
	n.networks = []network{}

	// result, err := net.NewNetwork().Scan(ctx)
	// if err != nil {
	// 	n.WriteMsg(err.Error())
	// 	return
	// }
	result, err := net.Items()
	if err != nil {
		n.WriteMsg(err.Error())
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

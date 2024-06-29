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
	GetAccessType() string
	GetMaxBitrate() string
	GetMacAddr() string
	GetChannel() string
	GetQuality() string
	GetFreq() string
	GetSsid() string
	GetMode() string
}

type netScan struct {
	LoggerInterface
	RenderInterface

	table    *style.Table
	networks []network
}

func (n *netScan) delay() uint8 {
	return 7
}

func (n *netScan) renderComponent() {
	n.table.Object.Clear()
	n.table.AddCellHeader(0, 0, "ssid")
	n.table.AddCellHeader(0, 1, "quality")
	n.table.AddCellHeader(0, 2, "bitrate")
	n.table.AddCellHeader(0, 3, "sectype")
	n.table.AddCellHeader(0, 4, "channel")
	n.table.AddCellHeader(0, 5, "freq")
	n.table.AddCellHeader(0, 6, "mac")
	n.table.AddCellHeader(0, 7, "mode")

	for row, network := range n.networks {
		n.table.AddCellContent(row+1, 0, network.GetSsid())
		n.table.AddCellContent(row+1, 1, network.GetQuality())
		n.table.AddCellContent(row+1, 2, network.GetMaxBitrate())
		n.table.AddCellContent(row+1, 3, network.GetAccessType())
		n.table.AddCellContent(row+1, 4, network.GetChannel())
		n.table.AddCellContent(row+1, 5, network.GetFreq())
		n.table.AddCellContent(row+1, 6, network.GetMacAddr())
		n.table.AddCellContent(row+1, 7, network.GetMode())

	}
	n.table.WithCount(len(n.networks))
}

func (n *netScan) reload(ctx context.Context) {
	defer n.renderComponent()
	n.networks = []network{}

	result, err := net.Driver.WirelessAccessPoints()
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
					func(ssid, password string) {
						err := net.Driver.WirelessConnect(ssid, password)
						if err != nil {
							self.WriteMsg(err.Error())
							return
						}
						self.WriteMsg(fmt.Sprintf("connected to: %s", ssid))
					},
				),
			)
		})
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

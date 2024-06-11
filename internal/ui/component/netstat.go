package component

import (
	"context"

	"github.com/rivo/tview"
)

func NetStat() ComponentInterface {
	return new("netStat", func() ComponentInterface {
		self := &netStat{
			text: tview.NewTextView(),
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type netStat struct {
	text   *tview.TextView
	status string
}

func (n *netStat) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.text, 0, 1, false)
}

func (n *netStat) delay() int8 {
	return 3
}

func (n *netStat) draw() {
	n.text.SetText(n.status)
}

func (n *netStat) reload(ctx context.Context) {
	n.status = "test"
	n.draw()
}
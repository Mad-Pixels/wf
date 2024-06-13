package component

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

func NetStat(synk *binding.Synk) ComponentInterface {
	return new("netStat", func() ComponentInterface {
		self := &netStat{
			Synk: synk,
			text: style.NewText(),
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type netStat struct {
	text   *style.Text
	status string
	*binding.Synk
}

func (n *netStat) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.text.Object, 0, 1, false)
}

func (n *netStat) delay() int8 {
	return 3
}

func (n *netStat) draw() {
	n.text.Object.SetText(fmt.Sprintf("\n%s\n", n.status))
}

func (n *netStat) reload(ctx context.Context) {
	defer n.draw()

	var status = "n/a"
	info, err := net.NewNetworkManager().Stat(ctx)
	switch {
	case err != nil:
		n.PutLog(err.Error())
	case info == nil:
		status = "no active connection"
	default:
		status = fmt.Sprintf("connected: %s", info.GetSsid())
	}
	n.status = status
}

func (n *netStat) triggerAppDraw() {
	n.TriggerAppDraw()
}

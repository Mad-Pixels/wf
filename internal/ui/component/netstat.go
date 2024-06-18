package component

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/Mad-Pixels/wf/internal/ui/extension"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type netStat struct {
	RenderInterface
	logger *extension.Logger

	text   *style.Text
	status string
}

func (n *netStat) delay() int8 {
	return 3
}

func (n *netStat) renderComponent() {
	n.text.Object.SetText(fmt.Sprintf("\n%s\n", n.status))
}

func (n *netStat) reload(ctx context.Context) {
	defer n.renderComponent()

	var status = "n/a"
	info, err := net.NewNetwork().Stat(ctx)
	switch {
	case err != nil:
		n.logger.Put(err.Error())
	case info == nil:
		status = "no active connection"
	default:
		status = fmt.Sprintf("connected: %s", info.GetSsid())
	}
	n.status = status
}

func (n *netStat) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, n)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(n.text.Object, 0, 1, false)
}

func NetStat(render RenderInterface, logger *extension.Logger) ComponentInterface {
	return new("netstat", func() ComponentInterface {
		self := &netStat{
			RenderInterface: render,
			logger:          logger,
			text:            style.NewText(),
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

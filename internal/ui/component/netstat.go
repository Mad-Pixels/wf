package component

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Mad-Pixels/wf/internal/net"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type netStat struct {
	LoggerInterface
	RenderInterface

	text   *style.Text
	status string
}

func (n *netStat) delay() uint8 {
	return 5
}

func (n *netStat) renderComponent() {
	n.text.Object.SetText(fmt.Sprintf("\n%s\n", n.status))
}

func (n *netStat) reload(ctx context.Context) {
	defer n.renderComponent()
	var status = "n/a"

	ap, err := net.Driver.CurrentConnetcion()
	switch {
	case err != nil:
		n.WriteMsg(err.Error())
	case reflect.ValueOf(ap).IsNil():
		status = "no active connection"
	default:
		status = fmt.Sprintf("connected: %s", ap.GetSsid())
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

func NetStat[R RenderInterface, L LoggerInterface](render R, logger L) ComponentInterface {
	return new("netstat", func() ComponentInterface {
		self := &netStat{
			LoggerInterface: logger,
			RenderInterface: render,

			text: style.NewText(),
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

package component

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/wf/internal/ui/extension"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type helper struct {
	hotKeys *[]extension.Keys
	render  *extension.Render

	table *style.Table
}

func (h *helper) delay() int8 {
	return 100
}

func (n *helper) renderRoot() {
	n.render.Root()
}

func (h *helper) renderComponent() {
	if h.hotKeys == nil {
		return
	}
	for row, key := range *h.hotKeys {
		h.table.AddCellPrimary(row, 0, fmt.Sprintf("<%s>", tcell.KeyNames[key.Shortcut]))
		h.table.AddCellSecondary(row, 1, key.Description)
	}
}

func (h *helper) reload(_ context.Context) {}

func (h *helper) FlexItem(ctx context.Context) *tview.Flex {
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(h.table.Object, 0, 1, false)
}

func Helper(render *extension.Render, hotKeys *[]extension.Keys) ComponentInterface {
	return new("helper", func() ComponentInterface {
		self := &helper{
			render:  render,
			table:   style.NewTable(),
			hotKeys: hotKeys,
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

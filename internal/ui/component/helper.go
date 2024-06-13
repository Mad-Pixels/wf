package component

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type helper struct {
	*binding.Synk

	table   *style.Table
	hotKeys *[]binding.Keys
}

func (h *helper) delay() int8 {
	return 100
}

func (n *helper) triggerAppDraw() {
	n.TriggerAppDraw()
}

func (h *helper) draw() {
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

func Helper(hotKeys *[]binding.Keys, synk *binding.Synk) ComponentInterface {
	return new("helper", func() ComponentInterface {
		self := &helper{
			Synk:    synk,
			table:   style.NewTable(),
			hotKeys: hotKeys,
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

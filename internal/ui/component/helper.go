package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Helper(hotKeys *[]binding.Keys, synk *binding.Synk) ComponentInterface {
	return new("helper", func() ComponentInterface {
		self := &helper{
			Synk:    synk,
			table:   tview.NewTable(),
			hotKeys: hotKeys,
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type helper struct {
	table   *tview.Table
	hotKeys *[]binding.Keys
	*binding.Synk
}

func (h *helper) FlexItem(ctx context.Context) *tview.Flex {
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(h.table, 0, 1, false)
}

func (h *helper) delay() int8 {
	return 100
}

func (h *helper) draw() {
	if h.hotKeys == nil {
		return
	}
	for row, key := range *h.hotKeys {
		h.table.SetCell(row, 0, tview.NewTableCell(tcell.KeyNames[key.Shortcut]))
		h.table.SetCell(row, 1, tview.NewTableCell(key.Description))
	}
}

func (h *helper) reload(_ context.Context) {}

func (n *helper) triggerAppDraw() {
	n.TriggerAppDraw()

}

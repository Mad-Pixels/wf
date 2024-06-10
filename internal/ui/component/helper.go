package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Helper(hotKeys *[]ui.HotKeys) ComponentInterface {
	return new("helper", func() ComponentInterface {
		self := &helper{
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
	hotKeys *[]ui.HotKeys
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
		h.table.SetCell(row, 0, tview.NewTableCell(tcell.KeyNames[key.Key]))
		h.table.SetCell(row, 1, tview.NewTableCell(key.Description))
	}
}

func (h *helper) reload(_ context.Context) {}

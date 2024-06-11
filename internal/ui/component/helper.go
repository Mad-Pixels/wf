package component

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type keyBinding interface {
	Description() string
	Shortcut() tcell.Key
}

func Helper(hotKeys *[]keyBinding) ComponentInterface {
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
	hotKeys *[]keyBinding
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
		h.table.SetCell(row, 0, tview.NewTableCell(tcell.KeyNames[key.Shortcut()]))
		h.table.SetCell(row, 1, tview.NewTableCell(key.Description()))
	}
}

func (h *helper) reload(_ context.Context) {}

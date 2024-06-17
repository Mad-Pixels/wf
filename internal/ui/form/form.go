package form

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type Form struct {
	form *tview.Form
}

func (f Form) Content(title string) *style.Flex {
	flex := style.NewFlex()
	flex.Object.
		AddItem(nil, 0, 1, false).
		AddItem(
			style.NewFlex().WithRowDirection().Object.
				AddItem(nil, 0, 2, false).
				AddItem(f.primitive(title).Object, 0, 1, true).
				AddItem(nil, 0, 2, false),
			0, 2, true,
		).
		AddItem(nil, 0, 1, false)
	return flex
}

func (f Form) primitive(title string) *style.Flex {
	flex := style.
		NewFlex().
		AsModal(title)
	flex.Object.AddItem(f.form, 0, 1, true)
	return flex
}

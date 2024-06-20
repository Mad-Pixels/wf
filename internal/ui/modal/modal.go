package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
)

// Modal ...
type Modal struct {
	form      *style.Form
	CloseFunc func()
	title     string
	height    int
}

// Content return tview.Flex as a modal object.
func (m Modal) Content() *style.Flex {
	flex := style.NewFlex()
	flex.Object.
		AddItem(nil, 0, 1, false).
		AddItem(
			style.NewFlex().WithRowDirection().Object.
				AddItem(nil, 0, 1, false).
				AddItem(m.primitive(m.title).Object, m.height, 0, true).
				AddItem(nil, 0, 1, false),
			0, 2, true,
		).
		AddItem(nil, 0, 1, false)
	return flex
}

func (m Modal) primitive(title string) *style.Flex {
	flex := style.
		NewFlex().
		AsModal(title)
	flex.Object.
		AddItem(m.form.Object, 0, 1, true)
	return flex
}

package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
)

type Modal struct {
	form *style.Form
}

func (m Modal) Content(title string) *style.Flex {
	flex := style.NewFlex()
	flex.Object.
		AddItem(nil, 0, 1, false).
		AddItem(
			style.NewFlex().WithRowDirection().Object.
				AddItem(nil, 0, 2, false).
				AddItem(m.primitive(title).Object, 0, 1, true).
				AddItem(nil, 0, 2, false),
			0, 2, true,
		).
		AddItem(nil, 0, 1, false)
	return flex
}

func (m Modal) primitive(title string) *style.Flex {
	flex := style.
		NewFlex().
		AsModal(title)
	flex.Object.AddItem(m.form.Object, 0, 1, true)
	return flex
}

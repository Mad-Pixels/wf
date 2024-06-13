package style

import "github.com/rivo/tview"

type Text struct {
	Object *tview.TextView
}

func NewText() *Text {
	return &Text{
		Object: tview.NewTextView(),
	}
}

func (t *Text) AsLogger() *Text {
	return t
}

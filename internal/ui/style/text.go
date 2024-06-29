package style

import "github.com/rivo/tview"

// Text ...
type Text struct {
	Object *tview.TextView
}

// NewText return custom Text object with predefined styles.
func NewText() *Text {
	return &Text{
		Object: tview.NewTextView(),
	}
}

// AsLogger add custom styles for text.
func (t *Text) AsLogger() *Text {
	return t
}

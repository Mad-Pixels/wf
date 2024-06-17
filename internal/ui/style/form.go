package style

import "github.com/rivo/tview"

// Form ...
type Form struct {
	Object *tview.Form
}

// NewForm return custom Form object with predefined styles.
func NewForm() *Form {
	return &Form{
		Object: tview.NewForm(),
	}
}

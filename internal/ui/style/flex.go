package style

import "github.com/rivo/tview"

// Flex ...
type Flex struct {
	Object *tview.Flex
}

// NewFlex return custom Flex object with predefined styles.
func NewFlex() *Flex {
	return &Flex{
		Object: tview.NewFlex(),
	}
}

// WithRowDirection add direction inside flex object.
func (f *Flex) WithRowDirection() *Flex {
	f.Object.
		SetDirection(tview.FlexRow)
	return f
}

// WithColumnDirection add direction inside flex object.
func (f *Flex) WithColumnDirection() *Flex {
	f.Object.
		SetDirection(tview.FlexColumn)
	return f
}

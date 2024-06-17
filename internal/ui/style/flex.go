package style

import (
	"fmt"

	"github.com/rivo/tview"
)

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

func (f *Flex) AsModal(title string) *Flex {
	f.Object.
		SetDirection(tview.FlexRow).
		SetTitleColor(ColorTitle).
		SetTitle(fmt.Sprintf(" %s ", title))

	f.Object.
		SetBorder(true).
		SetBorderColor(ColorContent)
	return f
}

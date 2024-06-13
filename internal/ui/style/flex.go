package style

import "github.com/rivo/tview"

type Flex struct {
	Object *tview.Flex
}

func NewFlex() *Flex {
	return &Flex{
		Object: tview.NewFlex(),
	}
}

func (f *Flex) WithRowDirection() *Flex {
	f.Object.
		SetDirection(tview.FlexRow)
	return f
}

func (f *Flex) WithColumnDirection() *Flex {
	f.Object.
		SetDirection(tview.FlexColumn)
	return f
}

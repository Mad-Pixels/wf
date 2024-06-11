package frame

import (
	"github.com/rivo/tview"
)

var App = tview.NewApplication()

type Config struct {
	Header  []*tview.Flex
	Content *tview.Flex
}

func Builder(c Config) *tview.Frame {
	h := tview.NewFlex().SetDirection(tview.FlexColumn)
	for col, item := range c.Header {
		h.AddItem(item, 0, col, false)
	}

	frame := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(h, 5, 1, false).
		AddItem(c.Content, 0, 5, true)

	return tview.NewFrame(frame)
}

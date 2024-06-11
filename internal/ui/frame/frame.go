package frame

import (
	"github.com/rivo/tview"
)

type Frame struct {
	App  *tview.Application
	Root *tview.Flex
}

func NewFrame() *Frame {
	return &Frame{
		App: tview.NewApplication(),
	}
}

func (f *Frame) SetHeader(header []*tview.Flex) {
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	for col, item := range header {
		headerFlex.AddItem(item, 0, col, false)
	}
	f.Root.AddItem(headerFlex, 5, 1, false)
}

func (f *Frame) SetContent(content *tview.Flex) {
	f.Root.AddItem(content, 0, 5, true)
}

func (f *Frame) ShowModal(text string, buttons []string, doneFunc func(buttonIndex int)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, l string) {
			doneFunc(buttonIndex)
			f.App.SetRoot(f.Root, true).SetFocus(f.Root)
		})
	f.App.SetRoot(modal, true).SetFocus(modal)
}

func (f *Frame) Run() error {
	return f.App.SetRoot(f.Root, true).Run()
}

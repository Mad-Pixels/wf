package frame

import (
	"github.com/rivo/tview"
)

type Frame struct {
	App   *tview.Application
	Root  *tview.Flex
	Pages *tview.Pages
}

func NewFrame() *Frame {
	pages := tview.NewPages()
	return &Frame{
		App:   tview.NewApplication(),
		Pages: pages,
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
	f.Pages.AddPage("main", f.Root, true, true)
}

func (f *Frame) ShowModal(text string, buttons []string, doneFunc func(buttonIndex int)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, l string) {
			doneFunc(buttonIndex)
			f.Pages.HidePage("modal")
		})
	f.Pages.AddPage("modal", modal, true, true).ShowPage("modal")
	f.App.SetRoot(f.Pages, true).SetFocus(modal)
}

func (f *Frame) Run() error {
	f.App.SetRoot(f.Pages, true).SetFocus(f.Root)

	return f.App.SetRoot(f.Root, true).Run()
}

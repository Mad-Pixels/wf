package frame

import (
	"github.com/rivo/tview"
)

type Page struct {
	Flex *tview.Flex
}

func NewPage() *Page {
	return &Page{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
	}
}

func (p *Page) SetHeader(header []*tview.Flex) {
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	for col, item := range header {
		headerFlex.AddItem(item, 0, col, false)
	}
	p.Flex.AddItem(headerFlex, 5, 1, false)
}

func (p *Page) SetContent(content *tview.Flex) {
	p.Flex.AddItem(content, 0, 5, true)
}

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
		headerFlex.AddItem(item, 0, col+1, false)
	}
	p.Flex.AddItem(headerFlex, 8, 1, false)
}

func (p *Page) SetContent(content *tview.Flex) {
	p.Flex.AddItem(content, 0, 5, true)
}

func (p *Page) SetFooter(footer *tview.Flex) {
	p.Flex.AddItem(footer, 4, 1, false)
}

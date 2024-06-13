package frame

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
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
	headerFlex := style.NewFlex().WithColumnDirection()
	for col, item := range header {
		headerFlex.Object.AddItem(item, 0, col+1, false)
	}
	p.Flex.AddItem(headerFlex.Object, 8, 1, false)
}

func (p *Page) SetContent(content *tview.Flex) {
	p.Flex.AddItem(content, 0, 5, true)
}

func (p *Page) SetFooter(footer *tview.Flex) {
	footerFlex := style.NewFlex().WithRowDirection()
	footerFlex.Object.AddItem(footer, 0, 1, false)
	p.Flex.AddItem(footerFlex.Object, 6, 1, false)
}

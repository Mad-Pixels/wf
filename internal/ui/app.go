package ui

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/component"
	"github.com/Mad-Pixels/wf/internal/ui/frame"
	"github.com/rivo/tview"
)

type ui struct {
	app   *tview.Application
	pages *tview.Pages
}

func NewUI() *ui {
	return &ui{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}
}

func (u *ui) Run() error {
	return u.app.SetRoot(u.pages, true).Run()
}

// TODO: struct to bindig and manage it from chan
func (u *ui) ShowModal(text string, buttons []string, f func(btnID int)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(btnID int, label string) {
			f(btnID)

			u.pages.RemovePage("modal")
		})

	u.pages.AddPage("modal", modal, true, true).ShowPage("modal")
	u.app.SetFocus(modal)
}

func (u *ui) Draw() {
	u.app.Draw()
}

func Run() {
	ctx := context.Background()

	cch := make(chan struct{}, 1)
	mch := make(chan string, 1)

	ui := NewUI()
	page := frame.NewPage()

	sync := binding.NewSynk(cch, mch)

	sysInfo := component.SysInfo(sync).FlexItem(ctx)
	netStat := component.NetStat(*sync).FlexItem(ctx)

	netScan := component.NetScan(sync)

	go func() {
		for {
			select {
			case <-cch:
				ui.Draw()
			case val := <-mch:
				ui.ShowModal(val, []string{"cancel"}, func(btnID int) {})
			}
		}
	}()

	page.SetHeader([]*tview.Flex{sysInfo, netStat})
	page.SetContent(netScan.FlexItem(ctx))

	ui.pages.AddPage("main", page.Flex, true, true)

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

package ui

import (
	"context"
	"os"

	"github.com/Mad-Pixels/wf/internal/ui/component"
	"github.com/Mad-Pixels/wf/internal/ui/extension"
	"github.com/Mad-Pixels/wf/internal/ui/frame"
	"github.com/gdamore/tcell/v2"
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

func (u *ui) ShowModal(data extension.ModalData) {
	// data.P.Form.Object.AddButton("cancel", func() {
	// 	u.app.SetRoot(u.pages.ShowPage("main"), true)
	// })
	data.M.CloseFunc = func() {
		u.app.SetRoot(u.pages.ShowPage("main"), true)
	}

	container := tview.NewPages().
		AddPage("main", u.pages.ShowPage("main"), true, true).
		AddPage("modal", data.M.Content("connet").Object, true, true)

	u.app.SetRoot(container, true)
}

func (u *ui) Draw() {
	u.app.Draw()
}

func Run() {
	ctx := context.Background()

	mch := make(chan extension.ModalData, 1)
	ich := make(chan string, 10)

	ui := NewUI()
	page := frame.NewPage()

	renderCh := make(chan struct{}, 1)
	renderTrigger := extension.NewRender(renderCh)
	modalTr := extension.NewTriggerModal(mch)
	loggerTr := extension.NewLogger(ich)

	keys := []extension.Keys{
		{
			Description: "exit",
			Shortcut:    tcell.KeyCtrlC,
			Action: func(ctx context.Context) {
				ui.app.Stop()
				os.Exit(0)
			},
		},
	}

	sysInfo := component.SysInfo(renderTrigger, loggerTr).FlexItem(ctx)
	netStat := component.NetStat(renderTrigger, loggerTr).FlexItem(ctx)
	helper := component.Helper(renderTrigger, loggerTr, &keys).FlexItem(ctx)
	std := component.StdOut(renderTrigger, loggerTr).FlexItem(ctx)

	netScan := component.NetScan(renderTrigger, loggerTr, modalTr)
	info := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sysInfo, 5, 0, false).
		AddItem(netStat, 5, 1, false)

	go func() {
		for {
			select {
			case <-renderCh:
				ui.Draw()
			case val := <-mch:
				ui.ShowModal(val)
			}
		}
	}()

	page.SetHeader([]*tview.Flex{info, helper})
	page.SetContent(netScan.FlexItem(ctx))
	page.SetFooter(std)

	ui.pages.AddPage("main", page.Flex, true, true)

	ui.app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			for _, k := range keys {
				if k.Shortcut == event.Key() {
					k.Action(ctx)
					break
				}
			}
			return event
		},
	)

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

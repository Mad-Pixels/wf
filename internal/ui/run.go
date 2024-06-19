package ui

import (
	"context"
	"os"

	"github.com/Mad-Pixels/wf/internal/ui/component"
	"github.com/Mad-Pixels/wf/internal/ui/extension"
	"github.com/Mad-Pixels/wf/internal/ui/frame"
	"github.com/Mad-Pixels/wf/internal/ui/modal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Run is a main endpoint for ui.
func Run(ctx context.Context) error {
	var (
		gCtx, cancel = context.WithCancel(ctx)
		ui           = NewUI()

		hotKeys = []extension.Keys{
			{
				Description: "close and exit cli",
				Shortcut:    tcell.KeyCtrlC,

				Action: func(ctx context.Context) {
					defer os.Exit(0)

					cancel()
					ui.app.Stop()
				},
			},
			{
				Description: "go to main view",
				Shortcut:    tcell.KeyESC,

				Action: func(ctx context.Context) {
					ui.app.SetRoot(ui.pages.ShowPage("main"), true)
				},
			},
		}

		chView   = make(chan *modal.Modal, 1)
		chRender = make(chan struct{}, 1)
		chlogger = make(chan string, 10)

		triggerRender = extension.NewRender(chRender)
		triggerLogger = extension.NewLogger(chlogger)
		triggerView   = extension.NewView(chView)

		componentNetScan = component.NetScan(triggerRender, triggerLogger, triggerView).FlexItem(gCtx)
		componentHelper  = component.Helper(triggerRender, triggerLogger, &hotKeys).FlexItem(gCtx)
		componentSysInfo = component.SysInfo(triggerRender, triggerLogger).FlexItem(gCtx)
		componentNetStat = component.NetStat(triggerRender, triggerLogger).FlexItem(gCtx)
		componentStdOut  = component.StdOut(triggerRender, triggerLogger).FlexItem(gCtx)
	)
	go func() {
		for {
			select {
			case <-chRender:
				ui.Draw()
			case view := <-chView:
				ui.ShowModal(view)
			}
		}
	}()

	page := frame.NewPage()
	page.SetHeader(
		[]*tview.Flex{
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(componentSysInfo, 5, 0, false).
				AddItem(componentNetStat, 5, 1, false),
			componentHelper,
		},
	)
	page.SetContent(componentNetScan)
	page.SetFooter(componentStdOut)

	ui.pages.AddPage("main", page.Flex, true, true)
	ui.app.SetInputCapture(
		func(e *tcell.EventKey) *tcell.EventKey {
			for _, k := range hotKeys {
				if k.Shortcut == e.Key() {
					k.Action(gCtx)
					break
				}
			}
			return e
		},
	)
	return ui.Run()
}

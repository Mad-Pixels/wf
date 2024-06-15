package ui

import (
	"context"
	"os"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/component"
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

func (u *ui) ShowModal(text string, buttons []string, f func(btnID int, ssid, password string)) {
	form := tview.NewForm().
		AddInputField("SSID", "", 20, nil, nil).
		AddInputField("Password", "", 20, nil, nil)
	form.AddButton("Connect", func() {
		ssid := form.GetFormItemByLabel("SSID").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		f(1, ssid, password)
		u.pages.RemovePage("modal")
	}).
		AddButton("Cancel", func() {
			f(0, "", "")
			u.pages.RemovePage("modal")
		})

	form.SetButtonsAlign(tview.AlignCenter)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(false), 1, 0, false). // Add some padding at the top
		AddItem(tview.NewTextView().SetText(text).SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(form, 0, 1, true)

	u.pages.AddPage("modal", layout, true, true).ShowPage("modal")
	u.app.SetFocus(form)
}

// TODO: struct to bindig and manage it from chan
// func (u *ui) ShowModal(text string, buttons []string, f func(btnID int, ssid, password string)) {
// 	form := tview.NewForm().
// 		AddInputField("SSID", "", 20, nil, nil).
// 		AddInputField("Password", "", 20, nil, nil)

// 	modal := tview.NewModal().
// 		SetText(text).
// 		AddButtons(buttons).
// 		SetDoneFunc(func(btnID int, label string) {
// 			ssid := form.GetFormItemByLabel("SSID").(*tview.InputField).GetText()
// 			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
// 			f(btnID, ssid, password)

// 			u.pages.RemovePage("modal")
// 		})

// 	layout := tview.NewFlex().SetDirection(tview.FlexRow).
// 		AddItem(form, 0, 1, true).
// 		AddItem(modal, 0, 1, false)

// 	u.pages.AddPage("modal", layout, true, true).ShowPage("modal")
// 	u.app.SetFocus(modal)
// }

func (u *ui) Draw() {
	u.app.Draw()
}

func Run() {
	ctx := context.Background()

	cch := make(chan struct{}, 1)
	mch := make(chan func(id int, ssid, password string), 1)
	ich := make(chan string, 10)

	ui := NewUI()
	page := frame.NewPage()

	sync := binding.NewSynk(cch, mch, ich)

	keys := []binding.Keys{
		{
			Description: "exit",
			Shortcut:    tcell.KeyCtrlC,
			Action: func(ctx context.Context) {
				ui.app.Stop()
				os.Exit(0)
			},
		},
	}

	sysInfo := component.SysInfo(sync).FlexItem(ctx) //component.SysInfo(sync).FlexItem(ctx)
	netStat := component.NetStat(sync).FlexItem(ctx)
	helper := component.Helper(&keys, sync).FlexItem(ctx)
	std := component.StdOut(sync).FlexItem(ctx)

	netScan := component.NetScan(sync)
	info := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sysInfo, 5, 0, false).
		AddItem(netStat, 5, 1, false)

	go func() {
		for {
			select {
			case <-cch:
				ui.Draw()
			case val := <-mch:
				ui.ShowModal("mock", []string{"connect", "cancel"}, val)
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

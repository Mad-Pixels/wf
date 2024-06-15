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

func (u *ui) ShowModal(data binding.TriggerModalData) {
	//form := tview.NewForm().
	//	AddInputField("SSID", "", 20, nil, nil).
	//	AddInputField("Password", "", 20, nil, nil)
	//form.AddButton("Connect", func() {
	//	ssid := form.GetFormItemByLabel("SSID").(*tview.InputField).GetText()
	//	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	//	f(1, ssid, password)
	//	u.pages.RemovePage("modal")
	//}).
	//	AddButton("Cancel", func() {
	//		f(0, "", "")
	//		u.pages.RemovePage("modal")
	//	})
	//
	//form.SetButtonsAlign(tview.AlignCenter)

	//layout := tview.NewFlex().SetDirection(tview.FlexRow).
	//	AddItem(tview.NewBox().SetBorder(false), 1, 0, false). // Add some padding at the top
	//	AddItem(tview.NewTextView().SetText(text).SetTextAlign(tview.AlignCenter), 1, 0, false).
	//	AddItem(form, 0, 1, true)
	//u.pages.AddPage("modal", data.P, true, true)

	//modal := tview.NewModal().
	//	SetText("Enter your SSID and Password").
	//	AddButtons([]string{"Close"}).
	//	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	//		u.pages.HidePage("modal")
	//	})

	//flex := tview.NewFlex().
	//	SetDirection(tview.FlexRow).
	//	AddItem(nil, 1, 1, false).
	//	AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
	//		AddItem(nil, 1, 1, false).
	//		AddItem(data.P, 0, 1, true).
	//		AddItem(nil, 1, 1, false), 0, 1, true).
	//	AddItem(nil, 1, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Пустой элемент для выравнивания по вертикали
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false). // Пустой элемент для выравнивания по горизонтали
			AddItem(data.P, 60, 1, true).
			AddItem(nil, 0, 1, false), 10, 1, true). // Пустой элемент для выравнивания по горизонтали
		AddItem(nil, 0, 1, false)

	u.pages.AddPage("modal", flex, true, false)
	//container := tview.NewPages().
	//	AddPage("background", u.pages, true, true).
	//	AddPage("modal", page, true, true)
	u.pages.ShowPage("modal")
	u.app.SetFocus(data.P)

	//u.pages.AddPage("modal", container, false, true).ShowPage("modal")
	//u.app.SetFocus(data.P)
	//u.app.SetFocus(form)
}

func (u *ui) Draw() {
	u.app.Draw()
}

func Run() {
	ctx := context.Background()

	cch := make(chan struct{}, 1)
	mch := make(chan binding.TriggerModalData, 1)
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

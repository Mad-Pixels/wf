package ui

import (
	"github.com/Mad-Pixels/wf/internal/ui/modal"
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

func (u *ui) ShowModal(data *modal.Modal) {
	data.CloseFunc = func() {
		u.app.SetRoot(u.pages.ShowPage("main"), true)
	}

	container := tview.NewPages().
		AddPage("main", u.pages.ShowPage("main"), true, true).
		AddPage("modal", data.Content("connect").Object, true, true)

	u.app.SetRoot(container, true)
}

func (u *ui) Draw() {
	u.app.Draw()
}

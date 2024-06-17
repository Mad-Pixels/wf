package form

import "github.com/rivo/tview"

func NewWiFiConn() *Form {
	f := tview.NewForm().
		AddInputField("SSID:", "", 20, nil, nil).
		AddInputField("Password:", "", 20, nil, nil)
	f.AddButton("connect", func() {})
	f.AddButton("cancel", func() {})

	return &Form{
		form: f,
	}
}

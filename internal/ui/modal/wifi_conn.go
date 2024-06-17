package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

// NewWifiConn is a modal object for WiFi connection.
func NewWiFiConn(ssid string, conn func(ssid string)) *Modal {
	modal := &Modal{
		CloseFunc: func() {},
	}

	form := style.NewForm()
	form.Object.
		AddInputField("SSID:", ssid, 0, nil, nil).
		AddInputField("Password:", "", 0, nil, nil)
	form.Object.
		AddButton("connect", func() {
			defer modal.CloseFunc()
			conn(form.Object.GetFormItemByLabel("SSID:").(*tview.InputField).GetText())
		}).
		AddButton("cancel", func() {
			modal.CloseFunc()
		})

	modal.form = form
	return modal
}

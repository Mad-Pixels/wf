package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

// NewWifiConn is a modal object for WiFi connection.
func NewWiFiConn(ssid string, conn func(ssid string)) *Modal {
	var (
		ssidField = "SSID:"
		passField = "password:"

		connectBtn = "connect"
		cancelBtn  = "cancel"
	)
	modal := &Modal{
		title:  "connect to WiFi",
		height: 9,

		CloseFunc: func() {},
	}

	form := style.NewForm()
	form.Object.
		AddInputField(ssidField, ssid, 0, nil, nil).
		AddInputField(passField, "", 0, nil, nil)
	form.Object.
		AddButton(connectBtn, func() {
			defer modal.CloseFunc()
			conn(form.Object.GetFormItemByLabel(ssidField).(*tview.InputField).GetText())
		}).
		AddButton(cancelBtn, func() { modal.CloseFunc() })
	modal.form = form
	return modal
}

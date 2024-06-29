package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

// NewWifiConn is a modal object for WiFi connection.
func NewWiFiConn(ssid string, conn func(ssid, password string)) *Modal {
	var (
		ssidField = "SSID:"
		passField = "Password:"

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
		AddPasswordField(passField, "", 0, '*', nil)
	form.Object.
		AddButton(connectBtn, func() {
			defer modal.CloseFunc()
			conn(
				form.Object.GetFormItemByLabel(ssidField).(*tview.InputField).GetText(),
				form.Object.GetFormItemByLabel(passField).(*tview.InputField).GetText(),
			)
		}).
		AddButton(cancelBtn, func() { modal.CloseFunc() })
	form.Object.
		SetButtonBackgroundColor(style.ColorDark).
		SetFieldBackgroundColor(style.ColorDark).
		SetFieldTextColor(style.ColorText).
		SetLabelColor(style.ColorText)
	modal.form = form
	return modal
}

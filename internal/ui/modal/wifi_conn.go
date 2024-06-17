package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
)

// NewWifiConn is a modal object for WiFi connection.
func NewWiFiConn(ssid string) *Modal {
	form := style.NewForm()
	form.Object.
		AddInputField("SSID:", ssid, 0, nil, nil).
		AddInputField("Password:", "", 0, nil, nil).
		AddButton("connect", func() {}).
		AddButton("cancel", func() {})
	return &Modal{
		form: form,
	}
}

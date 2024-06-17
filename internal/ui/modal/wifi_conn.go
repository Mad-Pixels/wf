package modal

import (
	"github.com/Mad-Pixels/wf/internal/ui/style"
)

func NewWiFiConn() *Modal {
	form := style.NewForm()
	form.Object.
		AddInputField("SSID:", "", 0, nil, nil).
		AddInputField("Password:", "", 0, nil, nil).
		AddButton("connect", func() {}).
		AddButton("cancel", func() {})
	return &Modal{
		form: form,
	}
}

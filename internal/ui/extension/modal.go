package extension

import "github.com/Mad-Pixels/wf/internal/ui/modal"

type ModalData struct {
	M *modal.Modal
}

// OpenModal ...
type TriggerModal struct {
	triggerCh chan ModalData
}

func NewTriggerModal(ch chan ModalData) *TriggerModal {
	return &TriggerModal{
		triggerCh: ch,
	}
}

// Root trigger for open modal window.
func (e TriggerModal) Root(data ModalData) {
	e.triggerCh <- data
}

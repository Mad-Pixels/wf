package binding

import (
	"github.com/Mad-Pixels/wf/internal/ui/modal"
)

type Synk struct {
	drawCh chan struct{}
	modal  chan TriggerModalData
	std    chan string
}

func NewSynk(ch chan struct{}, chm chan TriggerModalData, chio chan string) *Synk {
	return &Synk{
		drawCh: ch,
		modal:  chm,
		std:    chio,
	}
}

func (s *Synk) TriggerAppDraw() {
	s.drawCh <- struct{}{}
}

type TriggerModalData struct {
	Title string
	P     *modal.Modal
}

func (s *Synk) TriggerModal(tr TriggerModalData) {
	s.modal <- tr
}

func (s *Synk) PutLog(data string) {
	s.std <- data
}

func (s *Synk) ReadLog() chan string {
	return s.std
}

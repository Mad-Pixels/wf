package binding

type Synk struct {
	drawCh chan struct{}
	modal  chan string
}

func NewSynk(ch chan struct{}, chm chan string) *Synk {
	return &Synk{
		drawCh: ch,
		modal:  chm,
	}
}

func (s *Synk) TriggerAppDraw() {
	s.drawCh <- struct{}{}
}

func (s *Synk) TriggerModal(ssid string) {
	s.modal <- ssid
}

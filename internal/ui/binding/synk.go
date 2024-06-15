package binding

type Synk struct {
	drawCh chan struct{}
	modal  chan func(id int, ssid, password string)
	std    chan string
}

func NewSynk(ch chan struct{}, chm chan func(id int, ssid, password string), chio chan string) *Synk {
	return &Synk{
		drawCh: ch,
		modal:  chm,
		std:    chio,
	}
}

func (s *Synk) TriggerAppDraw() {
	s.drawCh <- struct{}{}
}

func (s *Synk) TriggerModal(f func(id int, ssid, password string)) {
	s.modal <- f
}

func (s *Synk) PutLog(data string) {
	s.std <- data
}

func (s *Synk) ReadLog() chan string {
	return s.std
}

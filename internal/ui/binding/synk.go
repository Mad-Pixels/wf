package binding

type Synk struct {
	drawCh chan struct{}
}

func NewSynk(ch chan struct{}) *Synk {
	return &Synk{
		drawCh: ch,
	}
}

func (s *Synk) TriggerAppDraw() {
	s.drawCh <- struct{}{}
}

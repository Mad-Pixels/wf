package extension

// TriggerDraw ...
type TriggerDraw struct {
	triggerCh chan struct{}
}

func NewTriggerDraw(ch chan struct{}) *TriggerDraw {
	return &TriggerDraw{
		triggerCh: ch,
	}
}

// Root trigger for refresh application frames.
func (e TriggerDraw) Root() {
	e.triggerCh <- struct{}{}
}

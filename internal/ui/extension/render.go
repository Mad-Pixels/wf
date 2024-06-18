package extension

// Render ...
type Render struct {
	triggerCh chan struct{}
}

func NewRender(ch chan struct{}) *Render {
	return &Render{
		triggerCh: ch,
	}
}

// Root trigger for refresh application frames.
func (e Render) DrawRootFrame() {
	e.triggerCh <- struct{}{}
}

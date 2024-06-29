package extension

import "github.com/Mad-Pixels/wf/internal/ui/modal"

// View ...
type View struct {
	ch chan *modal.Modal
}

// NewView return View object.
func NewView(ch chan *modal.Modal) *View {
	return &View{
		ch: ch,
	}
}

// Open is a trigger for open modal window.
func (e View) Open(modal *modal.Modal) {
	e.ch <- modal
}

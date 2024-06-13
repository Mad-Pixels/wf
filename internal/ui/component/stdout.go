package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

func StdOut(synk *binding.Synk) ComponentInterface {
	return new("stdout", func() ComponentInterface {
		self := &stdout{
			Synk: synk,
			text: style.NewText().AsLogger(),
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

type stdout struct {
	text *style.Text
	*binding.Synk
	draft string
}

func (s *stdout) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, s)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.text.Object, 0, 1, false)
}

func (s *stdout) delay() int8 {
	return 1
}

func (s *stdout) draw() {
	s.text.Object.SetScrollable(true).
		ScrollToEnd().
		SetDynamicColors(true)
}

func (s *stdout) reload(ctx context.Context) {
	for {
		select {
		case data := <-s.ReadLog():
			s.draft += (data + "\n")
			s.text.Object.SetText(s.draft)
		default:
			s.draw()
			return
		}
	}
}

func (s *stdout) triggerAppDraw() {
	s.TriggerAppDraw()
}
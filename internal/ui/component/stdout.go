package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type stdout struct {
	LoggerInterface
	RenderInterface

	text  *style.Text
	draft string
}

func (s *stdout) delay() uint8 {
	return 1
}

func (s *stdout) renderComponent() {
	s.text.Object.SetScrollable(true).
		ScrollToEnd().
		SetDynamicColors(true)
}

func (s *stdout) reload(ctx context.Context) {
	for {
		select {
		case data := <-s.ReadMsg():
			s.draft += (data + "\n")
			s.text.Object.SetText(s.draft)
		default:
			s.renderComponent()
			return
		}
	}
}

func (s *stdout) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, s)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.text.Object, 0, 1, false)
}

func StdOut[R RenderInterface, L LoggerInterface](render R, logger L) ComponentInterface {
	return new("stdout", func() ComponentInterface {
		self := &stdout{
			LoggerInterface: logger,
			RenderInterface: render,

			text: style.NewText().AsLogger(),
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

package component

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/extension"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type stdout struct {
	LoggerWriterInterface
	RenderInterface
	logger *extension.Logger

	text  *style.Text
	draft string
}

func (s *stdout) delay() int8 {
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
		case data := <-s.logger.LogCh():
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

func StdOut(render RenderInterface, logger *extension.Logger) ComponentInterface {
	return new("stdout", func() ComponentInterface {
		self := &stdout{
			RenderInterface: render,
			logger:          logger,
			text:            style.NewText().AsLogger(),
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

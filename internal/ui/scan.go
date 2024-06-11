package ui

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/component"
	"github.com/Mad-Pixels/wf/internal/ui/frame"
	"github.com/rivo/tview"
)

func pageScan(ctx context.Context) {
	pCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	view := frame.Builder(frame.Config{
		Header: []*tview.Flex{
			component.SysInfo().FlexItem(pCtx),
			component.NetStat().FlexItem(pCtx),
		},
		Content: component.NetScan().FlexItem(pCtx),
	})

	App.SetRoot(view, true).SetFocus(view).Run()

}

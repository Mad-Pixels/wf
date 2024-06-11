package ui

import (
	"context"

	"github.com/Mad-Pixels/wf/internal/ui/component"
	"github.com/Mad-Pixels/wf/internal/ui/frame"
	"github.com/rivo/tview"
)

func Run() {
	ctx := context.Background()

	f := frame.NewFrame()
	f.Root = tview.NewFlex().SetDirection(tview.FlexRow)

	sysInfo := component.SysInfo().FlexItem(ctx)
	netStat := component.NetStat().FlexItem(ctx)

	netScan := component.NetScan(f)

	f.SetHeader([]*tview.Flex{sysInfo, netStat})
	f.SetContent(netScan.FlexItem(ctx))

	if err := f.Run(); err != nil {
		panic(err)
	}
	//pageScan(ctx)

	//frame.MM(ctx)

}

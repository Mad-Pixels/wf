package component

import (
	"context"
	"os/user"
	"runtime"

	"github.com/Mad-Pixels/wf"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type sysInfo struct {
	LoggerInterface
	RenderInterface

	usr   string
	uid   string
	table *style.Table
}

func (s *sysInfo) delay() uint8 {
	return 254
}

func (s *sysInfo) renderComponent() {
	s.table.AddCellTitle(0, 0, "Version:")
	s.table.AddCellText(0, 1, wf.Version)
	s.table.AddCellTitle(1, 0, "OS:")
	s.table.AddCellText(1, 1, runtime.GOOS)
	s.table.AddCellTitle(2, 0, "Arch:")
	s.table.AddCellText(2, 1, runtime.GOARCH)
	s.table.AddCellTitle(3, 0, "User:")
	s.table.AddCellText(3, 1, s.usr)
	s.table.AddCellTitle(4, 0, "UID:")
	s.table.AddCellText(4, 1, s.uid)
}

func (s *sysInfo) reload(ctx context.Context) {
	defer s.renderComponent()
	var (
		uid = "n/a"
		usr = "n/a"
	)
	info, err := user.Current()
	if err != nil {
		s.WriteMsg(err.Error())
		s.usr = usr
		s.uid = uid
		return
	}
	s.usr = info.Username
	s.uid = info.Uid
}

func (s *sysInfo) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, s)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.table.Object, 0, 1, false)
}

func SysInfo[R RenderInterface, L LoggerInterface](render R, logger L) ComponentInterface {
	return new("sysinfo", func() ComponentInterface {
		self := &sysInfo{
			LoggerInterface: logger,
			RenderInterface: render,

			table: style.NewTable(),
		}
		self.reload(context.Background())
		self.renderComponent()
		return self
	})
}

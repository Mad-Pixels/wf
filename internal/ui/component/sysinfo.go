package component

import (
	"context"
	"os/user"
	"runtime"

	"github.com/Mad-Pixels/wf"
	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/Mad-Pixels/wf/internal/ui/style"
	"github.com/rivo/tview"
)

type sysInfo struct {
	*binding.Synk

	usr   string
	uid   string
	table *style.Table
}

func (s *sysInfo) delay() int8 {
	return 10
}

func (n *sysInfo) triggerAppDraw() {
	n.TriggerAppDraw()
}

func (s *sysInfo) draw() {
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
	defer s.draw()
	var (
		uid = "n/a"
		usr = "n/a"
	)
	info, err := user.Current()
	if err != nil {
		s.PutLog(err.Error())
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

func SysInfo(synk *binding.Synk) ComponentInterface {
	return new("sysinfo", func() ComponentInterface {
		self := &sysInfo{
			Synk:  synk,
			table: style.NewTable(),
		}
		self.reload(context.Background())
		self.draw()
		return self
	})
}

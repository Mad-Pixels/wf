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

func SysInfo(synk *binding.Synk) ComponentInterface {
	self := &sysInfo{
		Synk:  synk,
		table: style.BaseTable(),
	}
	self.reload(context.Background())
	self.draw()
	return self
}

type sysInfo struct {
	usr   string
	uid   string
	table *tview.Table
	*binding.Synk
}

func (s *sysInfo) FlexItem(ctx context.Context) *tview.Flex {
	go schedule(ctx, s)
	return tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(s.table, 0, 1, false)
}

func (s *sysInfo) delay() int8 {
	return 10
}

func (s *sysInfo) draw() {
	s.table.SetCell(0, 0, style.CellTitle("Version:"))
	s.table.SetCell(0, 1, style.CellText(wf.Version))
	s.table.SetCell(1, 0, style.CellTitle("OS:"))
	s.table.SetCell(1, 1, style.CellText(runtime.GOOS))
	s.table.SetCell(2, 0, style.CellTitle("Arch:"))
	s.table.SetCell(2, 1, style.CellText(runtime.GOARCH))
	s.table.SetCell(3, 0, style.CellTitle("User:"))
	s.table.SetCell(3, 1, style.CellText(s.usr))
	s.table.SetCell(4, 0, style.CellTitle("UID:"))
	s.table.SetCell(4, 1, style.CellText(s.uid))
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

func (n *sysInfo) triggerAppDraw() {
	n.TriggerAppDraw()
}

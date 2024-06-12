package component

import (
	"context"
	"os/user"

	"github.com/Mad-Pixels/wf/internal/ui/binding"
	"github.com/rivo/tview"
)

func SysInfo(synk *binding.Synk) ComponentInterface {
	self := &sysInfo{
		Synk:  synk,
		table: tview.NewTable(),
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
	s.table.SetCell(0, 0, tview.NewTableCell("User:"))
	s.table.SetCell(0, 1, tview.NewTableCell(s.usr))
	s.table.SetCell(1, 0, tview.NewTableCell("UID:"))
	s.table.SetCell(0, 1, tview.NewTableCell(s.uid))
}

func (s *sysInfo) reload(ctx context.Context) {
	var (
		uid = "error"
		usr = "error"
	)
	if u, err := user.Current(); err == nil {
		usr = u.Username
		uid = u.Uid
	}
	s.uid = uid
	s.usr = usr
	s.draw()
}

func (n *sysInfo) triggerAppDraw() {
	n.TriggerAppDraw()
}

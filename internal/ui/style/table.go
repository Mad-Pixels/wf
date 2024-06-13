package style

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	Object *tview.Table
	title  string
}

func NewTable() *Table {
	return &Table{
		Object: tview.NewTable(),
	}
}

func (t *Table) WithTitle(title string) *Table {
	t.title = title
	t.Object.
		SetTitle(fmt.Sprintf(" %s ", t.title)).
		SetTitleColor(ColorTitle)
	return t
}

func (t *Table) WithCount(count int) *Table {
	if t.title == "" {
		t.Object.
			SetTitle(fmt.Sprintf("[%d]", count))
		return t
	}
	t.Object.
		SetTitle(fmt.Sprintf(" %s [red][%d][white] ", t.title, count))
	return t
}

func (t *Table) AsContent() *Table {
	t.Object.
		SetSelectable(true, false).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(ColorContent).
		SetBorder(true)
	return t
}

func (t *Table) AddCellContent(r, c int, value string) {
	t.addCell(r, c, value, ColorContent)
}

func (t *Table) AddCellSecondary(r, c int, value string) {
	t.addCell(r, c, value, ColorSecondary)
}

func (t *Table) AddCellPrimary(r, c int, value string) {
	t.addCell(r, c, value, ColorPrimary)
}

func (t *Table) AddCellTitle(r, c int, value string) {
	t.addCell(r, c, value, ColorTitle)
}

func (t *Table) AddCellText(r, c int, value string) {
	t.addCell(r, c, value, ColorText)
}

func (t *Table) AddCellHeader(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			strings.ToUpper(value),
		).
			SetAlign(tview.AlignLeft).
			SetTextColor(ColorTitle).
			SetSelectable(false).
			SetExpansion(1),
	)
}

func (t *Table) addCell(r, c int, value string, color tcell.Color) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(color).
			SetAlign(tview.AlignLeft),
	)
}

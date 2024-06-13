package style

import (
	"strings"

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

func (t *Table) AddCellContent(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(ColorContent).
			SetAlign(tview.AlignLeft),
	)
}

func (t *Table) AddCellSecondary(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(ColorSecondary).
			SetAlign(tview.AlignLeft),
	)
}

func (t *Table) AddCellPrimary(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(ColorPrimary).
			SetAlign(tview.AlignLeft),
	)
}

func (t *Table) AddCellTitle(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(ColorTitle).
			SetAlign(tview.AlignLeft),
	)
}

func (t *Table) AddCellText(r, c int, value string) {
	t.Object.SetCell(
		r,
		c,
		tview.NewTableCell(
			value,
		).
			SetTextColor(ColorText).
			SetAlign(tview.AlignLeft),
	)
}

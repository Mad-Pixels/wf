package style

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Table ...
type Table struct {
	Object    *tview.Table
	title     string
	expansion bool
}

// NewTable return custom Table object with predifined styles.
func NewTable() *Table {
	return &Table{
		Object: tview.NewTable().
			ScrollToBeginning(),
	}
}

// WithTitle add title to table.
func (t *Table) WithTitle(title string) *Table {
	t.title = title
	t.Object.
		SetTitle(fmt.Sprintf(" %s ", t.title)).
		SetTitleColor(ColorTitle)
	return t
}

// WithExpansion set width: 100%.
func (t *Table) WithExpansion() *Table {
	t.expansion = true
	return t
}

// WithCount add to table title count of rows.
func (t *Table) WithCount(count int) *Table {
	switch {
	case t.title == "":
		t.Object.SetTitle(fmt.Sprintf("[%d]", count))
	default:
		t.Object.SetTitle(fmt.Sprintf(" %s [red][%d][white] ", t.title, count))
	}
	return t
}

// AsContent add custom styles for table.
func (t *Table) AsContent() *Table {
	t.Object.
		SetSelectable(true, false).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(ColorContent).
		SetBorder(true)
	return t
}

// AddCellContent add item [row, col] as table value with Content color style.
func (t *Table) AddCellContent(r, c int, value string) {
	t.addCell(r, c, value, ColorContent)
}

// AddCellSecondary add item [row, col] as table value with Secondary color style.
func (t *Table) AddCellSecondary(r, c int, value string) {
	t.addCell(r, c, value, ColorSecondary)
}

// AddCellPrimary add item [row, col] as table value with Primary color style.
func (t *Table) AddCellPrimary(r, c int, value string) {
	t.addCell(r, c, value, ColorPrimary)
}

// AddCellTitle add item [row, col] as table value with Title color style.
func (t *Table) AddCellTitle(r, c int, value string) {
	t.addCell(r, c, value, ColorTitle)
}

// AddCellText add item [row, col] as table value with Default color style.
func (t *Table) AddCellText(r, c int, value string) {
	t.addCell(r, c, value, ColorText)
}

// AddCellHeader add item [row, col] as table Header.
func (t *Table) AddCellHeader(r, c int, value string) {
	cell := tview.NewTableCell(strings.ToUpper(value)).
		SetAlign(tview.AlignLeft).
		SetTextColor(ColorTitle).
		SetSelectable(false)
	if t.expansion {
		cell.SetExpansion(1)
	}

	t.Object.SetCell(r, c, cell)
}

func (t *Table) addCell(r, c int, value string, color tcell.Color) {
	cell := tview.NewTableCell(value).
		SetAlign(tview.AlignLeft).
		SetTextColor(color)
	if t.expansion {
		cell.SetExpansion(1)
	}

	t.Object.SetCell(r, c, cell)
}

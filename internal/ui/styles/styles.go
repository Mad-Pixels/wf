package styles

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ColorTitle     = tcell.ColorOrange
	ColorText      = tcell.ColorWhiteSmoke
	ColorPrimary   = tcell.ColorDeepSkyBlue
	ColorSecondary = tcell.ColorDarkGray
	ColorContent   = tcell.ColorLightSkyBlue
)

func ContentTable() *tview.Table {
	table := tview.NewTable()
	table.
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(ColorContent).
		SetBorder(true)
	return table
}

func BaseTable() *tview.Table {
	table := tview.NewTable()
	return table
}

func CellTitle(val string) *tview.TableCell {
	return tview.
		NewTableCell(val).
		SetTextColor(ColorTitle).
		SetAlign(tview.AlignLeft)
}

func CellText(val string) *tview.TableCell {
	return tview.
		NewTableCell(val).
		SetTextColor(ColorText).
		SetAlign(tview.AlignLeft)
}

func CellPrimary(val string) *tview.TableCell {
	return tview.
		NewTableCell(val).
		SetTextColor(ColorPrimary).
		SetAlign(tview.AlignLeft)
}

func CellSecondary(val string) *tview.TableCell {
	return tview.
		NewTableCell(val).
		SetTextColor(ColorSecondary).
		SetAlign(tview.AlignLeft)
}

func CellContent(val string) *tview.TableCell {
	return tview.
		NewTableCell(val).
		SetTextColor(ColorContent).
		SetAlign(tview.AlignLeft)
}

func CellHeader(val string) *tview.TableCell {
	return tview.
		NewTableCell(strings.ToUpper(val)).
		SetAlign(tview.AlignLeft).
		SetTextColor(ColorTitle).
		SetSelectable(false).
		SetExpansion(1)
}

package ui

import (
	"context"

	"github.com/rivo/tview"
)

var App = tview.NewApplication()

func Run() {
	ctx := context.Background()

	pageScan(ctx)

	//frame.MM(ctx)

}

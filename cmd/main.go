package main

import (
	"context"
	"log"

	"github.com/Mad-Pixels/wf/internal/ui"
)

func main() {
	if err := ui.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

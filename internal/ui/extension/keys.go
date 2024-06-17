package extension

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

// Hotkeys is a structure for storing hot key data used in a GUI helper and hot key listener.
// Expects the HotKeys structure to be present in the provided context.
type Keys struct {
	// Action is a function to be executed when the hot key is activated.
	Action func(ctx context.Context)

	// Description is the text describing the hot key's purpose in the GUI.
	Description string

	// Shortcut is the key code in the format of tcell.Key.
	Shortcut tcell.Key
}

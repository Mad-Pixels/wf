package keys

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

// Hotkeys is a structure for storing hot key data used in a GUI helper and hot key listener.
// Expects the HotKeys structure to be present in the provided context.
type Keys struct {
	// Action is a function to be executed when the hot key is activated.
	action func(ctx context.Context)

	// Description is the text describing the hot key's purpose in the GUI.
	description string

	// Shortcut is the key code in the format of tcell.Key.
	shortcut tcell.Key
}

func (k Keys) Action() func(ctx context.Context) {
	return k.action
}

func (k Keys) Description() string {
	return k.description
}

func (k Keys) Shortcut() tcell.Key {
	return k.shortcut
}

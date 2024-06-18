package component

import (
	"context"
	"sync"
	"time"

	"github.com/rivo/tview"
)

// RenderInterface defines the interface for application frames render process.
type RenderInterface interface {
	DrawRootFrame()
}

// LoggerWriterInterface defines the interface for application log process.
type LoggerInterface interface {
	WriteMsg(string)
	ReadMsg() chan string
}

// ComponentInterface defines the interface for UI components.
type ComponentInterface interface {
	LoggerInterface
	RenderInterface

	FlexItem(context.Context) *tview.Flex

	reload(context.Context)
	renderComponent()
	delay() int8
}

var (
	// Map to hold single instances of types implementing ComponentInterface.
	impl = make(map[string]ComponentInterface)

	// global mutex.
	m sync.Mutex
)

// return a singleton instance of the type implementing ComponentInterface.
func new[T ComponentInterface](key string, factory func() T) T {
	m.Lock()
	defer m.Unlock()

	if obj, exist := impl[key]; exist {
		return obj.(T)
	}
	obj := factory()

	impl[key] = obj
	return obj
}

// schedule is a wrapper for refreshing component data of the type implementing ComponentInterface.
func schedule[T ComponentInterface](ctx context.Context, t T) {
	delay := t.delay()
	if delay < 1 {
		delay = 1
	}
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.reload(ctx)
			t.DrawRootFrame()
		case <-ctx.Done():
			return
		}
	}
}

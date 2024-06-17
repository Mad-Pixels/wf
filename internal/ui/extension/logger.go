package extension

// Logger ...
type Logger struct {
	logCh chan string
}

func NewLogger(ch chan string) *Logger {
	return &Logger{
		logCh: ch,
	}
}

// Put log data.
func (e Logger) Put(data string) {
	e.logCh <- data
}

// LogCh return channel with log data.
func (e Logger) LogCh() chan string {
	return e.logCh
}

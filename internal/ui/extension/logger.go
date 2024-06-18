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
func (e Logger) WriteMsg(data string) {
	e.logCh <- data
}

// LogCh return channel with log data.
func (e Logger) ReadMsg() chan string {
	return e.logCh
}

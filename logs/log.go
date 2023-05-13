package logs

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"io"
	"os"
	"sync"
	"time"
)

var defaultLogger *Logger
var defaultFormat = &Format{
	Types:       "json",
	Color:       color.White,
	IsColorable: true,
}

func init() {
	defaultLogger = &Logger{
		output: os.Stdout,
		level:  LevelInfo,
		mux:    &sync.Mutex{},
		f:      defaultFormat,
	}
}

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
)

const (
	ColorRed    = color.Red
	ColorBlue   = color.Blue
	ColorWhite  = color.White
	ColorGrey   = color.Gray
	ColorYellow = color.Yellow
	ColorBlack  = color.Black
)

type Event struct {
	time    time.Time
	content string
	level   int
	id      int
}

type Logger struct {
	output io.Writer
	level  int
	mux    *sync.Mutex
	f      *Format
}

type Format struct {
	Types       string
	Color       color.Color
	IsColorable bool
}

func NewLogger(format ...*Format) *Logger {
	if len(format) == 0 {
		return &Logger{
			output: os.Stdout,
			level:  LevelInfo,
			mux:    &sync.Mutex{},
			f:      nil,
		}
	}
	return &Logger{
		output: os.Stdout,
		level:  LevelInfo,
		mux:    &sync.Mutex{},
		f:      format[len(format)-1],
	}
}

func SetFormat(format *Format) {
	defaultLogger.f = format
}

func SetLevel(level int) {
	defaultLogger.level = level
}

func SetOutput(output io.Writer) {
	defaultLogger.output = output
}

func SetMultiOutput(outputs ...io.Writer) {
	defaultLogger.output = io.MultiWriter(outputs...)
}

func SetColor(color color.Color) {
	defaultLogger.f.Color = color
}

func Debug(format string, args ...interface{}) {
	defaultLogger.log(LevelDebug, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.log(LevelInfo, format, args...)
}

func Warning(format string, args ...interface{}) {
	defaultLogger.log(LevelWarning, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.log(LevelError, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(LevelWarning, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *Logger) log(level int, format string, args ...interface{}) {
	e := Event{
		content: fmt.Sprintf(format, args...),
		level:   level,
		time:    time.Now(),
	}
	l.out(&e)
}

func (l *Logger) out(e *Event) {
	if l.level > e.level {
		return
	}
	if l.output == nil {
		panic("output is nil, please set output")
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	cp := map[string]interface{}{
		"time":    e.time,
		"content": e.content,
		"level":   e.level,
	}
	switch l.f.Types {
	case "json", "JSON", "Json":
		bytes, err := json.Marshal(cp)
		if err != nil {
			return
		}
		if l.f.IsColorable {
			l.output.Write([]byte(l.f.Color.Sprint(string(bytes))))
			return
		}
		l.output.Write(bytes)
	}
}

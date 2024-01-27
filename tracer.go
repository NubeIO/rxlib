package rxlib

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func newLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

type Tracer struct {
	zLog    zerolog.Logger
	Loggers map[string]*Logger
}

type Opts struct {
}

func NewTracer(opts *Opts) *Tracer {
	return &Tracer{
		zLog:    newLogger(),
		Loggers: make(map[string]*Logger),
	}
}

type Logger struct {
	AppName       string   `json:"appName"`
	UUID          string   `json:"UUID"`
	Traces        []*Trace `json:"traces"`
	MaxTracers    int      `json:"maxTracers"`
	appNameColour string
	zLog          zerolog.Logger
}

func (l *Tracer) NewLogger(appName, uuid, appNameColour string, maxTracers int) (*Logger, error) {
	err := checkTextLength(appName)
	if err != nil {

		return nil, fmt.Errorf("appName: %v actual: %d", err, len(appName))
	}
	logger := &Logger{
		AppName:       appName,
		UUID:          uuid,
		zLog:          l.zLog,
		Traces:        make([]*Trace, 0),
		MaxTracers:    maxTracers,
		appNameColour: appNameColour,
	}
	l.Loggers[appName] = logger
	return logger, nil
}

// ------------- TRACE ---------------

type Trace struct {
	LoggerName          string `json:"-"`
	LoggerUUID          string `json:"-"`
	Name                string
	NextTraceID         int
	MaxMessagePerLogger int
	Messages            []*LoggerMessage
	zLog                zerolog.Logger
	debugColour         string
	appNameColour       string
}

func (l *Logger) NewTrace(traceName, debugColour string, maxTraceCount int) *Trace {
	trace := &Trace{
		LoggerName:          l.AppName,
		LoggerUUID:          l.UUID,
		Name:                traceName,
		NextTraceID:         1,
		MaxMessagePerLogger: maxTraceCount,
		Messages:            make([]*LoggerMessage, 0, maxTraceCount),
		zLog:                l.zLog,
		appNameColour:       l.appNameColour,
		debugColour:         debugColour,
	}
	l.Traces = append(l.Traces, trace)
	return trace
}

type LoggerMessage struct {
	ID        int
	Type      string
	Message   string
	Timestamp time.Time
}

func (t *Trace) Logf(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	messageType := "INFO"
	t.newLog(zerolog.InfoLevel, logMessage, messageType)
}

func (t *Trace) Debugf(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	messageType := "DEBUG"
	t.newLog(zerolog.DebugLevel, logMessage, messageType)
}

func (t *Trace) Errorf(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	messageType := "ERROR"
	t.newLog(zerolog.ErrorLevel, logMessage, messageType)
}

func (t *Trace) newLog(level zerolog.Level, logMessage, messageType string) {

	newMessage := &LoggerMessage{
		Type:      messageType, // You can customize the log level if needed
		Message:   logMessage,
		Timestamp: time.Now(),
	}
	t.Messages = append(t.Messages, newMessage)

	// Remove oldest message if the limit is reached
	if len(t.Messages) > t.MaxMessagePerLogger {
		t.Messages = t.Messages[1:]
	}

	//text := fmt.Sprintf("app=%s trace=%s message=%s", l.LoggerName, l.Name, newMessage.Message)
	text := fmt.Sprintf("[app=%-6s] | trace=%-15s | %-10s", strings.ToUpper(t.LoggerName), fmt.Sprintf("%s:%s", t.Name, t.LoggerUUID), newMessage.Message)

	switch level {
	case zerolog.DebugLevel:
		t.zLog.Debug().Msg(t.colorizeLogMessage(text))
	case zerolog.InfoLevel:
		t.zLog.Info().Msg(t.colorizeLogMessage(text))
	case zerolog.WarnLevel:
		t.zLog.Warn().Msg(t.colorizeLogMessage(text))
	case zerolog.ErrorLevel:
		t.zLog.Error().Msg(t.colorizeLogMessage(text))
	default:
		panic("unhandled default case")
	}
}

type ColorScheme struct {
	Colors      map[string]string
	ResetColor  string
	DisableFlag bool
}

var disableColour = false

type ColorSet struct {
	Pattern *regexp.Regexp
	Color   string // Directly using a string for the ANSI color code
}

func (t *Trace) colorizeLogMessage(msg string) string {
	if disableColour {
		return msg
	}

	appNamePattern := regexp.MustCompile(`app=([^ ]+)`)
	msg = appNamePattern.ReplaceAllStringFunc(msg, func(match string) string {
		appName := strings.TrimPrefix(match, "app=")
		return t.appNameColour + appName + ColorReset
	})

	colorSets := []ColorSet{
		{regexp.MustCompile(`(trace=)`), t.debugColour},
	}

	for _, set := range colorSets {
		msg = set.Pattern.ReplaceAllStringFunc(msg, func(match string) string {
			return set.Color + match + ColorReset
		})
	}

	return msg
}

const (
	ColorBrightGray      = "\033[37;1m"
	ColorBrightRedBg     = "\033[101m"
	ColorBrightGreenBg   = "\033[102m"
	ColorBrightYellowBg  = "\033[103m"
	ColorBrightBlueBg    = "\033[104m"
	ColorBrightMagentaBg = "\033[105m"
	ColorBrightCyanBg    = "\033[106m"
	ColorBrightWhiteBg   = "\033[107m"
	ColorGray            = "\033[90m"
	ColorLightRed        = "\033[91m"
	ColorLightGreen      = "\033[92m"
	ColorLightYellow     = "\033[93m"
	ColorLightBlue       = "\033[94m"
	ColorLightMagenta    = "\033[95m"
	ColorLightCyan       = "\033[96m"
	ColorLightGray       = "\033[37;2m"
	ColorLightRedBg      = "\033[101m"

	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"

	ColorReset = "\033[0m"
)

func checkTextLength(text string) error {
	length := len(text)
	if length != 7 {
		return errors.New("text length must != 7 characters")
	}
	return nil
}

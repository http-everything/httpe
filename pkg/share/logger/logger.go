package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelError LogLevel = 0
	LogLevelInfo  LogLevel = 1
	LogLevelDebug LogLevel = 2
)

// ParseLogLevel parses the string representation of a log level and returns the corresponding LogLevel enum value.
// It uses a map to map the string value to its corresponding LogLevel value.
// If the given string value is found in the map, it returns the corresponding LogLevel value.
// Otherwise, it returns LogLevelError and an error indicating that the log level is invalid.
func ParseLogLevel(str string) (LogLevel, error) {
	var m = map[string]LogLevel{
		"error": LogLevelError,
		"info":  LogLevelInfo,
		"debug": LogLevelDebug,
	}
	if result, ok := m[str]; ok {
		return result, nil
	}
	return LogLevelError, fmt.Errorf("invalid log level: %q", str)
}

// LogLevelStr returns the string representation of the given LogLevel.
// It uses a map to map the integer value of LogLevel to its corresponding string value.
// If the given LogLevel is found in the map, it returns the string value.
// Otherwise, it returns "unknown".
func LogLevelStr(level LogLevel) (levelStr string) {
	var m = map[int]string{
		0: "ERROR",
		1: "INFO",
		2: "DEBUG",
	}
	if result, ok := m[int(level)]; ok {
		return result
	}
	return "unknown"
}

type LogOutput struct {
	File     *os.File
	filePath string
}

func NewLogOutput(filePath string) LogOutput {
	return LogOutput{
		filePath: filePath,
	}
}

func (o *LogOutput) Start() error {
	if o.filePath == "" {
		o.File = os.Stdout
		return nil
	}

	var err error
	o.File, err = os.OpenFile(o.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("can't open log file %s: %w", o.filePath, err)
	}
	return nil
}

func (o *LogOutput) Shutdown() {
	if o.File != nil && o.File != os.Stdout {
		_ = o.File.Close()
	}
}

type Logger struct {
	prefix string
	logger *log.Logger
	output LogOutput
	level  LogLevel
}

func New(prefix string, filePath string, level string) (l *Logger, err error) {
	logLevel, err := ParseLogLevel(level)
	if err != nil {
		return nil, err
	}

	logOutput := NewLogOutput(filePath)
	err = logOutput.Start()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logOutput = LogOutput{File: os.Stderr}
		} else {
			return nil, err
		}
	}

	l = NewLogger(prefix, logOutput, logLevel)
	return l, nil
}

func NewLogger(prefix string, output LogOutput, level LogLevel) *Logger {
	l := &Logger{
		prefix: prefix,
		logger: log.New(output.File, "", log.Ldate|log.Ltime),
		output: output,
		level:  level,
	}
	return l
}

func (l *Logger) Errorf(f string, args ...interface{}) {
	l.Logf(LogLevelError, f, args...)
}

func (l *Logger) Infof(f string, args ...interface{}) {
	l.Logf(LogLevelInfo, f, args...)
}

func (l *Logger) Debugf(f string, args ...interface{}) {
	l.Logf(LogLevelDebug, f, args...)
}

func (l *Logger) Logf(severity LogLevel, f string, args ...interface{}) {
	if l.level >= severity {
		if l.prefix == "-" {
			l.logger.Printf(LogLevelStr(severity)+": "+f, args...)
		} else {
			l.logger.Printf(LogLevelStr(severity)+": "+l.prefix+": "+f, args...)
		}
	}
}

func (l *Logger) Fork(prefix string, args ...interface{}) *Logger {
	// slip the parent prefix at the front
	args = append([]interface{}{l.prefix}, args...)
	ll := NewLogger(fmt.Sprintf("%s: "+prefix, args...), l.output, l.level)
	return ll
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) Shutdown() {
	l.output.Shutdown()
}

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_globalLogger *Logger
	_globalMtx    sync.RWMutex
	_logPath      string = "./internal/logs/"
)

// Create new global logger
func NewGlobal(l *Logger) func() {
	_globalMtx.Lock()
	prev := _globalLogger
	_globalLogger = l
	_globalMtx.Unlock()
	return func() { NewGlobal(prev) }
}

// Global logger
func L(logLevel ...string) *Logger {
	_globalMtx.RLock()
	defer _globalMtx.RUnlock()
	l := _globalLogger
	if l == nil {
		level := ""
		if len(logLevel) > 0 {
			level = logLevel[0]
		}
		l = New(level)
	}
	return l
}

// Interface -.
type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zap.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	filename := "error"
	if level != "" {
		filename = level
	}
	filename = filepath.Join(_logPath, fmt.Sprintf("%s.log", level))
	// Check if the file not exists
	if _, err := os.Stat(filename); err != nil {
		// Create a new file
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	// Open log file
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create file logger config
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)

	// Create writer
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	// Create a new zap logger
	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		logger: logger,
	}
}

// Debug -.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug(message)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info(message)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message)
}

// Error -.
func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error(message)
}

// Fatal -.
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatal(message)

	os.Exit(1)
}

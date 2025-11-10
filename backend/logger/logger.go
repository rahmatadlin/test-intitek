package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// FileLogger implements logger.Logger interface for Wails
// It writes logs to both file and console
type FileLogger struct {
	logDir     string
	logFile    *os.File
	fileLog    *log.Logger
	consoleLog *log.Logger
	level      logger.LogLevel
}

// NewFileLogger creates a new file logger
func NewFileLogger(logDir string, level logger.LogLevel) (*FileLogger, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("warehouse-management-%s.log", timestamp)
	logFilePath := filepath.Join(logDir, logFileName)

	// Open log file (append mode)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer for both file and console
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	fileLogger := &FileLogger{
		logDir:     logDir,
		logFile:    logFile,
		fileLog:    log.New(logFile, "", log.LstdFlags|log.Lshortfile),
		consoleLog: log.New(multiWriter, "", log.LstdFlags|log.Lshortfile),
		level:      level,
	}

	return fileLogger, nil
}

// Close closes the log file
func (l *FileLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Print logs a message at Print level
func (l *FileLogger) Print(message string) {
	l.log(0, message) // Print level is 0 (always shown)
}

// Trace logs a message at Trace level
func (l *FileLogger) Trace(message string) {
	l.log(logger.TRACE, message)
}

// Debug logs a message at Debug level
func (l *FileLogger) Debug(message string) {
	l.log(logger.DEBUG, message)
}

// Info logs a message at Info level
func (l *FileLogger) Info(message string) {
	l.log(logger.INFO, message)
}

// Warning logs a message at Warning level
func (l *FileLogger) Warning(message string) {
	l.log(logger.WARNING, message)
}

// Error logs a message at Error level
func (l *FileLogger) Error(message string) {
	l.log(logger.ERROR, message)
}

// Fatal logs a message at Fatal level and exits
func (l *FileLogger) Fatal(message string) {
	l.log(logger.ERROR+1, message) // Fatal is ERROR+1
	os.Exit(1)
}

// log writes log message if level is enabled
func (l *FileLogger) log(level logger.LogLevel, message string) {
	// Check if this level should be logged
	// Print level (0) is always shown
	if level != 0 && level < l.level {
		return
	}

	levelName := getLevelName(level)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, levelName, message)

	// Write to both file and console
	l.consoleLog.Output(3, logMessage)
}

// getLevelName returns string representation of log level
func getLevelName(level logger.LogLevel) string {
	switch level {
	case 0:
		return "PRINT"
	case logger.TRACE:
		return "TRACE"
	case logger.DEBUG:
		return "DEBUG"
	case logger.INFO:
		return "INFO"
	case logger.WARNING:
		return "WARNING"
	case logger.ERROR:
		return "ERROR"
	case logger.ERROR + 1:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Log level constants (matching Wails logger levels)
// Wails logger levels: TRACE=1, DEBUG=2, INFO=3, WARNING=4, ERROR=5, FATAL=6
// We use the constants directly from Wails logger package

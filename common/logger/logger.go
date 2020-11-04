package logger

import (
	"errors"
	"log"
	"os"
)

type appLogger struct {
	level     int
	directory string
	logError  *log.Logger
	logInfo   *log.Logger
	logDebug  *log.Logger
}

//AppLogger - Logger interface
type AppLogger interface {
	SetLevel(logLevel int) error
	Level() int
	Error(a ...interface{})
	Info(a ...interface{})
	Debug(a ...interface{})
}

var logger AppLogger = nil

//SetLevel - Sets current log level
func (log *appLogger) SetLevel(logLevel int) error {

	if logLevel > 2 {
		return errors.New("Invalid log level value")
	}
	log.level = logLevel
	return nil
}

//Level - Gets log level
func (log *appLogger) Level() int {

	return log.level
}

//Error - Logs error
func (log *appLogger) Error(a ...interface{}) {
	if log.level >= 0 {
		log.logError.Println(a...)
	}

}

//Info - Logs info
func (log *appLogger) Info(a ...interface{}) {
	if log.level >= 1 {
		log.logInfo.Println(a...)
	}

}

//Debug - Logs debug
func (log *appLogger) Debug(a ...interface{}) {
	if log.level == 2 {
		log.logDebug.Println(a...)
	}

}

//Get - gets instance of a AppLogger
func Get() AppLogger {
	return logger
}

//NewTestLogger - returns instance of a logger for tests
func NewTestLogger() AppLogger {

	if logger == nil {
		logger = &appLogger{
			level:    2,
			logDebug: log.New(os.Stderr, "DEBUG:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
			logInfo:  log.New(os.Stdout, "INFO:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
			logError: log.New(os.Stdout, "ERROR:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
		}
	}

	return logger

}

// NewLogger - get single instance of logger
func NewLogger(directory string, logLevel int) AppLogger {

	if logger == nil {
		logger = &appLogger{
			level:    logLevel,
			logDebug: log.New(os.Stderr, "DEBUG:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
			logInfo:  log.New(os.Stdout, "INFO:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
			logError: log.New(os.Stdout, "ERROR:", log.Ldate|log.LUTC|log.Ltime|log.Lmicroseconds),
		}
	}

	return logger

}

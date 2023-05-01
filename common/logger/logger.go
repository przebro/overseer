package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerConfiguration - configuration for logger
type LoggerConfiguration struct {
	LogLevel     int    `json:"logLevel" validate:"min=0,max=5"`
	SizeLimit    int    `json:"sizeLimit" validate:"gte=1024"`
	LogDirectory string `json:"logDirectory"`
	FilePrefix   string `json:"prefix"`
}

// LogConfiguration - provides required log settings
type LogConfiguration interface {
	Level() int
	Directory() string
	Prefix() string
	Limit() int
}

// Level - returns current log level
func (l LoggerConfiguration) Level() int {
	return l.LogLevel
}

// Directory - returns current log directory
func (l LoggerConfiguration) Directory() string {

	return l.LogDirectory
}

// Prefix - returns log file prefix
func (l LoggerConfiguration) Prefix() string {
	return l.FilePrefix
}

// Limit - returns log size limit
func (l LoggerConfiguration) Limit() int {
	return l.SizeLimit

}

func Configure(name string, conf LoggerConfiguration) {

	zerolog.SetGlobalLevel(zerolog.Level(conf.LogLevel))
	log.Logger = log.With().Str("app", name).Logger()
	if conf.FilePrefix != "" {
		rr, _ := newRotateSync(conf.LogDirectory, conf.FilePrefix, conf.SizeLimit*1024)
		multi := zerolog.MultiLevelWriter(os.Stderr, rr)
		log.Logger = zerolog.New(multi).With().Str("app", name).Logger()
	}

}

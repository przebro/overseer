package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//AppLogger - Logger interface
type AppLogger interface {
	Error(a ...interface{})
	Info(a ...interface{})
	Debug(a ...interface{})
	Desugar() *zap.Logger
}

//LogConfiguration - provides required log settings
type LogConfiguration interface {
	Level() int
	Directory() string
	Prefix() string
	Limit() int
}

//NewTestLogger - returns instance of a logger for tests
func NewTestLogger() AppLogger {

	lg, _ := zap.NewDevelopment()

	return lg.Sugar()
}

// NewLogger - get a new instance of AppLogger
func NewLogger(conf LogConfiguration) (*zap.SugaredLogger, error) {

	var lg *zap.Logger
	var rsync *fileRotateSync
	var err error
	var jfile zapcore.Core
	var cons zapcore.Core

	zapconfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "lvl",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	cEncoder := zapcore.NewConsoleEncoder(zapconfig)

	f := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {

		return lvl >= zapcore.Level(conf.Level()-1)
	})

	if conf.Prefix() != "" {

		jEncoder := zapcore.NewJSONEncoder(zapconfig)

		if rsync, err = newRotateSync(conf.Directory(), conf.Prefix(), conf.Limit()*1024); err != nil {
			return nil, fmt.Errorf("fatal, can't create WriteSyncer:%v", err)
		}

		jfile = zapcore.NewCore(jEncoder, rsync, f)

	} else {
		jfile = zapcore.NewNopCore()
	}

	cons = zapcore.NewCore(cEncoder, os.Stdout, f)
	core := zapcore.NewTee(jfile, cons)

	lg = zap.New(core, zap.AddCaller())

	return lg.Sugar(), nil

}

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.SugaredLogger

func init() {
	cores := NewCores()
	Logger = zap.New(cores, zap.AddCaller()).Sugar()
	defer Logger.Sync()
}

func NewCores() zapcore.Core {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewConsoleEncoder(config)

	mainLog := zapcore.NewCore(enc, zapcore.AddSync(&lumberjack.Logger{
		Filename: "./log/main.log",
	}), zapcore.InfoLevel)
	verbLog := zapcore.NewCore(enc, zapcore.AddSync(&lumberjack.Logger{
		Filename: "./log/verbose.log",
		MaxSize:  5,
	}), zapcore.DebugLevel)
	termLog := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	return zapcore.NewTee(mainLog, verbLog, termLog)
}

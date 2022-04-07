package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger zap.SugaredLogger

func initLogger() {
	logConf := zap.NewProductionConfig()
	logConf.Encoding = "console"
	logConf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConf.DisableStacktrace = true
	unsugared, err := logConf.Build()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	logger = *unsugared.Sugar()
}

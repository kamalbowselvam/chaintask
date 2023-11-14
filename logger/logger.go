package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func init() {
    var err error
    // config := zap.NewProductionConfig()
    enccoderConfig := zap.NewDevelopmentEncoderConfig()
    enccoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

    //zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
    //zapcore.AddSync(colorable.NewColorableStdout())
    //zapcore.NewConsoleEncoder(enccoderConfig)

    //enccoderConfig.StacktraceKey = "" // to hide stacktrace info
    //config.EncoderConfig = enccoderConfig
    //zapLog, err = config.Build(zap.AddCallerSkip(1))


    zapLog = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(enccoderConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

    if err != nil {
        panic(err)
    }
}

func Info(message string, fields ...zap.Field) {
    zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
    zapLog.Debug(message, fields...)
}


func Warn(message string, fields ...zap.Field) {
    zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
    zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
    zapLog.Fatal(message, fields...)
}
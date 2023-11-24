package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	var err error

	/*
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"

	enccoderConfig := zap.NewDevelopmentEncoderConfig()

	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	zapcore.AddSync(colorable.NewColorableStdout())
	zapcore.NewConsoleEncoder(enccoderConfig)

	//enccoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = enccoderConfig
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	*/
	//zapLog = zap.New(zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(enccoderConfig),
	//	zapcore.AddSync(colorable.NewColorableStdout()),
	//	zapcore.DebugLevel,
	//))

	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	sugar = zapLog.Sugar()

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

func Panic(message string, fields ...zap.Field) {
	zapLog.Panic(message, fields...)
}

func DPanic(message string, fields ...zap.Field) {
	zapLog.DPanic(message, fields...)
}

func Warnf(message string, args ...interface{}) {
	sugar.Warnf(message, args)
}

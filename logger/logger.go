package logger

import (
	"runtime/debug"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	var err error

	stdout := zapcore.AddSync(colorable.NewColorableStdout())

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
	})
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var gitRevision string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}
	// log to multiple destinations (console and file)
	// extra fields are added to the JSON output alone
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), stdout, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, file, zapcore.DebugLevel).
			With(
				[]zapcore.Field{
					zap.String("git_revision", gitRevision),
					zap.String("go_version", buildInfo.GoVersion),
				},
			),
	)

	zapLog = zap.New(core)

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

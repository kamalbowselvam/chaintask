package logger

import (
	"context"
	"os"

	"github.com/mattn/go-colorable"
	//"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger
var sugar *zap.SugaredLogger

type ctxKey struct{}

func init() {
	//var err error

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

	/*aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	sugar = zapLog.Sugar()

	if err != nil {
		panic(err)
	}*/
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder :=  zapcore.NewConsoleEncoder(encoderConfig)
	if os.Getenv("APP_ENV") == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(colorable.NewColorableStdout()),
		zap.DebugLevel,
	)
	zapLog = zap.New(core)
	sugar = zapLog.Sugar()
}

func Get() *zap.Logger {
	return zapLog
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


// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := zapLog; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
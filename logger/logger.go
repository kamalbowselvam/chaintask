package logger

import (
	"context"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger
var sugar *zap.SugaredLogger


type ctxKey struct{}

const ginCtxKey = "logger"

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

	zapLog = zap.New(core, zap.AddCallerSkip(1))

	sugar = zapLog.Sugar()

	if err != nil {
		panic(err)
	}
}

func Get() *zap.Logger{
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

func WithGinCtx(ctx *gin.Context, l *zap.Logger) {
    if l == FromGinCtx(ctx){
		return;
	}
	ctx.Set(ginCtxKey, l)
}


func FromGinCtx(ctx *gin.Context) *zap.Logger {
	if value, ok := ctx.Get(ginCtxKey); ok {
		if l, ok := value.(*zap.Logger); ok {
			return l
		} else if l := zapLog; l != nil {
			return l
		}
	}

	return zap.NewNop()
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

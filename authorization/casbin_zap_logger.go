package authorization

import (
	"sync/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CasbinLogger is the implementation for a CasbinLogger using golang log.
type CasbinLogger struct {
	enabled int32
	logger  zap.Logger
}

type stringMatrix [][]string

func (matrix stringMatrix) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, vector := range matrix {
		if err := enc.AppendArray(zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, item := range vector {
				enc.AppendString(item)
			}
			return nil
		})); err != nil {
			return err
		}
	}
	return nil
}

func NewCasbinLogger(enabled bool, logger zap.Logger) *CasbinLogger {
	var enable int32;
	if enabled{
		enable = 1
	}else{
		enable = 0
	}
	return &CasbinLogger{
		enabled: enable,
		logger:  logger,
	}
}

func (l *CasbinLogger) EnableLog(enable bool) {
	var enab int32
	if enable {
		enab = 1
	}
	atomic.StoreInt32(&l.enabled, enab)
}

func (l *CasbinLogger) IsEnabled() bool {
	return atomic.LoadInt32(&l.enabled) == 1
}

func (l *CasbinLogger) LogModel(model [][]string) {
	if !l.IsEnabled() {
		return
	}

	l.logger.Info("LogModel", zap.Array("model", stringMatrix(model)))
}

func (l *CasbinLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !l.IsEnabled() {
		return
	}

	l.logger.Info(
		"LogEnforce",
		zap.String("matcher", matcher),
		zap.Array("request", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, v := range request {
				if err := enc.AppendReflected(v); err != nil {
					return err
				}
			}
			return nil
		})),
		zap.Bool("result", result),
		zap.Array("explains", stringMatrix(explains)),
	)
}

func (l *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if !l.IsEnabled() {
		return
	}

	l.logger.Info("LogPolicy", zap.Object("policy", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for k, v := range policy {
			if err := enc.AddArray(k, stringMatrix(v)); err != nil {
				return err
			}
		}
		return nil
	})))
}

func (l *CasbinLogger) LogRole(roles []string) {
	if !l.IsEnabled() {
		return
	}

	l.logger.Info("LogRole", zap.Strings("roles", roles))
}

func (l *CasbinLogger) LogError(err error, msg ...string) {
	if !l.IsEnabled() {
		return
	}

	l.logger.Error("LogError", zap.Error(err), zap.Strings("msg", msg))
}
package authorization

import (
	"os"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/casbin/casbin/v2/util"
	"go.uber.org/zap"
)

type Loaders struct {
	Adapter  *pgadapter.Adapter
	Enforcer *casbin.Enforcer
	Logger zap.Logger
}

type FakeLoader struct {
	Adapter  *fileadapter.Adapter
	Enforcer *casbin.Enforcer
	Logger   zap.Logger
}

var singleInstance *Loaders

func Load(source string, conf string, logger zap.Logger) (*Loaders, error) {
	if singleInstance == nil {
		adapter, err := pgadapter.NewAdapter(source, "d8p077445kq414")
		if err != nil {
			panic(err)
		}
		// Load model configuration file and policy store adapter
		_, err = os.Stat(conf)
		if err != nil {
			conf = "." + conf
		}
		enforcer, err := casbin.NewEnforcer(conf, adapter)
		if err != nil {
			logger.Fatal("", zap.Error(err))
			return nil, err
		}
		enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch)
		enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.RegexMatch)
		enforcer.EnableLog(true)
		casbin_logger := NewCasbinLogger(true, logger)
		enforcer.SetLogger(casbin_logger)
		singleInstance = &Loaders{
			Adapter:  adapter,
			Enforcer: enforcer,
			Logger: logger,
		}
	}
	return singleInstance, nil
}

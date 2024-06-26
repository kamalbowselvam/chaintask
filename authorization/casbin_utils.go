package authorization

import (
	"os"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/casbin/casbin/v2/util"
	"github.com/go-pg/pg/v10"
	"github.com/kamalbowselvam/chaintask/logger"
	"go.uber.org/zap"
)

type Loaders struct {
	Adapter  *pgadapter.Adapter
	Enforcer *casbin.Enforcer
}

type FakeLoader struct {
	Adapter  *fileadapter.Adapter
	Enforcer *casbin.Enforcer
}

var singleInstance *Loaders

func Load(source string, conf string) (*Loaders, error) {
	if singleInstance == nil {

		opts, _ := pg.ParseURL(source)
		db := pg.Connect(opts)
		//defer db.Close()

		adapter, err := pgadapter.NewAdapterByDB(db)
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
		casbin_logger := NewCasbinLogger(true, logger.Get())
		enforcer.SetLogger(casbin_logger)

		singleInstance = &Loaders{
			Adapter:  adapter,
			Enforcer: enforcer,

		}
	}
	return singleInstance, nil
}

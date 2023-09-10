package authorization

import (
	"log"
	"os"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/casbin/casbin/v2/util"
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
		adapter, err := pgadapter.NewAdapter(source)
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
			log.Fatal(err)
			return nil, err
		}
		enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch)
		enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.RegexMatch)
		singleInstance = &Loaders{
			Adapter:  adapter,
			Enforcer: enforcer,
		}
	}
	return singleInstance, nil
}

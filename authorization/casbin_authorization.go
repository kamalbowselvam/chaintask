package authorization

import (
	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/kamalbowselvam/chaintask/token"
)

type CasbinAuthorization struct {
	Adapter  *pgadapter.Adapter
	Enforcer *casbin.Enforcer
}

func NewCasbinAuthorization(loader Loaders) (AuthorizationService, error) {
	authorize := &CasbinAuthorization{
		Adapter:  loader.Adapter,
		Enforcer: loader.Enforcer,
	}

	return authorize, nil
}
func (authorize *CasbinAuthorization) Enforce(sub token.Payload, obj string, act string) (bool, error) {
	return false, nil
}

func (authorize *CasbinAuthorization) LoadPolicy() error {
	return authorize.Enforcer.LoadPolicy()
}

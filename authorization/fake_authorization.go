package authorization

import (
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/kamalbowselvam/chaintask/token"
)

type FakeCasbinAuthorization struct {
	Adapter  *fileadapter.Adapter
	Enforcer *casbin.Enforcer
}

func NewFakeCasbinAuthorization(loader FakeLoader) (AuthorizationService, error) {
	authorize := &FakeCasbinAuthorization{
		Adapter:  loader.Adapter,
		Enforcer: loader.Enforcer,
	}

	return authorize, nil
}
func (authorize *FakeCasbinAuthorization) Enforce(sub *token.Payload, obj string, act string) (bool, error) {
	// No enforcement is done, so that previous tests do not break
	// could be replaced by a setupAuthorirization function for each test cases ?
	return true, nil
}
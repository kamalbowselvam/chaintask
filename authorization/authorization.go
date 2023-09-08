package authorization

import "github.com/kamalbowselvam/chaintask/token"

type AuthorizationService interface {
	LoadPolicy() error
	Enforce(token.Payload, string, string) (bool, error)
}

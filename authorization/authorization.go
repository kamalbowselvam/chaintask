package authorization

import "github.com/kamalbowselvam/chaintask/token"

type AuthorizationService interface {
	Enforce(*token.Payload, string, string) (bool, error)
}

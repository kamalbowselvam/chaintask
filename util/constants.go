package util

import "strings"

const (
	PIPE = "|"
	DEFAULT_ADMIN = "kselvamADMIN"
)

var (
	ROLES = map[int64](string){
		1: "CLIENT",
		2: "MANAGER",
		3: "ADMIN",
	}

	ROLES_INVERT = map[string](int64){
		"CLIENT":        1,
		"MANAGER":       2,
		"ADMIN":         3,
	}
)

func GenerateRoleString(roles ...string) string{
	return strings.Join(roles, PIPE)
}

package util

import "strings"

const (
	PIPE = "|"
	DEFAULT_SUPER_ADMIN = "chaintask"
)



var (
	ROLES = map[int64](string){
		1: "CLIENT",
		2: "RESPOSIBLE",
		3: "MANAGER",
		4: "ADMIN",
	}
)

func GenerateRoleString(roles ...string) string{
	return strings.Join(roles, PIPE)
}

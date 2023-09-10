package util

const (
	READ   = "read"
	WRITE  = "write"
	DELETE = "delete"
	UPDATE = "update"
	TASK   = "tasks"
	PROJECT = "projects"
	USER = "users"
	PIPE = "|"
	DEFAULT_ADMIN = "kselvamADMIN"
)

var (
	ROLES = map[int64](string){
		1: "USER",
		2: "WORKS_MANAGER",
		3: "ADMIN",
	}

	ROLES_INVERT = map[string](int64){
		"USER":          1,
		"WORKS_MANAGER": 2,
		"ADMIN":         3,
	}
)

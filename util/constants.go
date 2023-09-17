package util

const (
	READ   = "read"
	WRITE  = "write"
	DELETE = "delete"
	UPDATE = "update"
	TASK   = "task"
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

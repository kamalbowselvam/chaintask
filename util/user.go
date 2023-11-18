package util

// Constants for all supported currencies
const (
	ADMIN = "ADMIN"
	MANAGER = "MANAGER"
	RESPONSIBLE = "RESPONSIBLE"
	CLIENT = "CLIENT"
	SUPERADMIN = "SUPERADMIN"
)

// IsSupportedRolereturns true if the role is supported
func IsSupportedRole(role string) bool {
	switch role {
	case ADMIN, MANAGER, RESPONSIBLE, CLIENT, SUPERADMIN:
		return true
	}
	return false
}
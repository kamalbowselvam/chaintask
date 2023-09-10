package authorization

type PolicyManagementService interface {
	RemoveUserPolicies(string) error
	RemoveProjectPolicies(int64, string, string) error
	CreateAdminPolicies(string) error
	CreateUserPolicies(string, string) error
	CreateProjectPolicies(int64, string, string) error
}

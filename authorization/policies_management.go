package authorization

type PolicyManagementService interface {
	RemoveUserPolicies(string) error
	RemoveProjectPolicies(int64, string, string) error
	CreateAdminPolicies(string) error
	CreateUserPolicies(int64, string, string) error
	CreateProjectPolicies(int64, string, string) error
}

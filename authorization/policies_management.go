package authorization

type PolicyManagementService interface {
	RemoveUserPolicies(string) error
	RemoveTaskPolicies(int64) error
	RemoveProjectPolicies(int64) error
	RemoveAdminPolicies(string) error
	CreateAdminPolicies(string) error
	CreateUserPolicies(string, string) error
	CreateTaskPolicies(int64, int64, string, string) error
	CreateProjectPolicies(int64, string, string) error
}

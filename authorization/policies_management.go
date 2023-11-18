package authorization

type PolicyManagementService interface {
	RemoveUserPolicies(string) error
	RemoveProjectPolicies(int64, string, string) error
	CreateAdminPolicies(string) error
	CreateUserPolicies(string, string) error
	CreateProjectPolicies(int64, string, string) error
	CreateTaskPolicies(int64, int64, string) error
	RemoveTaskPolicies(int64, int64, string) error
	CreateCompanyPolicies(int64, string) error
	RemoveCompanyPolicies(int64, string) error
}

package authorization

type PolicyManagementService interface {
	RemoveUserPolicies(string) error
	RemoveProjectPolicies(int64, string, string, int64) error
	CreateSuperAdminPolicies(string) error
	CreateAdminPolicies(string, int64) error
	CreateUserPolicies(string, string, int64) error
	CreateProjectPolicies(int64, string, string, int64) error
	CreateTaskPolicies(int64, int64, string, int64) error
	RemoveTaskPolicies(int64, int64, string, int64) error
	CreateCompanyPolicies(int64, string) error
	RemoveCompanyPolicies(int64, string) error
}

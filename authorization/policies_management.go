package authorization

type PolicyManagementService interface{
	AddPolicy()
	AddPolicies()
	RemovePolicy()
	RemovePolicies()
	RemoveUserPolicies()
	RemoveTaskPolicies()
	RemoveProjectPolicies()
	RemoveAdminPolicies()
	CreateAdminPolicies()
	CreateUserPolicies()
	CreateTaskPolicies()
	CreateProjectPolocies()
}
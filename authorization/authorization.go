package authorization
type AuthorizationService interface {
	CreateEnforcer()
	CreateAdapter()
	LoadAdminPolicies()
	AddPolicy()
	AddPolicies()
	RemovePolicy()
	RemovePolicies()
	Enforce()
}
package authorization

import (
	"log"
	"strings"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/kamalbowselvam/chaintask/util"
)

type CasbinManagement struct {
	Adapter  *pgadapter.Adapter
	Enforcer *casbin.Enforcer
}

func NewCasbinManagement(loader Loaders) (PolicyManagementService, error) {
	management := &CasbinManagement{
		Adapter:  loader.Adapter,
		Enforcer: loader.Enforcer,
	}

	return management, nil
}
func (management *CasbinManagement) CreateAdminPolicies() {
	// FIXME Change to HTTP verb ? it should be better
	rights := strings.Join([]string{util.READ, util.WRITE, util.UPDATE, util.DELETE}, util.PIPE)
	rules := [][]string{
		// FIXME
		{"p", util.ROLES[3], util.TASK, "*", rights},
		{"p", util.ROLES[3], util.PROJECT, "*", rights},
		{"p", util.ROLES[3], util.USER, "*", rights},
	}

	_, err := management.Enforcer.AddPoliciesEx(rules)
	if err != nil {
		log.Fatal(err)
	}
	_, err = management.Enforcer.AddGroupingPoliciesEx([][]string{{"g", util.ROLES[3], "kselvamADMIN"}})
	if err != nil {
		log.Fatal(err)
	}
}
func (management *CasbinManagement) AddPolicy()
func (management *CasbinManagement) AddPolicies()
func (management *CasbinManagement) RemovePolicy()
func (management *CasbinManagement) RemovePolicies()
func (management *CasbinManagement) RemoveUserPolicies()
func (management *CasbinManagement) RemoveTaskPolicies()
func (management *CasbinManagement) RemoveProjectPolicies()
func (management *CasbinManagement) RemoveAdminPolicies()
func (management *CasbinManagement) CreateUserPolicies()
func (management *CasbinManagement) CreateTaskPolicies()
func (management *CasbinManagement) CreateProjectPolocies()

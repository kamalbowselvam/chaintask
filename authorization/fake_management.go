package authorization

import (
	"log"
	"strings"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/kamalbowselvam/chaintask/util"
)

type FakeCasbinManagement struct {
	Adapter  *fileadapter.Adapter
	Enforcer *casbin.Enforcer
}

func NewFakeCasbinManagement(loader FakeLoader) (PolicyManagementService, error) {
	management := &FakeCasbinManagement{
		Adapter:  loader.Adapter,
		Enforcer: loader.Enforcer,
	}

	return management, nil
}
func (management *FakeCasbinManagement) CreateAdminPolicies() {
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
func (management *FakeCasbinManagement) AddPolicy()
func (management *FakeCasbinManagement) AddPolicies()
func (management *FakeCasbinManagement) RemovePolicy()
func (management *FakeCasbinManagement) RemovePolicies()
func (management *FakeCasbinManagement) RemoveUserPolicies()
func (management *FakeCasbinManagement) RemoveTaskPolicies()
func (management *FakeCasbinManagement) RemoveProjectPolicies()
func (management *FakeCasbinManagement) RemoveAdminPolicies()
func (management *FakeCasbinManagement) CreateUserPolicies()
func (management *FakeCasbinManagement) CreateTaskPolicies()
func (management *FakeCasbinManagement) CreateProjectPolocies()

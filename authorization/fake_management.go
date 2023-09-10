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
func (management *FakeCasbinManagement) CreateAdminPolicies(adminName string) error {
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
	_, err = management.Enforcer.AddGroupingPoliciesEx([][]string{{"g", util.ROLES[3], adminName}})
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func (management *FakeCasbinManagement) RemoveUserPolicies(string) error
func (management *FakeCasbinManagement) RemoveTaskPolicies(int64) error
func (management *FakeCasbinManagement) RemoveProjectPolicies(int64) error
func (management *FakeCasbinManagement) RemoveAdminPolicies(string) error
func (management *FakeCasbinManagement) CreateUserPolicies(string, string) error
func (management *FakeCasbinManagement) CreateTaskPolicies(int64, int64, string, string) error
func (management *FakeCasbinManagement) CreateProjectPolicies(int64, string, string) error

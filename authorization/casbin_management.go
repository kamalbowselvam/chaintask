package authorization

import (
	"log"
	"net/http"
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
func (management *CasbinManagement) CreateAdminPolicies(adminName string) error {
	rights := strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut}, util.PIPE)
	rules := [][]string{
		{"p", util.ROLES[3], "/tasks/*", "*", rights},
		{"p", util.ROLES[3], "/users/*", "*", rights},
		{"p", util.ROLES[3], "/projects/*", "*", rights},
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
func (management *CasbinManagement) RemoveUserPolicies(string) error
func (management *CasbinManagement) RemoveTaskPolicies(int64) error
func (management *CasbinManagement) RemoveProjectPolicies(int64) error
func (management *CasbinManagement) RemoveAdminPolicies(string) error
func (management *CasbinManagement) CreateUserPolicies(string, string) error
func (management *CasbinManagement) CreateTaskPolicies(int64, int64, string, string) error
func (management *CasbinManagement) CreateProjectPolicies(int64, string, string) error

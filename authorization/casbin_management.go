package authorization

import (
	"fmt"
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
		{"p", util.ROLES[3], "/users/*", rights},
		{"p", util.ROLES[3], "/projects/*", rights},
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
func (management *CasbinManagement) RemoveUserPolicies(username string) error {
	affected, err := management.Enforcer.DeleteUser(username)
	if !affected {
		log.Fatalf("%s not present in policies", username)
	}
	return err
}
func (management *CasbinManagement) RemoveProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d", projectId)
	management.RemovePolicies(resource, client)
	management.RemovePolicies(resource, responsible)
	resource = fmt.Sprintf("/projects/%d/*", projectId)
	management.RemovePolicies(resource, client)
	management.RemovePolicies(resource, responsible)
	return nil
}
func (management *CasbinManagement) CreateUserPolicies(userId int64, username string, role string) error {
	var err error
	if role != util.ROLES[3] {
		resource := fmt.Sprintf("/users/%d", userId)
		err = management.AddPolicies(resource, username, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	} else {
		err = management.CreateAdminPolicies(username)
	}
	return err
}
func (management *CasbinManagement) CreateProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d", projectId)
	management.AddPolicies(resource, client, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	management.AddPolicies(resource, responsible, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	resource = fmt.Sprintf("/projects/%d/*", projectId)
	management.AddPolicies(resource, client, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut))
	management.AddPolicies(resource, responsible, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut))
	return nil
}
func (management *CasbinManagement) RemovePolicies(resource string, username string) error {
	policies := management.Enforcer.GetPermissionsForUser(username)
	policiesToRemove := [][]string{}
	for _, policy := range policies {
		if policy[2] == resource {
			policiesToRemove = append(policiesToRemove, policy)
		}
	}
	affected, err := management.Enforcer.RemovePolicies(policiesToRemove)
	if !affected {
		log.Fatalf("problem while removing policies %s due to %s", policiesToRemove, err)
	}
	return err
}
func (management *CasbinManagement) AddPolicies(resource string, username string, rights string) error {
	rules := [][]string{
		{"p", username, resource, rights},
	}
	_, err := management.Enforcer.AddPoliciesEx(rules)
	return err
}

package authorization

import (
	"fmt"
	"log"
	"net/http"
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
	rights := strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut}, util.PIPE)
	rules := [][]string{
		{util.ROLES[3], "/users/*", rights},
		{util.ROLES[3], "/projects/*", rights},
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
func (management *FakeCasbinManagement) RemoveUserPolicies(username string) error {
	affected, err := management.Enforcer.DeleteUser(username)
	if !affected {
		log.Fatalf("%s not present in policies", username)
	}
	return err
}
func (management *FakeCasbinManagement) RemoveProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d", projectId)
	management.RemovePolicies(resource, client)
	management.RemovePolicies(resource, responsible)
	resource = fmt.Sprintf("/projects/%d/*", projectId)
	management.RemovePolicies(resource, client)
	management.RemovePolicies(resource, responsible)
	return nil
}
func (management *FakeCasbinManagement) CreateUserPolicies(username string, role string) error {
	var err error
	if role != util.ROLES[3] {
		resource := fmt.Sprintf("/users/%s", username)
		err = management.AddPolicies(resource, username, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	} else {
		err = management.CreateAdminPolicies(username)
	}
	return err
}
func (management *FakeCasbinManagement) CreateProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d", projectId)
	management.AddPolicies(resource, client, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	management.AddPolicies(resource, responsible, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	resource = fmt.Sprintf("/projects/%d/*", projectId)
	management.AddPolicies(resource, client, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut))
	management.AddPolicies(resource, responsible, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut))
	return nil
}
func (management *FakeCasbinManagement) RemovePolicies(resource string, username string) error {
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
func (management *FakeCasbinManagement) AddPolicies(resource string, username string, rights string) error {
	rules := [][]string{
		{"p", username, resource, rights},
	}
	_, err := management.Enforcer.AddPoliciesEx(rules)
	return err
}

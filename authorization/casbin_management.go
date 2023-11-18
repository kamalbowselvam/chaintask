package authorization

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/util"
	"go.uber.org/zap"
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
		{adminName, "/users*", rights},
		{adminName, "/projects*", rights},
	}

	_, err := management.Enforcer.AddPoliciesEx(rules)
	if err != nil {
		logger.Warn("", zap.Error(err))
	}
	return err
}
func (management *CasbinManagement) RemoveUserPolicies(username string) error {
	affected, err := management.Enforcer.DeleteUser(username)
	if !affected {
		logger.Warnf("%s not present in policies", username)
	}
	return err
}
func (management *CasbinManagement) RemoveProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d*", projectId)
	return errors.Join(management.RemovePolicies(resource, client), management.RemovePolicies(resource, responsible))
}
func (management *CasbinManagement) CreateUserPolicies(username string, role string) error {
	var err error
	if role != util.ROLES[3] {
		resource := fmt.Sprintf("/users/%s", username)
		err = management.AddPolicies(resource, username, util.GenerateRoleString(http.MethodGet, http.MethodPut))
	} else {
		err = management.CreateAdminPolicies(username)
	}
	return err
}
func (management *CasbinManagement) CreateProjectPolicies(projectId int64, client string, responsible string) error {
	resource := fmt.Sprintf("/projects/%d/tasks/", projectId)
	return errors.Join(management.AddPolicies(resource, client, util.GenerateRoleString(http.MethodGet, http.MethodPost)),
	management.AddPolicies(resource, responsible, util.GenerateRoleString(http.MethodGet, http.MethodPost)))

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
		logger.Warnf("problem while removing policies %s due to %s", policiesToRemove, err)
	}
	return err
}
func (management *CasbinManagement) AddPolicies(resource string, username string, rights string) error {
	rules := [][]string{
		{username, resource, rights},
	}
	_, err := management.Enforcer.AddPoliciesEx(rules)
	if err != nil {
		logger.Warnf("could not create policies %s", err)
	}
	return err
}
func (management *CasbinManagement) CreateTaskPolicies(taskId int64, projectId int64, author string) error {
	resource := fmt.Sprintf("/projects/%d/tasks/%d", projectId, taskId)
	return management.AddPolicies(resource, author, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete))
}
func (management *CasbinManagement) RemoveTaskPolicies(taskId int64, projectId int64, author string) error {
	resource := fmt.Sprintf("/projects/%d/tasks/%d", projectId, taskId)
	return management.RemovePolicies(resource, author)
}

func (management *CasbinManagement) CreateCompanyPolicies(companyId int64, user string) error {
	resource := fmt.Sprintf("/company/%d", companyId)
	return management.AddPolicies(resource, user, util.GenerateRoleString(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete))
}

func (management *CasbinManagement) RemoveCompanyPolicies(companyId int64, user string) error {
	resource := fmt.Sprintf("/company/%d", companyId)
	return management.RemovePolicies(resource, user)
}

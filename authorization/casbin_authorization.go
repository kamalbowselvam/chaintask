package authorization

import (
	"log"
	"os"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/kamalbowselvam/chaintask/util"
)

type CasbinAuthorization struct {
	Adapter  *pgadapter.Adapter
	Enforcer casbin.Enforcer
}

func NewCasbinAuthorization(source string, conf string) (AuthorizationService, error) {
	adapter, err := pgadapter.NewAdapter(source)
	if err != nil {
		panic(err)
	}
	// Load model configuration file and policy store adapter
	conf_file_path := "./config/rbac_model.conf"
	_, err = os.Stat(conf_file_path)
	if err != nil {
		conf_file_path = "." + conf_file_path
	}
	enforcer, err := casbin.NewEnforcer(conf_file_path, adapter)
	if err != nil {
		log.Fatal(err)
	}
	authorize := &CasbinAuthorization{
		Adapter:  adapter,
		Enforcer: *enforcer,
	}

	return authorize, nil
}
func (authorize *CasbinAuthorization) CreateEnforcer()
func (authorize *CasbinAuthorization) CreateAdapter()
func (authorize *CasbinAuthorization) LoadAdminPolicies() {
	rules := [][]string{
		// FIXME
		[]string{"p", util.ROLES[3], util.TASK, "*", util.READ},
		[]string{"p", util.ROLES[3], util.PROJECT, "*", util.WRITE},
		[]string{"p", util.ROLES[3], util.USER, "*", util.DELETE},
	}

	_, err := authorize.Enforcer.AddPoliciesEx(rules)
	if err != nil {
		log.Fatal(err)
	}
	_, err = authorize.Enforcer.AddGroupingPoliciesEx([][]string{{"g", util.ROLES[3], "kselvamADMIN"}})
	if(err != nil){
		log.Fatal(err)
	}
}
func (authorize *CasbinAuthorization) AddPolicy()
func (authorize *CasbinAuthorization) AddPolicies()
func (authorize *CasbinAuthorization) RemovePolicy()
func (authorize *CasbinAuthorization) RemovePolicies()
func (authorize *CasbinAuthorization) Enforce()

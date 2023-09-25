package authorization

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/stretchr/testify/require"
)

func TestEnforcing(t *testing.T) {

	policies := `p, kselvamADMIN, /users*, GET|POST|DELETE|PUT
p, kselvamADMIN, /projects*, GET|POST|DELETE|PUT
p, toto1,       /projects/1*, GET|POST|DELETE|PUT
`
	file, err := ioutil.TempFile("", "policies.*.csv")
	defer os.RemoveAll(file.Name())
	file.WriteString(policies)
	require.NoError(t, err)
	enforcer, err := casbin.NewEnforcer("../config/rbac_model.conf", file.Name())
	require.NoError(t, err)
	require.NotEmpty(t, enforcer)
	enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch)
	enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.RegexMatch)
	enforcer.EnableLog(true)
	cond, err := enforcer.Enforce("kselvamADMIN", "/projects/1", "GET")
	require.NoError(t, err)
	require.True(t, cond)
}

package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
}

func TestMainTestSuite(testContext *testing.T) {
	verbose = true
	suite.Run(testContext, new(MainTestSuite))
}

func (suite *MainTestSuite) SetupTest() {
	fakeOnline() // uses keys_test fakeserver
	config = *new(configuration)
	viper.Reset()
	configure()
}

// TODO test config not found error
// TODO test flag

func setTestPermissionsPermissive() {
	gitPerms := []string{"asherhawk"}
	allPerms := []string{"duncanblack"}
	config.Permissions = map[string][]string{"all": allPerms, "git": gitPerms}
}

func (suite *MainTestSuite) TestAuthorizedKeysOf() {
	assert := assert.New(suite.T())
	// before setting explicit perms, verify the bootstrap keys is returned
	autoPerms := authorizedKeysOf("root")
	assert.Contains(autoPerms, config.BootstrapKey)

	// use explicit perms to verify the permissions
	setTestPermissionsPermissive()
	gitAuthorizedKeys := authorizedKeysOf("git")
	assert.Contains(gitAuthorizedKeys, "Asher1")
	assert.Contains(gitAuthorizedKeys, "Asher2")
	assert.Contains(gitAuthorizedKeys, "Duncan1")
	assert.NotContains(gitAuthorizedKeys, config.BootstrapKey)
	
}

// ensure the special permissions of `all` work when present
// and that when missing randmo users get 0 access
func (suite *MainTestSuite) TestWildcardPermissions() {
	assert := assert.New(suite.T())
	
	// verify that the bootstrap key isn't present with explicit perms
	rootPerms := []string{"asherhawk"}
	config.Permissions = map[string][]string{"root": rootPerms}

	// first test the specific case of the allowed users
	rootAuthorizedKeys := authorizedKeysOf("root")
	assert.NotContains(rootAuthorizedKeys, config.BootstrapKey)
	assert.Contains(rootAuthorizedKeys, "Asher1")
	assert.Contains(rootAuthorizedKeys, "Asher2")

	// check for other [random] users and make sure there are 0 keys
	assert.Equal(len(permittedAccountsFor("abc123")), 0)
	assert.Empty(authorizedKeysOf("abc123"))

	assert.Equal(len(permittedAccountsFor("nobody")), 0)
	assert.Empty(authorizedKeysOf("nobody"))

	assert.Equal(len(permittedAccountsFor("git")), 0)
	assert.Empty(authorizedKeysOf("git"))

}

func (suite *MainTestSuite) TestPermittedAccountsFor() {
	assert := assert.New(suite.T())
	// without any config
	for _, user := range []string{"root", "git", "something"} {
		permitted := permittedAccountsFor(user)
		assert.Equal(len(permitted), 0)
	}

	// with specific perms
	setTestPermissionsPermissive()

	gitPermitted := permittedAccountsFor("git")
	assert.Equal(len(gitPermitted), 2)
	gitPermittedList := strings.Join(gitPermitted, " ")
	assert.Contains(gitPermittedList, "asherhawk")
	assert.Contains(gitPermittedList, "duncanblack")

	rootPermitted := permittedAccountsFor("root")
	assert.Equal(len(rootPermitted), 1)

	whateverPermitted := permittedAccountsFor("whatever")
	assert.Equal(len(whateverPermitted), 1)
	assert.Equal(whateverPermitted[0], "duncanblack")

}

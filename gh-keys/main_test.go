package main

import (
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/spf13/viper"
	"testing"
	"strings"
)

type MainTestSuite struct {
	suite.Suite
}

func TestMainTestSuite(testContext *testing.T){
	suite.Run(testContext, new(MainTestSuite))
}

func (suite *MainTestSuite) SetupTest() {
	config = *new(configuration)
	viper.Reset()
	configure()
}

func (suite *MainTestSuite) TestPermittedAccountsFor() {
	assert := assert.New(suite.T())
	// without any config
	for _, user := range []string{"root","git","something"} {
		permitted := permittedAccountsFor(user)
		assert.Contains(permitted[0],config.BootstrapKey)	
		assert.Equal(len(permitted),1)	
	}

	// with specific perms
	gitPerms := []string{"asherhawk"}
	allPerms := []string{"duncanblack"}
	config.Permissions = map[string][]string{"all":allPerms, "git": gitPerms}

	gitPermitted := permittedAccountsFor("git")
	assert.Equal(len(gitPermitted),2)
	gitPermittedList := strings.Join(gitPermitted, " ")

	assert.Contains(gitPermittedList, "asherhawk")
	assert.Contains(gitPermittedList, "duncanblack")

	rootPermitted := permittedAccountsFor("root")
	assert.Equal(len(rootPermitted),1)

	whateverPermitted :=  permittedAccountsFor("whatever")
	assert.Equal(len(whateverPermitted),1)
	assert.Equal(whateverPermitted[0],"duncanblack")

}

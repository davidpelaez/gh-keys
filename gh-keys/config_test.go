package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(testContext *testing.T) {
	verbose = true
	suite.Run(testContext, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) SetupTest() {
	config = *new(configuration)
	viper.Reset()
	configure()
}

// TODO test bad config error

func (suite *ConfigTestSuite) TestDefaults() {
	assert := assert.New(suite.T())

	config = *new(configuration)
	configure() // make sure configuration is alway reset

	assert.Equal(len(config.Permissions), 0) // by default no perms
	assert.True(config.AllowPanicMode)
	assert.Equal(config.TTL, 300)
	assert.Empty(config.BootstrapKeyFile)
	assert.Empty(config.ConfigFile)
	assert.Contains(config.BootstrapKey, builtinPublicKey)
	pwd, error := os.Getwd()
	check(error)
	assert.Contains(config.KeysDir, pwd)
}

func exampleFilePath(filename string) string {
	configPath, error := filepath.Abs("./config/" + filename)
	check(error)
	return configPath
}

func (suite *ConfigTestSuite) TestConfigFileParsing() {
	assert := assert.New(suite.T())
	config.ConfigFile = exampleFilePath("config.yaml")
	configure()
	assert.Equal(len(config.Permissions), 2)
	assert.Equal(len(config.Permissions["git"]), 2)
	assert.Equal(len(config.Permissions["all"]), 1)
	assert.Equal(config.Permissions["all"][0], "asherhawk")
	assert.False(config.AllowPanicMode)
	assert.Equal(config.TTL, 10)
	assert.False(config.AllowPanicMode)
	assert.Contains(config.KeysDir, "test/config/keys")
	// TODO check alternative bootstrap key
}

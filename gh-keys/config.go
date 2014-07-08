package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type configuration struct {
	AllowPanicMode   bool
	TTL              int
	BootstrapKeyFile string
	APIToken         string
	Permissions      map[string][]string
	BootstrapKey     string
	ConfigFile       string
	KeysDir          string
	InternetTestURL  string
}

const builtinPublicKey string = `ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ== vagrant insecure public key`

var config configuration

func configure() {

	// use config file if set from the flags
	if config.ConfigFile != "" {
		debugPrint("Using config file " + config.ConfigFile)
		viper.SetConfigFile(config.ConfigFile)
	}

	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/ghk")
	viper.AddConfigPath("$HOME/.ghk")
	viper.AddConfigPath(".")

	viper.SetDefault("AllowPanicMode", true)
	viper.SetDefault("TTL", 1)
	viper.SetDefault("InternetTestURL", "http://icanhazip.com/")

	err := viper.ReadInConfig()
	switch err.(type) {
	case viper.UnsupportedConfigError:
		debugPrint("No config file, using defaults")
	default:
		check(err)
	}

	config.ConfigFile = viper.ConfigFileUsed()
	if config.ConfigFile != "" {
		debugPrint("Resolving KeysDir relative to config file")
		// KeysDir defined relative to the config file
		keysDir := filepath.Dir(config.ConfigFile) + "/keys"
		keysDir, err = filepath.Abs(keysDir)
		check(err)
		config.KeysDir = keysDir
	} else {
		debugPrint("Resolving KeysDir relative to binary")
		// KeysDir is relative to the working dir
		binaryPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		config.KeysDir = binaryPath + "/keys"
	}

	err = viper.GetAllIntoStruct(&config)
	check(err)

	if config.BootstrapKeyFile != "" {
		debugPrint("Using config defined bootstrap key file")
		dat, err := ioutil.ReadFile(config.BootstrapKeyFile)
		check(err)
		config.BootstrapKey = strings.TrimSpace(string(dat))
	} else {
		config.BootstrapKey = builtinPublicKey
	}

	// ensure keys dir exists
	debugPrint("Creating " + config.KeysDir)
	mkdirError := os.MkdirAll(config.KeysDir, 0700)
	if mkdirError != nil {
		debugPrint("Error trying to create " + config.KeysDir)
		check(mkdirError)
	}
}

func logConfigItem(key string, rawValue interface{}) {
	fmt.Printf("%19s %v\n", key+":", rawValue)
}

func printConfigSummary() {
	logConfigItem("ConfigFile", config.ConfigFile)
	logConfigItem("KeysDir", config.KeysDir)
	logConfigItem("AllowPanicMode", config.AllowPanicMode)
	logConfigItem("TTL", config.TTL)
	logConfigItem("BootstrapKeyFile", config.BootstrapKeyFile)
	logConfigItem("APIToken", config.APIToken)

	perms := ""
	for username, accounts := range config.Permissions {
		perms = perms + fmt.Sprintf("%s: %v ", username, accounts)
	}

	logConfigItem("Permissions", perms)
}

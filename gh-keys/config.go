package main

import (	
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"fmt"
	"path/filepath"
	"os"
)

type configuration struct {
	AllowPanicMode bool
	TTL int
	BootstrapKeyFile string
	APIToken string
	Permissions map[string][]string
	BootstrapKey string
	ConfigFile string
	KeysDir string
}

const BuiltinKeyFilename string = "vagrant.pub"

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
	
	err := viper.ReadInConfig()
	switch err.(type){
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
	}else{
		debugPrint("Resolving KeysDir relative to binary")
		// KeysDir is relative to the working dir
		binaryPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		config.KeysDir = binaryPath + "/keys"
	}


	err = viper.GetAllIntoStruct(&config); check(err)

	if config.BootstrapKeyFile != "" {
		debugPrint("Using config defined bootstrap key file")
		dat, err := ioutil.ReadFile(config.BootstrapKeyFile); check(err)
		config.BootstrapKey = strings.TrimSpace(string(dat))
	} else {
		config.BootstrapKeyFile = BuiltinKeyFilename + " (builtin)"
		keyBytes, err := Asset(BuiltinKeyFilename); check(err)
		config.BootstrapKey = strings.TrimSpace(string(keyBytes))
	}

}

func logConfigItem(key string, rawValue interface{}) {
	fmt.Printf("%20s %v\n", key + ":" , rawValue)
}

func printConfigSummary() {
	logConfigItem("ConfigFile", config.ConfigFile)
	logConfigItem("KeysDir", config.KeysDir)
	logConfigItem("AllowPanicMode", config.AllowPanicMode)
	logConfigItem("TTL", config.TTL)
	logConfigItem("BootstrapKeyFile", config.BootstrapKeyFile)
	logConfigItem("APIToken", config.APIToken)
	logConfigItem("Permissions", config.Permissions)
}

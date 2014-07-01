package main

import (	
	"github.com/spf13/viper"
	"io/ioutil"
)

type Configuration struct {
	AllowPanicMode bool
	TTL int
	BootstrapKeyFile string
	BootstrapKey string
	APIToken string
	Permissions map[string]string
}

func configure() {
	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/ghk")
	viper.AddConfigPath("$HOME/.ghk")
	viper.AddConfigPath(".")

	viper.SetDefault("allow_panic_mode", true)
	viper.SetDefault("ttl", 300)
	
	viper.ReadInConfig()

	err := viper.GetAllIntoStruct(&config); check(err)

	// TODO trim key
	if config.BootstrapKey != "" {
		verbosePrint("Using config defined bootstrap key file")
		dat, err := ioutil.ReadFile(config.BootstrapKeyFile); check(err)
		config.BootstrapKey = string(dat)
	} else {
		keyBytes, err := Asset("vagrant.pub"); check(err)
		config.BootstrapKey = string(keyBytes)
	}
}

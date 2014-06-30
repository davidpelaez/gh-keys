package main


import (
  //  "os"
    "fmt"
    "io/ioutil"
    "github.com/spf13/viper"
)


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func setDefaultConfiguration() {
    viper.SetConfigName("config")

    viper.AddConfigPath("/etc/ghk")
    viper.AddConfigPath("$HOME/.ghk")
    viper.AddConfigPath(".")

    viper.SetDefault("allow_panic_mode", true)
    viper.SetDefault("ttl", "300") // 5 minutes
}

func permittedAccountsFor(username string) []string {
    // determine who can access this username
    return []string{}
}

func syncKeys(accounts []string) {
    // get the key in memory. If no connection and no panic on bad download returns an error for panic
    // confirm format, then write to file
}

func keysOf(account string) string {
    // send to stdout the key IF still valid
    // otherwise verify if we're in panic mode, aka no internet
    // and then if true print
    return ""
}

func printConfigSummary() {
    
    fmt.Println("config_file:", viper.ConfigFileUsed())
    
    fmt.Println("allow_panic_mode: ", viper.GetBool("allow_panic_mode"))
    
    fmt.Println("allow: ", viper.GetStringMap("allow"))
    fmt.Println("allow: ", viper.GetStringMapString("allow"))
    
    fmt.Println("ttl: ", viper.GetString("ttl"))
}

func printBootstrapKey() {
    if viper.IsSet("bootstrap_key_file"){
        keyPath := viper.GetString("bootstrap_key_file")
        dat, err := ioutil.ReadFile(keyPath)
        check(err)
        fmt.Print(string(dat))
    }else{
        keyBytes, err := Asset("vagrant.pub")
        check(err)
        vagrantKey := string(keyBytes)
        fmt.Println(vagrantKey)
    }
}

func printUsage() {
    fmt.Println("Usage...")
}
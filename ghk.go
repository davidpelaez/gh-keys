package main

import (
    "os"
    "fmt"
    "github.com/spf13/viper"
)

func main() {

    if len(os.Args) < 2 {
        printUsage()
    }

    // todo check if it's command
    

    setDefaultConfiguration()

    viper.ReadInConfig()

    username := os.Args[1]
    fmt.Println("Getting auth keys for", username)
    // todo append to the log

    printConfigSummary()

    if viper.IsSet("allow"){
        // permittedAccountsFor
        // syncKeys, if error check for panic mode, then continue, or exit
        // print the keys of the applicable users, calling printKeysOf
    }else{
        printBootstrapKey()
    }

}

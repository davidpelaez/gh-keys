package main

import (
	"fmt"
	"os"
	"flag"
)

// globally accessible configuration
var config = new(Configuration)
var verbose bool

func check(e error) {
	if e != nil {
		//failWith(e)
		failWith("pending...")
	}
}

func failWith(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

func verbosePrint(msg string){
	if verbose {
		fmt.Println(msg)
	}
}

func printUsage() {
	fmt.Println("Usage...")
}

func authorize(username string) {

	verbosePrint("Getting auth keys for" + username)
	fmt.Println(keysOf(username))
	os.Exit(128)
	// todo append to the log

	// permittedAccountsFor
	// syncKeys, if error check for panic mode, then continue, or exit
	// print the keys of the applicable users, calling printKeysOf
}

func logConfigItem(key string, rawValue interface{}) {
	printableValue, ok := rawValue.(string)
	if !ok {
		failWith("couldnt convert " + key)
	}

	//if printableValue != "" {
		fmt.Println(key + ": " , printableValue)
	//	}else{
	//		fmt.Println(key, "none")
	//	}
}

func printConfigSummary() {
	logConfigItem("AllowPanicMode", config.AllowPanicMode)
	logConfigItem("TTL", config.TTL)
	logConfigItem("BootstrapKeyFile", config.BootstrapKeyFile)
	logConfigItem("BootstrapKey", config.BootstrapKey)
	logConfigItem("APIToken", config.APIToken)
	logConfigItem("Permissions", config.Permissions)
}

func main() {

	info := flag.Bool("i", false, "Display configuration summary")
	flag.BoolVar(&verbose,"v", false, "Verbose mode")
	flag.StringVar(&config.APIToken, "t", "", "Github API token")
	
	flag.Parse()
	configure()
	arguments := flag.Args()

	switch {
	case *info:
		printConfigSummary()
	case len(arguments) == 1:
		authorize(arguments[0])
	default:
		printUsage()
	}
}

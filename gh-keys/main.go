package main

import (
	"fmt"
	"os"
	"flag"
)

var verbose bool

func check(e error) {
	if e != nil {
		failWith(e.Error())
	}
}

func failWith(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

func debugPrint(msg string){
	if verbose {
		fmt.Println("[debug]",msg)
	}
}

func printUsage() {
	fmt.Println("Usage...")
}

func main() {

	info := flag.Bool("i", false, "Display configuration summary")

	flag.BoolVar(&verbose,"v", false, "Verbose mode")
	flag.StringVar(&config.APIToken, "t", "", "Github API token")
	flag.StringVar(&config.ConfigFile, "c", "", "Config file to use")
	
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


func authorize(username string) {

	debugPrint("Getting auth keys for" + username)
	fmt.Println(keysOf(username))
	os.Exit(128)
	// todo append to the log

	// permittedAccountsFor
	// syncKeys, if error check for panic mode, then continue, or exit
	// print the keys of the applicable users, calling printKeysOf
}

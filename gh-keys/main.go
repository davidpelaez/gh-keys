package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var verbose bool
var panicMode = false

func check(e error) {
	if e != nil {
		failWith(e.Error())
	}
}

func failWith(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func debugPrint(msg string) {
	if verbose {
		// output aligned with go test verbose format
		fmt.Println("    [debug]", msg)
	}
}

func printUsage() {
	fmt.Println("Usage...")
}

func main() {

	info := flag.Bool("i", false, "Display configuration summary")

	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&config.APIToken, "t", "", "Github API token")
	flag.StringVar(&config.ConfigFile, "c", "", "Config file to use")

	flag.Parse()
	configure()
	arguments := flag.Args()

	switch {
	case *info:
		printConfigSummary()
	case len(arguments) == 1:
		fmt.Println(authorizedKeysOf(arguments[0]))
	default:
		printUsage()
	}
}

func permittedAccountsFor(username string) []string {
	permittedAccounts := make([]string, 0)

	if admins, ok := config.Permissions["all"]; ok {
		permittedAccounts = append(permittedAccounts, admins...)
	}

	if userSpecific, ok := config.Permissions[username]; ok {
		permittedAccounts = append(permittedAccounts, userSpecific...)
	}

	if len(permittedAccounts) == 0 {
		return []string{config.BootstrapKey}
	}

	return permittedAccounts
}

func authorizedKeysOf(username string) string {
	debugPrint("testing connectivity")
	panicMode = !online()
	if panicMode && !config.AllowPanicMode {
		debugPrint("Internet connectivity failed and panic mode isn't allowed")
		os.Exit(1)
	}
	authorizedKeys := make([]string, 0)
	debugPrint("getting keys")
	for _, permittedAccount := range permittedAccountsFor(username) {
		authorizedKeys = append(authorizedKeys, printableKeysOf(permittedAccount))
	}
	// TODO trim empty lines
	emptyLines := regexp.MustCompile("^$")
	keys := strings.Join(authorizedKeys, "\n")
	return emptyLines.ReplaceAllString(keys, "")
}

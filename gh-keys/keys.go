package main

import (
	//"github.com/gorilla/http"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// define outside a function to be reused in the pkg
var client = &http.Client{Timeout: 2 * time.Second}
var githubAPI = "https://api.github.com/users/"

type PublicKey struct {
	Id  int
	Key string
}

func getKeysOf(account string) []string {

	url := githubAPI + account + "/keys"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if config.APIToken != "" {
		// use personal token from the config if available
		debugPrint("Using auth token")
		req.SetBasicAuth(config.APIToken, "x-oauth-basic ")
	}

	resp, err := client.Do(req)

	if err != nil {
		// handle protocol error
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic("Request wan't 200!")
	}

	// verify response code
	defer resp.Body.Close()
	keysBody, err := ioutil.ReadAll(resp.Body)

	keysCollection := make([]PublicKey, 0)
	json.Unmarshal(keysBody, &keysCollection)

	keys := make([]string, 0)
	for _, key := range keysCollection {
		keys = append(keys, key.Key)
	}

	return keys
}

func syncKeys(accounts []string) {
	// get the key in memory. If no connection and no panic on bad download returns an error for panic
	// confirm format, then write to file
}

func keysOf(account string) string {
	// send to stdout the key IF still valid
	// otherwise verify if we're in panic mode, aka no internet
	// and then if true print
	keys := getKeysOf(account)
	for _, k := range keys {
		debugPrint("  - " + k)
	}

	return "" //string(body)
}

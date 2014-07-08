package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var timeout = 5 * time.Second
var client = &http.Client{Timeout: timeout}
var githubAPI = "https://api.github.com/users/"

type PublicKey struct {
	Id  int
	Key string
}

func online() bool {
	req, _ := http.NewRequest("GET", config.InternetTestURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		// a protocol error is a if there were no internet
		debugPrint("Proto error in online check with " + config.InternetTestURL)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		debugPrint("Unexpected code in online check: " + resp.Status)
		return false
	}

}

func getAPIKeysOf(account string) ([]string, error) {

	url := githubAPI + account + "/keys"

	offlineError := errors.New("Offline API")

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// use personal token from the config if available
	if config.APIToken != "" {
		req.SetBasicAuth(config.APIToken, "x-oauth-basic")
	}

	resp, err := client.Do(req)
	if err != nil {
		// a protocol error is a if the API is offline
		return nil, offlineError
	}
	defer resp.Body.Close()

	keys := make([]string, 0)

	switch resp.StatusCode {
	case 200:
		break // continue to JSON parsing
	case 404:
		debugPrint("Github user " + account + " wasn't found, ignoring")
		return keys, nil // empty keys array, no such user
	default:
		// any other code is treated as if the API is offline
		return nil, offlineError
	}

	keysBody, err := ioutil.ReadAll(resp.Body)

	keysCollection := make([]PublicKey, 0)
	json.Unmarshal(keysBody, &keysCollection)

	for _, key := range keysCollection {
		keys = append(keys, key.Key)
	}

	return keys, nil
}

func keyFilepath(account string) string {
	filename := account + ".pub"
	return path.Join(config.KeysDir, filename)
}

func readKeyFileOf(account string) (string, bool) {
	filepath := keyFilepath(account)
	if _, err := os.Stat(filepath); err != nil && os.IsNotExist(err) {
		debugPrint("No key file found for " + account)
		return "", false // no stored keys
	} else {
		info, error := os.Stat(filepath)
		check(error)
		now := time.Now().Unix()
		modTime := info.ModTime().Unix()
		expirationTime := modTime + int64(config.TTL)

		isValid := false
		if expirationTime > now {
			debugPrint("Key file of " + account + " is within TTL")
			isValid = true
		} else {
			ago := strconv.FormatInt(now-expirationTime, 2)
			debugPrint("Key file of " + account + "expired " + ago + "s ago")
		}

		keys, error := ioutil.ReadFile(filepath)
		check(error)
		return string(keys), isValid
	}
}

func deleteKeyFile(account string) {
	filepath := keyFilepath(account)
	if _, err := os.Stat(filepath); err != nil && !os.IsNotExist(err) {
		debugPrint("Deleting " + filepath)
		removeError := os.Remove(filepath)
		check(removeError)
	}
}

func printableKeysOf(account string) string {
	cachedKeys, stillValid := readKeyFileOf(account)
	if stillValid || panicMode {
		if panicMode {
			debugPrint("Using expired key(s) of " + account + " in panic mode")
		}
		return cachedKeys
	} else {
		keys, error := getAPIKeysOf(account)
		if error != nil {
			time.Sleep(2 * time.Second) // retry once after sleep
			keys, error = getAPIKeysOf(account)
			check(error)
		}
		filepath := keyFilepath(account)
		content := strings.Join(keys, "\n")
		// don't write empty file that could end up being cached
		if len(keys) != 0 {
			deleteKeyFile(account)
			// TODO verify format before writing to file
			error = ioutil.WriteFile(filepath, []byte(content), 0600)
			check(error)
		}
		return content
	}
}

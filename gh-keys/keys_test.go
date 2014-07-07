package main

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/go-martini/martini"
	"net/http"
	"net/http/httptest"
	"testing"
	//"net/http_test"
)

var fakeServer httptest.Server

func prepareFakeServer() {

	asherHawkKeys := `[ { "id": 1, "key": "ssh-rsa Asher1..." }, 
		{ "id": 2, "key": "ssh-rsa Asher2..." } ]`

	duncanBlackKeys := `[ { "id": 3, "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ1111....." }]`

	fakeKeysHandler := func (params martini.Params, r *http.Request) (int, string) {
		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			return 500, "Missing API version header"
		}
		switch params["username"] {
		case "asherhawk":
			return 200, asherHawkKeys
		case "duncanblack":
			return 200, duncanBlackKeys
		default:
			return 404, "Not found"
		}
	}

	m := martini.Classic()
	m.Get("/users/:username/keys", fakeKeysHandler)
	fakeServer = *httptest.NewServer(http.HandlerFunc(m.ServeHTTP))
	githubAPI = fakeServer.URL + "/users/"
}

// set fake http server
// modify the keys API url

func TestSuccessfulKeyRetrieval(t *testing.T) {
	// with a successful request, verify keys are parsed and given in the right array
	// verify they format of each key and length
	prepareFakeServer()
	keys := getKeysOf("asherhawk")

	assert := assert.New(t)
	assert.Equal(len(keys), 2, "returned keys don't match the expected length")
	
	for _, key := range keys {
		// we use this to detect any error in the fake server's asnwers
		assert.Contains(key, "ssh-rsa", "returned keys don't match the expected format")
	}
}

func TestKeySyncing(t *testing.T) {
	// test expiration of files in disk
	// writing to the right location
	assert.True(t, true, "pending...")
}

func TestBadAccountRetrieval(t *testing.T) {
	// verify what happens when an account isn't found
	assert.True(t, true, "pending...")
}

func TestOfflineSync(t *testing.T) {
	// verify errors are returned when the api is offline
	assert.True(t, true, "pending...")
}

func TestPermittedKeys(t *testing.T) {
	// verify the format of all keys
	// verify the length
	assert.True(t, true, "pending...")
}

func TestAPIToken(t *testing.T) {
	// ensure the token is used when provided in the config
	assert.True(t, true, "pending...")
}

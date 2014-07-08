package main

import (
	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	//"time"
	"io/ioutil"
)

var receivedAuthorization = "--"
var badURL = "http://localhost:1111/users/" // made up url = offline
var fakeServer httptest.Server

type KeysTestSuite struct {
	suite.Suite
}

func (suite *KeysTestSuite) SetupTest() {
	fakeOnline()
}

func TestKeysTestSuite(testContext *testing.T) {
	verbose = true
	configure() // make sure configuration defaults are set
	fakeServer = buildFakeServer()
	suite.Run(testContext, new(KeysTestSuite))
}

func buildFakeServer() httptest.Server {
	asherHawkKeys := `[ { "id": 1, "key": "ssh-rsa Asher1..." }, 
		{ "id": 2, "key": "ssh-rsa Asher2..." } ]`

	duncanBlackKeys := `[ { "id": 3, "key": "ssh-rsa Duncan1" }]`

	fakeKeysHandler := func(params martini.Params, r *http.Request) (int, string) {
		debugPrint("Fakeserver received req. for account " + params["username"])
		receivedAuthorization = r.Header.Get("Authorization")
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

	router := martini.NewRouter()
	raw_martini := martini.New()
	//raw_martini.Use(martini.Logger())
	//	raw_martini.Use(martini.Recovery())
	raw_martini.MapTo(router, (*martini.Routes)(nil))
	raw_martini.Action(router.Handle)

	m := martini.ClassicMartini{raw_martini, router}

	m.Get("/users/:username/keys", fakeKeysHandler)
	server := *httptest.NewServer(http.HandlerFunc(m.ServeHTTP))
	githubAPI = server.URL + "/users/"
	return server
}

func fakeOffline() {
	githubAPI = badURL
}

func fakeOnline() {
	githubAPI = fakeServer.URL + "/users/"
}

func (suite *KeysTestSuite) TestSuccessfulRetrieval() {

	assert := assert.New(suite.T())

	keys, _ := getAPIKeysOf("duncanblack")
	assert.Equal(len(keys), 1)

	keys, _ = getAPIKeysOf("asherhawk")
	assert.Equal(len(keys), 2)

	for _, key := range keys {
		// we use this to detect any error in the fake server's asnwers
		assert.Contains(key, "ssh-rsa", "returned keys don't match the expected format")
	}
}

func (suite *KeysTestSuite) TestBadAccountRetrieval() {
	keys, _ := getAPIKeysOf("johndoe")
	assert.Equal(suite.T(), len(keys), 0)
}

func (suite *KeysTestSuite) TestOfflineRetrieval() {
	fakeOffline()
	_, error := getAPIKeysOf("duncanblack")
	assert.NotNil(suite.T(), error)
}

// ensure the token is used when provided in the config
func (suite *KeysTestSuite) TestAPIToken() {
	config.APIToken = "aloha123"
	tokenBase64 := "YWxvaGExMjM6eC1vYXV0aC1iYXNpYw==" //encoded version
	getAPIKeysOf("tokentester")
	assert.Contains(suite.T(), receivedAuthorization, tokenBase64)
}

func createKeyFile(account string) {
	keyfile := keyFilepath(account)
	deleteKeyFile(account)
	debugPrint("Writing " + keyfile)
	content := []byte(config.BootstrapKey)
	error := ioutil.WriteFile(keyfile, content, 0600)
	check(error)
}

func (suite *KeysTestSuite) TestReadKeyFileOf() {
	assert := assert.New(suite.T())
	createKeyFile("johndoe")

	config.TTL = -1 // make file expiration be in the past
	_, valid := readKeyFileOf("johndoe")
	assert.False(valid)

	createKeyFile("janedoe")

	config.TTL = 3600 // make file expiration be in the future
	keys, valid := readKeyFileOf("janedoe")
	assert.True(valid)
	assert.Equal(builtinPublicKey, keys)
}

func purgeStoredKeys() {
	error := os.RemoveAll(config.KeysDir + "/*.pub")
	if error != nil && !os.IsNotExist(error) {
		debugPrint("Tried to remove " + config.KeysDir + " it exists but got error")
		check(error)
	}
	debugPrint("Deleted all keys in " + config.KeysDir)
}

func (suite *KeysTestSuite) TestPrintableKeysOf() {
	assert := assert.New(suite.T())

	createKeyFile("johndoe")

	// cached version
	config.TTL = 3600 // make file expiration be in the future
	assert.Equal(printableKeysOf("johndoe"), builtinPublicKey)

	// getting unknown user
	config.TTL = -1 // make file expiration be in the past
	// no such user in the API and key expired, hence should be empty
	assert.Equal(printableKeysOf("johndoe"), "")

	// getting valid user from scratch

	purgeStoredKeys()

	keyFile := keyFilepath("duncanblack")

	_, err := os.Stat(keyFile)
	assert.True(os.IsNotExist(err)) // ensure the file's absent

	assert.Contains(printableKeysOf("duncanblack"), "Duncan1")

	_, err = os.Stat(keyFilepath("duncanblack"))
	assert.Nil(err) // ensure the file was created

	// on panic mode
	fakeOffline()
	panicMode = true

	createKeyFile("johndoe")
	config.TTL = -1 // make file expiration be in the past
	assert.Equal(builtinPublicKey, printableKeysOf("johndoe"))

	// when a user has no local keyfile and api offline, empty keys returned
	assert.Equal("", printableKeysOf("unknown"))

}

func (suite *KeysTestSuite) TestOnlineCheck() {
	assert := assert.New(suite.T())
	config.InternetTestURL = badURL
	assert.False(online())
	config.InternetTestURL = "http://icanhazip.com"
	assert.True(online())
}

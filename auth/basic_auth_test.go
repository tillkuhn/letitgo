package auth

import (
	"github.com/tillkuhn/letitgo/shared/ioutil"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var successResponse = "you made it"
var errorResponse = "this call is a no go"

func TestBasicAuth(t *testing.T) {
	user := "hase"
	pw := "popase"

	wrappedProtectedFunc := BasicAuthMiddleware(protectedFunc, user, pw)
	testServer := httptest.NewServer(wrappedProtectedFunc)
	defer testServer.Close()

	// auth not set
	resAuthNotSet, err := http.Post(testServer.URL+"/secret", "text/plain", http.NoBody)
	assert.NoError(t, err)
	defer ioutil.SafeClose(resAuthNotSet.Body)
	body, err := io.ReadAll(resAuthNotSet.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resAuthNotSet.StatusCode)
	assert.Contains(t, string(body), "Unauthorized")

	// auth invalid credentials
	reqAuthWrong, err := http.NewRequest(http.MethodPost, testServer.URL+"/secret", http.NoBody)
	assert.NoError(t, err)
	reqAuthWrong.SetBasicAuth("alice", "*******")
	respAuthWrong, err := http.DefaultClient.Do(reqAuthWrong)
	assert.NoError(t, err)
	defer ioutil.SafeClose(respAuthWrong.Body)
	body, err = io.ReadAll(respAuthWrong.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, respAuthWrong.StatusCode)
	assert.Contains(t, string(body), "Unauthorized")

	// auth valid credentials
	reqAuthOK, err := http.NewRequest(http.MethodPost, testServer.URL+"/secret", http.NoBody)
	assert.NoError(t, err)
	reqAuthOK.SetBasicAuth(user, pw)
	respAuthOK, err := http.DefaultClient.Do(reqAuthOK)
	assert.NoError(t, err)
	defer ioutil.SafeClose(respAuthOK.Body)
	body, err = io.ReadAll(respAuthOK.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAuthOK.StatusCode)
	assert.Contains(t, string(body), successResponse)

}

func protectedFunc(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "/secret") {
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(successResponse))
	} else {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte(errorResponse))
	}
}

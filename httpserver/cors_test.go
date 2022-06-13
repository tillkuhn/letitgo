package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/test", nil)
	assert.NoError(t, err)
	rw := httptest.NewRecorder()
	setupCORSHeader(rw, req)
	// assert.Contains(t, rw.Body.String(), "<html")
	assert.Equal(t, uiRoot, rw.Header().Get("Access-Control-Allow-Origin"))
}

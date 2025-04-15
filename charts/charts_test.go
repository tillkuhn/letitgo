package charts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/test", nil)
	assert.NoError(t, err)
	rw := httptest.NewRecorder()
	writeLines(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)
}

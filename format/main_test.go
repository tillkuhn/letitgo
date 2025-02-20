package format

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenExpiry(t *testing.T) {
	claims := make(map[string]interface{})
	now := time.Now().Unix()

	claims[JwtClaimExpiry] = float64(now + 3600 + 60) // more than 1h
	expiryInfo, err := HumanizeTokenExpiry(claims)
	assert.Equal(t, "1 hour from now ⚠️", expiryInfo)
	assert.NoError(t, err)

	claims[JwtClaimExpiry] = float64(now + 604800 + 60) // more than a week
	expiryInfo, err = HumanizeTokenExpiry(claims)
	assert.Equal(t, "1 week from now 👍", expiryInfo)
	assert.NoError(t, err)

	delete(claims, JwtClaimExpiry)
	expiryInfo, err = HumanizeTokenExpiry(claims)
	assert.Equal(t, "", expiryInfo)
	assert.Error(t, err)

}

func TestTemplateString(t *testing.T) {
	assert.Contains(t, TemplateString(), "*****")
}

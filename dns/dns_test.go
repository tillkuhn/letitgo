package dns

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNS(t *testing.T) {
	ips, err := net.LookupIP("google.com")
	assert.NoError(t, err)
	assert.True(t, len(ips) > 0)
	for _, ip := range ips {
		t.Logf("google.com. IN A %s\n", ip.String())
	}
}

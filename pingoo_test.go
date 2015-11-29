// Run via `sudo go test -cover`
package pingoo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err1 := New("http://whatsthis")

	assert.Error(t, err1, "Should return an error on bad address")

	p2, err2 := New("127.0.0.1")

	assert.NoError(t, err2, "Should not get an error with good address")
	assert.NotNil(t, p2.FastPing, "FastPing should not be nil after good instantiation")
	assert.Equal(t, "127.0.0.1", p2.IPv4Address.String(), "IP address should match")
}

func TestPing(t *testing.T) {
	rtt := time.Duration(100) * time.Millisecond

	p1, pErr1 := New("127.0.0.1")
	assert.NoError(t, pErr1, "Unexpected error: %v", pErr1)

	r1, rErr1 := p1.Ping(5, rtt)
	assert.NoError(t, rErr1, "Should not get error with good ping: %v", rErr1)
	assert.Equal(t, 5, r1)

	p2, pErr2 := New("127.0.0.254")
	assert.NoError(t, pErr2, "Unexpected error: %v", pErr2)

	r2, rErr2 := p2.Ping(5, rtt)
	assert.NoError(t, rErr2, "Should not get an error pinging an IP that does not ping")
	assert.Equal(t, 0, r2)
}

package whatsup

import (
	"net"
	"testing"
)

func TestAddrError(t *testing.T) {
	// Open up listener on (hopefully) an unused port
	_, err := net.Listen("tcp", ":12345")
	if err != nil {
		t.Errorf("Unknown error, perhaps 12345 is already in use: %s", err)
	}

	// Re open listener on previously opened port, this should cause an EADDRINUSE error
	_, err = net.Listen("tcp", ":12345")

	// Test that our EADDRINUSE detector works
	if !addrInUse(err) {
		t.Errorf("Should correctly catch EADDRINUSE error!")
	}
}

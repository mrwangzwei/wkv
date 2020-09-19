package dns

import "testing"

func TestNewServer(t *testing.T) {
	a := NewServer(9901)
	a.Listen()
}

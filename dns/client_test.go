package dns

import (
	"testing"
)

func TestSend(t *testing.T) {
	Send("127.0.0.1:9991", "aaaaaaaaa")
}

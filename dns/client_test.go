package dns

import "testing"

func TestSend(t *testing.T) {
	Send("114.114.114.114:53", "www.baidu.com")
}

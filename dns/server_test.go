package dns

import (
	"fmt"
	"testing"
)

//sync.Map{} 提供并发锁安全锁接口

func TestNewServer(t *testing.T) {
	s, err := NewServer(9991, 10, ClientMode)
	fmt.Println(err)
	err = s.Listen()
	fmt.Println(err)
}

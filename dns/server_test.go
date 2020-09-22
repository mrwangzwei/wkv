package dns

import (
	"fmt"
	"testing"
)

//sync.Map{} 提供并发锁安全锁接口

func TestNewServer(t *testing.T) {
	s, err := NewServer(9991, 30, WeightMode)
	fmt.Println(err)
	err = s.AddWeightIpInfo("wzw", "127.0.0.1", 1)
	err = s.AddWeightIpInfo("wzw", "127.0.0.2", 2)
	err = s.AddWeightIpInfo("wzw", "127.0.0.3", 3)
	fmt.Println(table)
	err = s.Listen()
	fmt.Println(err)
}

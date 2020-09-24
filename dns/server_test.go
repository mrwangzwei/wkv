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

func TestWeightServer(t *testing.T) {
	s, err := NewServer(9991, 10, WeightMode)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = s.AddWeightIpInfo("wzw", "127.0.0.1", 1); err != nil {
		fmt.Println(err)
		return
	}
	if err = s.AddWeightIpInfo("wzw", "127.0.0.2", 3); err != nil {
		fmt.Println(err)
		return
	}
	if err = s.AddWeightIpInfo("wzw", "127.0.0.3", 4); err != nil {
		fmt.Println(err)
		return
	}
	err = s.Listen()
	fmt.Println(err)
}

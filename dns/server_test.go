package dns

import (
	"fmt"
	"testing"
)

//sync.Map{} 提供并发锁安全锁接口

func TestNewServer(t *testing.T) {
	fmt.Println(aaa(1))
}

func bb() {
	defer func() {
		fmt.Println("aaaa")
	}()

	defer func() {
		fmt.Println("bbbb")
	}()
}

func aaa(n int) (r int) {
	defer func() {
		r += n
		fmt.Println("r += n", r, n)
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var f func()

	defer f()

	f = func() {
		r += 2
	}

	return n + 1
}

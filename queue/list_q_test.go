package queue

import (
	"fmt"
	"runtime"
	"testing"
)

func TestNewListQueue(t *testing.T) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	fmt.Println(m)

	q, err := NewListQueue(1000000)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 1000000; i++ {
		if err := q.Push([]byte("aaaaaaaaaaaaaaaaaa")); err != nil {
			fmt.Println(err)
		}
	}
	runtime.ReadMemStats(m)
	fmt.Println(m)

	for j := 0; j < 1000000; j++ {
		if _, err := q.Pull(); err != nil {
			fmt.Println(err)
		}
	}
	runtime.ReadMemStats(m)
	fmt.Println(m)

	for i := 0; i < 10000; i++ {
		if err := q.Push([]byte("aaaaaaaaaaaaaaaaaa")); err != nil {
			fmt.Println(err)
		}
	}

	//time.Sleep(121 * time.Second)
	runtime.ReadMemStats(m)
	fmt.Println(m)

	cont, err := q.Pull()
	fmt.Println(string(cont), err)
}

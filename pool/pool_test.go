package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p, err := NewPool(3)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 20; i++ {
		err = p.AddTask(p.NewTask(testFunA, i, errHandle))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second) //加这个测一下协程数是不是递增的
	}
	p.Wait() //等所有task执行完
}

func testFunA(a interface{}) error {
	if a == 3 {
		panic("33333333333333333333333")
	}
	fmt.Println("testFunAAAAAAAAAA", a)
	return nil
}

func errHandle(params, err interface{}) {
	fmt.Println("errHandle", params, err)
}

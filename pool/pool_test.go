package pool

import (
	"fmt"
	"sync"
	"testing"
)

func Test_Pool(t *testing.T) {
	withPool()
}

func Test_NoPool(t *testing.T) {
	noPool()
}

func noPool() {
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func(i interface{}) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			_ = testFunA(i)
		}(i)
	}
	wg.Wait()
}

func withPool() {
	p, err := NewPool(20)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 2000; i++ {
		p.AddTask(p.NewTask(testFunA, i, errHandle))
	}
	p.Wait() //等所有task执行完
}

func testFunA(a interface{}) error {
	if a.(int) == 3 {
		panic("33333333333333333333333")
	}
	fmt.Println("testFunAAAAAAAAAA", a)
	return nil
}

func errHandle(params interface{}, err string) {
	fmt.Println("errHandle", params, err)
}

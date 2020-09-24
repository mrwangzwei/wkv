package main

import (
	"fmt"
	"wkv/dns"
)

func main() {

	dns.SendJson("127.0.0.1:9991", "wzw")

	////测试recover
	//aaa()

	////优雅退出测试
	//sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//go func() {
	//	sig := <-sigs
	//	fmt.Println()
	//	fmt.Println(sig)
	//	done <- true
	//}()
	//fmt.Println("awaiting signal")
	//<-done
	//fmt.Println("exiting")

	////WaitGroup测试
	//wait := sync.WaitGroup{}
	//for i := 0; i < 10; i++ {
	//	wait.Add(1)
	//	go func(wait *sync.WaitGroup, i int) {
	//		//defer wait.Done()
	//		fmt.Println(i)
	//	}(&wait, i)
	//}
	//wait.Wait()

	////channel测试
	//ch := make(chan int)
	//defer close(ch)
	//for i := 0; i < 20; i++ {
	//	go func(i int) {
	//		ch <- i
	//	}(i)
	//}
	//
	//go func(ch chan int) {
	//	for i := range ch {
	//		fmt.Println(i)
	//	}
	//}(ch)
	//time.Sleep(time.Duration(1) * time.Second)
}

func aaa() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(1111)
			fmt.Println(err)
		}
	}()
	panic("aaaaaaaaa")
}

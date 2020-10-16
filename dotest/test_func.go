package dotest

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ChannelUse() {
	sc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sc
		fmt.Println("退出信号", sig)
		fmt.Println("退出前做点什么...")
		if sig.String() == "interrupt" {
			done <- true
		}
	}()
	fmt.Println("正常业务代码部分...")
	<-done
	fmt.Println("退出")
}

func Aaaaa() {
	var ichan = make(chan int, 1)
	var str string

	go func() {
		str = "hello world"
		time.Sleep(2 * time.Second)
		ichan <- 0
	}()
	aa := <-ichan //这里有值，下面才会运行

	fmt.Println(str, aa)
}

func Bbbb() {
	defer func() {
		v := recover()
		fmt.Println(v)
	}()
	go func() {
		panic("aaa")
	}()
	time.Sleep(3 * time.Second)
}

func ListenInput() {
	for {
		aa := bufio.NewReader(os.Stdin)
		str, _, err := aa.ReadLine()
		fmt.Println("bufio.NewReader", string(str), err)

		var strr string
		cc := bufio.NewScanner(os.Stdin)
		if cc.Scan() {
			strr = cc.Text()
		} else {
			strr = "Find input error"
		}
		fmt.Println("bufio.NewScanner", strr)

		var bb []byte
		_, err = fmt.Scan(&bb)
		fmt.Println("fmt.Scan", string(bb), err)
	}
}

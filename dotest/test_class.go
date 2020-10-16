package dotest

import "fmt"

type ClassA struct {
}

func (c ClassA) Aaa() {
	fmt.Println("ClassA Aaa")
}

type ClassB struct {
	ClassA
}

func (c ClassB) Aaa() {
	fmt.Println("ClassB Aaa")
}

func (c ClassB) Bbb() {
	fmt.Println("ClassB Bbb")
}

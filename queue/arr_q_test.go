package queue

import (
	"fmt"
	"testing"
)

func TestNewArrQueue(t *testing.T) {
	q := NewArrQueue(3)
	err := q.Push([]byte("a"))
	fmt.Println(err)

	err = q.Push([]byte("b"))
	fmt.Println(err)

	err = q.Push([]byte("c"))
	fmt.Println(err)

	d, err := q.Pull()
	fmt.Println(string(d), err)

	d, err = q.RPull()
	fmt.Println(string(d), err)

	d, err = q.RPull()
	fmt.Println(string(d), err)

}

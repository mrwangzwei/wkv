package wkv_a

import (
	"fmt"
	"testing"
)

func TestAaaa(t *testing.T) {
	Kvt = &Boot{}
	err := Kvt.Set("aaaaa", "111111")
	err = Kvt.Set("bbbbb", "222222")
	err = Kvt.Set("ccccc", "333333")
	err = Kvt.Set("ddddd", "444444")
	err = Kvt.Set("eeeee", "555555")
	err = Kvt.Set("fffff", "666666")
	err = Kvt.Set("ggggg", "7777777")
	err = Kvt.Set("hhhhh", "888888")
	fmt.Println(err)
	ccc, err := Kvt.Get("ccccc")
	fmt.Println(ccc)
	ddd, err := Kvt.Get("ddddd")
	fmt.Println(ddd)
	eee, err := Kvt.Get("eeeee")
	fmt.Println(eee)
	err = Kvt.Del("ddddd")
	ccc, err = Kvt.Get("ccccc")
	fmt.Println(ccc)
	ddd, err = Kvt.Get("ddddd")
	fmt.Println(err)
	eee, err = Kvt.Get("eeeee")
	fmt.Println(eee)



}

package bar_graph

import (
	"fmt"
	"wkv/vars"
)

type BarTable struct {
	data []int
	x, y int
}

func NewBarTable(arr []int) *BarTable {
	t := &BarTable{
		data: arr,
	}
	t.init()
	return t
}

func (t *BarTable) init() {
	t.x = len(t.data)
	t.y = t.data[0]
	for _, item := range t.data {
		if item > t.y {
			t.y = item
		}
	}
}

func (t *BarTable) SwapAndDraw(i, j int) {
	draw(t.data, t.x, t.y, false, i, j)
	if i < t.x && j < t.y {
		t.data[i], t.data[j] = t.data[j], t.data[i]
	}
	draw(t.data, t.x, t.y, false, i, j)
	draw(t.data, t.x, t.y, false)
}

func (t *BarTable) Draw(pre bool) {
	draw(t.data, t.x, t.y, pre)
}

func draw[T vars.Integer](arr []T, x, y T, pre bool, lightIndex ...T) {
	lightMap := make(map[T]struct{})
	for _, item := range lightIndex {
		lightMap[item] = struct{}{}
	}
	fmt.Println("")
	var i, j T
	for i = 0; i < y; i++ {
		str := ""
		for j = 0; j < x; j++ {
			p := " █ "
			if pre {
				p = " □ "
			}
			if _, ok := lightMap[j]; ok {
				p = " | "
			}
			if arr[int(j)]+i < y {
				str += "   "
			} else {
				str += p
			}
		}
		fmt.Println(str)
	}
	fmt.Println("")
}

func DrawBarGraph[T vars.Integer](arr []T, pre bool, light ...T) {
	x := T(len(arr))
	y := arr[0]
	for _, item := range arr {
		if item > y {
			y = item
		}
	}
	draw(arr, x, y, pre, light...)
}

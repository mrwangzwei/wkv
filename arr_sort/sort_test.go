package arr_sort

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test1aa(t *testing.T) {
	arr := randArr(10000)
	start := time.Now().UnixNano()
	QuickSortTwo(arr)
	end := time.Now().UnixNano()
	fmt.Println(end - start)

	arr2 := randArr(10000)
	start = time.Now().UnixNano()
	QuickSortPhp(arr2)
	end = time.Now().UnixNano()
	fmt.Println(end - start)

}

func randArr(len int) []int {
	var arr []int
	for len > 0 {
		len--
		arr = append(arr, rand.Intn(100))
	}
	return arr
}

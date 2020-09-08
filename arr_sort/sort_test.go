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
	QuickSort(arr)
	end := time.Now().UnixNano()
	fmt.Println(end - start)

	arr2 := randArr(10000)
	start = time.Now().UnixNano()
	QuickSortPhp(arr2)
	end = time.Now().UnixNano()
	fmt.Println(end - start)

}

func randArr(len int) []int {
	arr := make([]int, len)
	for len > 0 {
		len--
		arr = append(arr, rand.Intn(100))
	}
	return arr
}

func TestQuickSortDesc(t *testing.T) {
	aaa := []int{2, 3, 5, 1, 7, 2, 324, 678, 2, 54, 1567, 7893, 23, 7, 4}
	QuickSortDesc(aaa)
	fmt.Println(aaa)
}

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
	fmt.Println("QuickSort Nano", end-start)

	arr2 := randArr(10000)
	start = time.Now().UnixNano()
	QuickSortPhp(arr2)
	end = time.Now().UnixNano()
	fmt.Println("QuickSortPhp Nano", end-start)

}

func randArr(len int) []int {
	arr := make([]int, len)
	for len > 0 {
		len--
		arr = append(arr, rand.Intn(100))
	}
	return arr
}

func TestQuickSort(t *testing.T) {
	sli := []int{1, 2, 6, 9, 6, 7, 8, 9, 10, 2, 3, 4, 5}
	QuickSort(sli)
	fmt.Println(sli)
}

func TestQuickSortSlice(t *testing.T) {
	sli := []int{1, 2, 6, 9, 6, 7, 8, 9, 10, 2, 3, 4, 5}
	quickSort(sli, func(i, j int) bool {
		return sli[i] < sli[j]
	})
	fmt.Println(sli)
}

func TestSortBubbleSlice(t *testing.T) {
	sli := []int{1, 2, 6, 9, 6, 7, 8, 9, 10, 2, 3, 4, 5}
	SortBubbleSlice(sli, func(i, j int) bool {
		return sli[i] > sli[j]
	})
	fmt.Println(sli)
}

func TestSelectionSort(t *testing.T) {
	sli := []int{1, 2, 6, 9, 6, 7, 8, 9, 10, 2, 3, 4, 5}
	SelectionSort(sli, func(i, j int) bool {
		return sli[i] < sli[j]
	})
	fmt.Println(sli)
}

func TestInsertSort(t *testing.T) {
	sli := []int{1, 2, 6, 9, 6, 7, 8, 9, 10, 2, 3, 4, 5}
	InsertSort(sli, func(i, j int) bool {
		return sli[i] > sli[j]
	})
	fmt.Println(sli)
}

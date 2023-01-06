package arr_sort

import (
	"wkv/bar_graph"
	"wkv/vars"
)

func QuickSort[T vars.Integer](sli []T) {
	length := len(sli)
	if length <= 1 {
		return
	}
	i, mid := 1, sli[0]
	head, tail := 0, length-1
	for head < tail {
		if sli[i] > mid {
			sli[i], sli[tail] = sli[tail], sli[i]
			tail--
		} else {
			sli[i], sli[head] = sli[head], sli[i]
			head++
			i++
		}
	}
	sli[head] = mid
	QuickSort(sli[:head])
	QuickSort(sli[head+1:])
}

func QuickSortDesc[T vars.Integer](sli []T) {
	if len(sli) <= 1 {
		return
	}
	mid, i := sli[0], 1
	head, tail := 0, len(sli)-1
	for head < tail {
		if sli[i] > mid {
			sli[i], sli[head] = sli[head], sli[i]
			head++
			i++
		} else {
			sli[i], sli[tail] = sli[tail], sli[i]
			tail--
		}
	}
	sli[head] = mid
	QuickSortDesc(sli[:head])
	QuickSortDesc(sli[head+1:])
}

// php写法
func QuickSortPhp[T vars.Integer](arr []T) []T {
	length := len(arr)
	if length <= 1 {
		return arr
	}
	mid := arr[0]
	var left []T
	var right []T
	for i := 1; i < length; i++ {
		if arr[i] < mid {
			left = append(left, arr[i])
		} else {
			right = append(right, arr[i])
		}
	}
	left = QuickSortPhp(left)
	right = QuickSortPhp(right)
	return append(append(left, mid), right...)
}

func SortBubbleSlice[T vars.Integer](sli []T, less func(i, j int) bool) {
	length := len(sli)
	for i := 0; i < length-1; i++ {
		for j := i + 1; j < length; j++ {
			if !less(i, j) {
				bar_graph.DrawBarGraph(sli, true, T(i), T(j))
				sli[j], sli[i] = sli[i], sli[j]
				bar_graph.DrawBarGraph(sli, false, T(i), T(j))
			}
		}
	}
}

func SelectionSort[T vars.Integer](sli []T, less func(i, j int) bool) {
	length := len(sli)
	for i := 0; i < length-1; i++ {
		min := i
		for j := i + 1; j < length; j++ {
			if !less(min, j) {
				min = j
			}
		}
		sli[i], sli[min] = sli[min], sli[i]
	}
}

func InsertSort[T vars.Integer](sli []T, less func(i, j int) bool) {
	for i := range sli {
		for pre := i - 1; pre >= 0 && !less(pre, pre+1); pre-- {
			sli[pre], sli[pre+1] = sli[pre+1], sli[pre]
		}
	}
}

func quickSort(arr []int, less func(i, j int) bool) []int {
	return _quickSort(arr, 0, len(arr)-1, less)
}

func _quickSort(arr []int, left, right int, less func(i, j int) bool) []int {
	if left < right {
		partitionIndex := partition(arr, left, right, less)
		_quickSort(arr, left, partitionIndex-1, less)
		_quickSort(arr, partitionIndex+1, right, less)
	}
	return arr
}

func partition(arr []int, left, right int, less func(i, j int) bool) int {
	pivot := left
	index := pivot + 1

	for i := index; i <= right; i++ {
		if less(i, pivot) {
			swap(arr, i, index)
			index += 1
		}
	}
	swap(arr, pivot, index-1)
	return index - 1
}

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

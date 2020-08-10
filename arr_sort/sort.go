package arr_sort

func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	flag, leftFlag, rightFlag := arr[left], left, right
	for right != left {
		for right > left && arr[right] >= flag {
			right--
		}
		for left < right && arr[left] <= flag {
			left++
		}
		if right > left {
			arr[right], arr[left] = arr[left], arr[right]
		}
	}
	arr[right], arr[leftFlag] = flag, arr[right]
	quickSort(arr, leftFlag, left)
	quickSort(arr, right+1, rightFlag)
}

func QuickSortTwo(arr []int) {
	if len(arr) <= 1 {
		return
	}
	mid, i := arr[0], 1
	head, tail := 0, len(arr)-1
	for head < tail {
		if arr[i] > mid {
			arr[i], arr[tail] = arr[tail], arr[i]
			tail--
		} else {
			arr[i], arr[head] = arr[head], arr[i]
			head++
			i++
		}
	}
	arr[head] = mid
	QuickSortTwo(arr[:head])
	QuickSortTwo(arr[head+1:])
}


//php写法
func QuickSortPhp(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := arr[0]
	var left []int
	var right []int
	for i := 1; i < len(arr); i++ {
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

func SortBubble(arr []int) {
	if len(arr) <= 1 {
		return
	}
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[j], arr[j] = arr[i], arr[j]
			}
		}
	}
}

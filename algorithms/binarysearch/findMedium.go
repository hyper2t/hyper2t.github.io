package main

import (
	"fmt"
	"math"
)

func findMedium(nums []int) float64 {
	var left = 0
	var right = len(nums) - 1
	for left < right {
		var mid = left + int(math.Floor(float64(right-left)/2))
		if nums[mid] > nums[right] {
			// 中间数字大于右边数字，比如[3,4,5,1,2]，则左侧是有序上升的，最小值在右侧
			left = mid + 1
		} else if nums[mid] < nums[right] {
			// 中间数字小于右边数字，比如[6,7,1,2,3,4,5]，则右侧是有序上升的，最小值在左侧
			right = mid
		} else {
			// 中间数字等于右边数字，比如[2,3,1,1,1]或者[4,1,2,3,3,3,3]
			// 则重复数字可能为最小值，也可能最小值在重复值的左侧
			// 所以将right左移一位
			right -= 1
		}
	}

	// 寻找中位数
	var minIndex = left
	var midIndex = int(math.Floor(float64(len(nums)) / 2))
	if len(nums)%2 == 1 {
		if minIndex < midIndex {
			var medium = float64(nums[minIndex+midIndex])
			return medium
		} else {
			var medium = float64(nums[minIndex+midIndex-len(nums)])
			return medium
		}
	} else {
		if minIndex < midIndex {
			var medium = float64(nums[minIndex+midIndex]+nums[minIndex+midIndex-1]) / 2
			return medium
		} else {
			var medium = float64(nums[minIndex+midIndex-len(nums)]+nums[minIndex+midIndex-len(nums)-1]) / 2
			return medium
		}
	}
}

func main() {
	fmt.Println(findMedium([]int{2, 3, 4, 5, 1}))
	fmt.Println(findMedium([]int{6, 7, 1, 2, 3, 4, 5}))
	fmt.Println(findMedium([]int{2, 3, 4, 1}))
	fmt.Println(findMedium([]int{2, 3, 4, 5, 6, 1}))
}

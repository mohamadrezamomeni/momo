package utils

import "math"

func Min(nums ...int) int {
	mi := math.Inf(1)
	for _, num := range nums {
		if mi > float64(num) {
			mi = float64(num)
		}
	}
	return int(mi)
}

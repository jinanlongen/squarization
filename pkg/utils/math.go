// utils/helper.go
package utils

func Max(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	maxVal := vals[0]
	for _, val := range vals {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

func Min(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	minVal := vals[0]
	for _, v := range vals[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

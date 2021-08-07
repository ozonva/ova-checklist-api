package utils

import "fmt"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ToSliceOfChunks(values []int, size int) [][]int {
	if size <= 0 {
		return [][]int{}
	}

	capacity := len(values) / size
	if len(values)%size != 0 {
		capacity++
	}

	// Allocate enough space before slicing in order to avoid additional allocations
	result := make([][]int, 0, capacity)
	for begin := 0; begin < len(values); begin += size {
		end := min(begin+size, len(values))
		result = append(result, values[begin:end])
	}
	return result
}

func InvertMap(dict map[string]int) map[int]string {
	result := make(map[int]string)
	for k, v := range dict {
		if _, exists := result[v]; exists {
			panic(fmt.Sprintf("collision for key: %v", v))
		}
		result[v] = k
	}
	return result
}

var blacklist = []int{1, 2, 4, 8, 16, 32, 64, 128}

func FilterByBlacklist(values []int) []int {
	result := make([]int, 0)
	for _, value := range values {
		isAllowedValue := true
		for _, badValue := range blacklist {
			if value == badValue {
				isAllowedValue = false
				break
			}
		}
		if isAllowedValue {
			result = append(result, value)
		}
	}
	return result
}

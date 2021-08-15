package utils

import (
	"errors"
	"fmt"

	"ova-checklist-api/internal/types"
)

var (
	ErrUserIdCollision = errors.New("two or more user IDs are being used as map keys")
)

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
	result := make(map[int]string, len(dict))
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

func SplitToChunks(checklists []types.Checklist, size int) [][]types.Checklist {
	if size <= 0 {
		return [][]types.Checklist{}
	}

	capacity := len(checklists) / size
	if len(checklists)%size != 0 {
		capacity++
	}

	// Allocate enough space before slicing in order to avoid additional allocations
	result := make([][]types.Checklist, 0, capacity)
	for begin := 0; begin < len(checklists); begin += size {
		end := min(begin+size, len(checklists))
		result = append(result, checklists[begin:end])
	}
	return result
}

func MapChecklistsByUserId(checklists []types.Checklist) (map[uint64]types.Checklist, error) {
	result := make(map[uint64]types.Checklist, len(checklists))
	for _, checklist := range checklists {
		if _, exists := result[checklist.UserId]; exists {
			return nil, ErrUserIdCollision
		}
		result[checklist.UserId] = checklist
	}
	return result, nil
}

package utils

import (
	"errors"
	"fmt"

	"ova-checklist-api/internal/types"
)

var (
	ErrUserIdCollision = errors.New("two or more user IDs are being used as map keys")
)

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func ToSliceOfChunks(values []int, chunkSize uint) [][]int {
	if chunkSize == 0 {
		return [][]int{}
	}

	valuesSize := uint(len(values))
	capacity := valuesSize / chunkSize
	if valuesSize % chunkSize != 0 {
		capacity++
	}

	// Allocate enough space before slicing in order to avoid additional allocations
	result := make([][]int, 0, capacity)
	for begin := uint(0); begin < valuesSize; begin += chunkSize {
		end := min(begin + chunkSize, valuesSize)
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

func SplitToChunks(checklists []types.Checklist, chunkSize uint) [][]types.Checklist {
	if chunkSize == 0 {
		return [][]types.Checklist{}
	}

	checklistsSize := uint(len(checklists))
	capacity := checklistsSize / chunkSize
	if checklistsSize % chunkSize != 0 {
		capacity++
	}

	// Allocate enough space before slicing in order to avoid additional allocations
	result := make([][]types.Checklist, 0, capacity)
	for begin := uint(0); begin < checklistsSize; begin += chunkSize {
		end := min(begin + chunkSize, checklistsSize)
		result = append(result, checklists[begin:end])
	}
	return result
}

func MapChecklistsByUserId(checklists []types.Checklist) (map[uint64]types.Checklist, error) {
	result := make(map[uint64]types.Checklist, len(checklists))
	for _, checklist := range checklists {
		if _, exists := result[checklist.UserID]; exists {
			return nil, ErrUserIdCollision
		}
		result[checklist.UserID] = checklist
	}
	return result, nil
}

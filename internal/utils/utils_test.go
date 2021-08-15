package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"ova-checklist-api/internal/types"
)

func TestToSliceOfChunks(t *testing.T) {
	tests := []struct {
		values   []int
		size     int
		expected [][]int
	}{
		{values: []int{1, 2, 3, 4, 5, 6, 7},       size: 0,     expected: [][]int{}},
		{values: []int{1, 2, 3, 4, 5, 6, 7},       size: -1,    expected: [][]int{}},
		{values: nil,                              size: -1,    expected: [][]int{}},
		{values: []int{},                          size: 2,     expected: [][]int{}},
		{values: nil,                              size: 2,     expected: [][]int{}},
		{values: []int{},                          size: -1,    expected: [][]int{}},
		{values: []int{1},                         size: 2,     expected: [][]int{{1}}},
		{values: []int{1, 2},                      size: 3,     expected: [][]int{{1, 2}}},
		{values: []int{1},                         size: 1,     expected: [][]int{{1}}},
		{values: []int{1, 2},                      size: 1,     expected: [][]int{{1}, {2}}},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8},    size: 2,     expected: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}},
		{values: []int{1, 2, 3},                   size: 2,     expected: [][]int{{1, 2}, {3}}},
		{values: []int{1, 2, 3, 4},                size: 3,     expected: [][]int{{1, 2, 3}, {4}}},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8},    size: 3,     expected: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}}},
	}

	for testId, ctx := range tests {
		t.Run(fmt.Sprintf("TestToSliceOfChunks_%d", testId), func(t *testing.T) {
			assert.Equal(t, ctx.expected, ToSliceOfChunks(ctx.values, ctx.size))
		})
	}
}

func TestInvertMap(t *testing.T) {
	tests := []struct {
		in  map[string]int
		out map[int]string
	}{
		{in: map[string]int{},                          out: map[int]string{}},
		{in: nil,                                       out: map[int]string{}},
		{in: map[string]int{"hello": 1, "world": 2},    out: map[int]string{1: "hello", 2: "world"}},
	}

	for testId, ctx := range tests {
		t.Run(fmt.Sprintf("TestInvertMap_%d", testId), func(t *testing.T) {
			assert.Equal(t, ctx.out, InvertMap(ctx.in))
		})
	}
}

func TestInvertMap_WithCollisions(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("Collision is not detected, panic expected")
		}
	}()

	result := InvertMap(map[string]int{
		"hello": 1,
		"world": 1,
	})
	assert.Equal(t, result, map[int]string{
		1: "hello",
	})
}

func TestFilterByBlacklist(t *testing.T) {
	tests := []struct {
		in  []int
		out []int
	}{
		{in: []int{1, 2, 3, 4, 5, 8, 9},    out: []int{3, 5, 9}},
		{in: []int{1},                      out: []int{}},
		{in: []int{3},                      out: []int{3}},
		{in: []int{},                       out: []int{}},
		{in: nil,                           out: []int{}},
	}

	for testId, ctx := range tests {
		t.Run(fmt.Sprintf("TestFilterByBlacklist_%d", testId), func(t *testing.T) {
			assert.Equal(t, ctx.out, FilterByBlacklist(ctx.in))
		})
	}
}

func TestSplitToChunks(t *testing.T) {
	tests := []struct {
		checklists   []types.Checklist
		size         int
		expected     [][]types.Checklist
	}{
		{checklists: buildChecklistSlice(7),       size: 0,     expected: [][]types.Checklist{}},
		{checklists: buildChecklistSlice(7),       size: -1,    expected: [][]types.Checklist{}},
		{checklists: nil,                          size: -1,    expected: [][]types.Checklist{}},
		{checklists: []types.Checklist{},          size: 2,     expected: [][]types.Checklist{}},
		{checklists: nil,                          size: 2,     expected: [][]types.Checklist{}},
		{checklists: []types.Checklist{},          size: -1,    expected: [][]types.Checklist{}},
		{checklists: buildChecklistSlice(1),       size: 2,     expected: buildChecklistChunks(1, 1)},
		{checklists: buildChecklistSlice(2),       size: 3,     expected: buildChecklistChunks(2, 3)},
		{checklists: buildChecklistSlice(1),       size: 1,     expected: buildChecklistChunks(1, 1)},
		{checklists: buildChecklistSlice(2),       size: 1,     expected: buildChecklistChunks(2, 1)},
		{checklists: buildChecklistSlice(8),       size: 2,     expected: buildChecklistChunks(8, 2)},
		{checklists: buildChecklistSlice(3),       size: 2,     expected: buildChecklistChunks(3, 2)},
		{checklists: buildChecklistSlice(4),       size: 3,     expected: buildChecklistChunks(4, 3)},
		{checklists: buildChecklistSlice(8),       size: 3,     expected: buildChecklistChunks(8, 3)},
	}

	for testId, ctx := range tests {
		t.Run(fmt.Sprintf("TestSplitToChunks_%d", testId), func(t *testing.T) {
			assert.Equal(t, ctx.expected, SplitToChunks(ctx.checklists, ctx.size))
		})
	}
}

func TestMapChecklistsByUserId(t *testing.T) {
	tests := []struct {
		checklists    []types.Checklist
		expected      map[uint64]types.Checklist
		expectedError error
	}{
		{checklists: buildChecklistSlice(0), expected: buildChecklistMap(0), expectedError: nil},
		{checklists: buildChecklistSlice(1), expected: buildChecklistMap(1), expectedError: nil},
		{checklists: buildChecklistSlice(8), expected: buildChecklistMap(8), expectedError: nil},
		{
			checklists: []types.Checklist{buildChecklist(0), buildChecklist(0)},
			expected: nil,
			expectedError: ErrUserIdCollision,
		},
	}

	for testId, ctx := range tests {
		t.Run(fmt.Sprintf("TestMapChecklistsByUserId_%d", testId), func(t *testing.T) {
			actual, err := MapChecklistsByUserId(ctx.checklists)
			assert.Equal(t, ctx.expectedError, err)
			if err == nil {
				assert.Equal(t, ctx.expected, actual)
			}
		})
	}
}

func buildChecklist(userId uint64) types.Checklist {
	return types.Checklist{
		UserId:      userId,
		Title:       "Default checklist",
		Description: "Testing checklist utils",
		Items: []types.ChecklistItem{
			{"Step 1", false},
		},
	}
}

func buildChecklistSlice(size int) []types.Checklist {
	result := make([]types.Checklist, size)
	for i := range result {
		result[i] = buildChecklist(uint64(i))
	}
	return result
}

func buildChecklistMap(size int) map[uint64]types.Checklist {
	result := make(map[uint64]types.Checklist, size)
	for i := 0; i < size; i++ {
		result[uint64(i)] = buildChecklist(uint64(i))
	}
	return result
}

func buildChecklistChunks(checklistCount, chunkSize int) [][]types.Checklist {
	result := make([][]types.Checklist, 0)
	for {
		if checklistCount <= 0 {
			break
		}
		sliceSize := min(checklistCount, chunkSize)
		result = append(result, buildChecklistSlice(sliceSize))
		checklistCount -= sliceSize
	}

	var userId uint64 = 0
	for chunkId := range result {
		for checklistId := range result[chunkId] {
			result[chunkId][checklistId].UserId = userId
			userId++
		}
	}

	return result
}

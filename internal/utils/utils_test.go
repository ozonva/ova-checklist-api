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
		{
			checklists: []types.Checklist{checklist(0), checklist(1)},
			size: 0,
			expected: [][]types.Checklist{},
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(1)},
			size: -1,
			expected: [][]types.Checklist{},
		},
		{
			checklists: nil,
			size: -1,
			expected: [][]types.Checklist{},
		},
		{
			checklists: []types.Checklist{},
			size: 2,
			expected: [][]types.Checklist{},
		},
		{
			checklists: nil,
			size: 2,
			expected: [][]types.Checklist{},
		},
		{
			checklists: []types.Checklist{},
			size: -1,
			expected: [][]types.Checklist{},
		},
		{
			checklists: []types.Checklist{checklist(0)},
			size: 2,
			expected: [][]types.Checklist{
				{checklist(0)},
			},
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(1)},
			size: 3,
			expected: [][]types.Checklist{
				{checklist(0), checklist(1)},
			},
		},
		{
			checklists: []types.Checklist{checklist(0)},
			size: 1,
			expected: [][]types.Checklist{
				{checklist(0)},
			},
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(1)},
			size: 1,
			expected: [][]types.Checklist{
				{checklist(0)},
				{checklist(1)},
			},
		},
		{
			checklists: []types.Checklist{
				checklist(0), checklist(1),
				checklist(2), checklist(3),
				checklist(4), checklist(5),
				checklist(6), checklist(7),
			},
			size: 2,
			expected: [][]types.Checklist{
				{checklist(0), checklist(1)},
				{checklist(2), checklist(3)},
				{checklist(4), checklist(5)},
				{checklist(6), checklist(7)},
			},
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(1), checklist(2)},
			size: 2,
			expected: [][]types.Checklist{
				{checklist(0), checklist(1)},
				{checklist(2)},
			},
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(1), checklist(2), checklist(3)},
			size: 3,
			expected: [][]types.Checklist{
				{checklist(0), checklist(1), checklist(2)},
				{checklist(3)},
			},
		},
		{
			checklists: []types.Checklist{
				checklist(0), checklist(1),
				checklist(2), checklist(3),
				checklist(4), checklist(5),
				checklist(6), checklist(7),
			},
			size: 3,
			expected: [][]types.Checklist{
				{checklist(0), checklist(1), checklist(2)},
				{checklist(3), checklist(4), checklist(5)},
				{checklist(6), checklist(7)},
			},
		},
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
		{
			checklists: []types.Checklist{},
			expected: map[uint64]types.Checklist{},
			expectedError: nil,
		},
		{
			checklists: []types.Checklist{checklist(0)},
			expected: map[uint64]types.Checklist{
				0: checklist(0),
			},
			expectedError: nil,
		},
		{
			checklists: []types.Checklist{
				checklist(0), checklist(1),
				checklist(2), checklist(3),
				checklist(4), checklist(5),
				checklist(6), checklist(7),
			},
			expected: map[uint64]types.Checklist{
				0: checklist(0),
				1: checklist(1),
				2: checklist(2),
				3: checklist(3),
				4: checklist(4),
				5: checklist(5),
				6: checklist(6),
				7: checklist(7),
			},
			expectedError: nil,
		},
		{
			checklists: []types.Checklist{checklist(0), checklist(0)},
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

func checklist(userId uint64) types.Checklist {
	return types.Checklist{
		UserID:      userId,
		Title:       "Default checklist",
		Description: "Testing checklist utils",
		Items: []types.ChecklistItem{
			{"Step 1", false},
		},
	}
}

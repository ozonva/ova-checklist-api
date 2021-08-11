package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

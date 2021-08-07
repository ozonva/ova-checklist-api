package utils

import (
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, result, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%v != %v", result, expected)
	}
}

func TestToSliceOfChunks_NotPositiveSize(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7}
	assertEquals(t, ToSliceOfChunks(values, 0), [][]int{})
	assertEquals(t, ToSliceOfChunks(values, -1), [][]int{})
}

func TestToSliceOfChunks_EmptySlice(t *testing.T) {
	assertEquals(t, ToSliceOfChunks([]int{}, 2), [][]int{})
	assertEquals(t, ToSliceOfChunks(nil, 2), [][]int{})
	assertEquals(t, ToSliceOfChunks([]int{}, -1), [][]int{})
}

func TestToSliceOfChunks_LenLesserThanSize(t *testing.T) {
	result := ToSliceOfChunks([]int{1}, 2)
	assertEquals(t, result, [][]int{
		{1},
	})

	result = ToSliceOfChunks([]int{1, 2}, 3)
	assertEquals(t, result, [][]int{
		{1, 2},
	})
}

func TestToSliceOfChunks_LenDivisibleBySize(t *testing.T) {
	result := ToSliceOfChunks([]int{1}, 1)
	assertEquals(t, result, [][]int{
		{1},
	})

	result = ToSliceOfChunks([]int{1, 2}, 1)
	assertEquals(t, result, [][]int{
		{1}, {2},
	})

	result = ToSliceOfChunks([]int{1, 2, 3, 4, 5, 6, 7, 8}, 2)
	assertEquals(t, result, [][]int{
		{1, 2}, {3, 4}, {5, 6}, {7, 8},
	})
}

func TestToSliceOfChunks_LenNotDivisibleBySize(t *testing.T) {
	result := ToSliceOfChunks([]int{1, 2, 3}, 2)
	assertEquals(t, result, [][]int{
		{1, 2}, {3},
	})

	result = ToSliceOfChunks([]int{1, 2, 3, 4}, 3)
	assertEquals(t, result, [][]int{
		{1, 2, 3}, {4},
	})

	result = ToSliceOfChunks([]int{1, 2, 3, 4, 5, 6, 7, 8}, 3)
	assertEquals(t, result, [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8},
	})
}

func TestInvertMap_EmptyMap(t *testing.T) {
	assertEquals(t, InvertMap(map[string]int{}), map[int]string{})
	assertEquals(t, InvertMap(nil), map[int]string{})
}

func TestInvertMap_WithoutCollisions(t *testing.T) {
	result := InvertMap(map[string]int{
		"hello": 1,
		"world": 2,
	})
	assertEquals(t, result, map[int]string{
		1: "hello",
		2: "world",
	})
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
	assertEquals(t, result, map[int]string{
		1: "hello",
	})
}

func TestFilterByBlacklist(t *testing.T) {
	assertEquals(t, FilterByBlacklist([]int{1, 2, 3, 4, 5, 8, 9}), []int{3, 5, 9})
	assertEquals(t, FilterByBlacklist([]int{1}), []int{})
	assertEquals(t, FilterByBlacklist([]int{3}), []int{3})
	assertEquals(t, FilterByBlacklist([]int{}), []int{})
	assertEquals(t, FilterByBlacklist(nil), []int{})
}

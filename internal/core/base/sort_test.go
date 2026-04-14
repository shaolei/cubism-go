package base

import (
	"sort"
	"testing"
)

func TestGetSortedDrawableIndicesLogic(t *testing.T) {
	t.Parallel()

	// Test the sorting logic used in GetSortedDrawableIndices
	// by simulating the algorithm with different order arrays

	t.Run("already sorted", func(t *testing.T) {
		orders := []int32{0, 1, 2, 3}
		result := sortByOrders(orders)
		for i, idx := range result {
			if idx != i {
				t.Errorf("sortByOrders(%v)[%d] = %d, want %d", orders, i, idx, i)
			}
		}
	})

	t.Run("reverse sorted", func(t *testing.T) {
		orders := []int32{3, 2, 1, 0}
		result := sortByOrders(orders)
		expected := []int{3, 2, 1, 0}
		for i, idx := range result {
			if idx != expected[i] {
				t.Errorf("sortByOrders(%v)[%d] = %d, want %d", orders, i, idx, expected[i])
			}
		}
	})

	t.Run("duplicated orders", func(t *testing.T) {
		orders := []int32{5, 3, 5, 1}
		result := sortByOrders(orders)
		// Sorted by order: 1(idx3), 3(idx1), 5(idx0), 5(idx2)
		expected := []int{3, 1, 0, 2}
		for i, idx := range result {
			if idx != expected[i] {
				t.Errorf("sortByOrders(%v)[%d] = %d, want %d", orders, i, idx, expected[i])
			}
		}
	})

	t.Run("single element", func(t *testing.T) {
		orders := []int32{42}
		result := sortByOrders(orders)
		if len(result) != 1 || result[0] != 0 {
			t.Errorf("sortByOrders(%v) = %v, want [0]", orders, result)
		}
	})

	t.Run("empty", func(t *testing.T) {
		orders := []int32{}
		result := sortByOrders(orders)
		if len(result) != 0 {
			t.Errorf("sortByOrders(empty) should return empty, got %v", result)
		}
	})

	t.Run("negative orders", func(t *testing.T) {
		orders := []int32{-1, 0, -5, 2}
		result := sortByOrders(orders)
		// Sorted: -5(idx2), -1(idx0), 0(idx1), 2(idx3)
		expected := []int{2, 0, 1, 3}
		for i, idx := range result {
			if idx != expected[i] {
				t.Errorf("sortByOrders(%v)[%d] = %d, want %d", orders, i, idx, expected[i])
			}
		}
	})
}

// sortByOrders replicates the sorting logic in GetSortedDrawableIndices
func sortByOrders(orders []int32) []int {
	count := len(orders)
	type orderEntry struct {
		index int
		order int32
	}
	entries := make([]orderEntry, count)
	for i := 0; i < count; i++ {
		entries[i] = orderEntry{index: i, order: orders[i]}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].order < entries[j].order
	})

	rs := make([]int, count)
	for i, e := range entries {
		rs[i] = e.index
	}
	return rs
}

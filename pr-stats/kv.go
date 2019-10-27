package main

import (
	"sort"
)

// kv is used to sort maps with strings as keys and integers as values.
type kv struct {
	Key   string
	Value int
}

// orderMap take a map with strings as keys and integers as values and
// converts it into a sorted slice of kv instances. These kv instances will be
// order by value (highest values first) and then by key (lexicographical
// order).
func orderMap(m map[string]int) []kv {
	var orderedKV []kv
	for k, v := range m {
		orderedKV = append(orderedKV, kv{Key: k, Value: v})
	}
	// Golang uses the quick sort algorithm. O(log n) space complexity and
	// average O(n log n) time complexity (O(n^2) in worst case).
	sort.Slice(orderedKV, func(i, j int) bool {
		if orderedKV[i].Value > orderedKV[j].Value {
			return true
		}
		if orderedKV[i].Value < orderedKV[j].Value {
			return false
		}
		return orderedKV[i].Key < orderedKV[j].Key
	})

	return orderedKV
}

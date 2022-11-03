// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package lru

import (
	"fmt"
	"testing"
)

func TestBlobLru(t *testing.T) {
	lru := NewSizeConstraiedLRU(100)
	var want uint64
	// Add 11 items of 10 byte each. First item should be swapped out
	for i := 0; i < 11; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("value-%04d", i)
		lru.Add(k, v)
		want += uint64(len(v))
		if want > 100 {
			want = 100
		}
		if have := lru.size; have != want {
			t.Fatalf("size wrong, have %d want %d", have, want)
		}
	}
	// Zero:th should be evicted
	{
		k := fmt.Sprintf("key-%d", 0)
		if val := lru.Get([]byte(k)); val != nil {
			t.Fatalf("should be evicted: %v", k)
		}
	}
	// Elems 1-11 should be present
	for i := 1; i < 11; i++ {
		k := fmt.Sprintf("key-%d", i)
		want := fmt.Sprintf("value-%04d", i)
		have := lru.Get([]byte(k))
		if have == nil {
			t.Fatalf("missing key %v", k)
		}
		if string(have) != want {
			t.Fatalf("wrong value, have %v want %v", have, want)
		}
	}
}

// TestBlobLruOverflow tests what happens when inserting an element exceeding
// the max size
func TestBlobLruOverflow(t *testing.T) {
	lru := NewSizeConstraiedLRU(100)
	// Add 10 items of 10 byte each, filling the cache
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("value-%04d", i)
		lru.Add(k, v)
	}
	// Add one single large elem. We expect it to swap out all entries.
	{
		k := fmt.Sprintf("large-%d", 0)
		v := make([]byte, 200)
		lru.Add(k, string(v))
	}
	// Elems 0-9 should be missing
	for i := 1; i < 10; i++ {
		k := fmt.Sprintf("key-%d", i)
		if val := lru.Get([]byte(k)); val != nil {
			t.Fatalf("should be evicted: %v", k)
		}
	}
	// The size should be accurate
	if have, want := lru.size, uint64(200); have != want {
		t.Fatalf("size wrong, have %d want %d", have, want)
	}
	// Adding one small item should swap out the large one
	{
		i := 0
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("value-%04d", i)
		lru.Add(k, v)
		if have, want := lru.size, uint64(10); have != want {
			t.Fatalf("size wrong, have %d want %d", have, want)
		}
	}
}

// TestBlobLruSameItem tests what happens when inserting the same k/v multiple times.
func TestBlobLruSameItem(t *testing.T) {
	lru := NewSizeConstraiedLRU(100)
	// Add one 10 byte-item 10 times
	k := fmt.Sprintf("key-%d", 0)
	v := fmt.Sprintf("value-%04d", 0)
	for i := 0; i < 10; i++ {
		lru.Add(k, v)
	}
	// The size should be accurate
	if have, want := lru.size, uint64(10); have != want {
		t.Fatalf("size wrong, have %d want %d", have, want)
	}
}

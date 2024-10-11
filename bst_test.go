// Package bst
// A concurrent safe, lockless binary search tree
// Copyright (C) Alex Gaetano Padula
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package bst

import (
	"flag"
	"fmt"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	bst := New()
	if bst == nil {
		t.Fatal("bst is nil")
	}
}

func TestBST_Put(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))
	bst.Put([]byte("key44"), []byte("value"))
	bst.Put([]byte("key2"), []byte("value"))

	//bst.Print()
}

func TestBST_Get(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))
	bst.Put([]byte("key44"), []byte("value 44"))
	bst.Put([]byte("key44"), []byte("value 44 2"))

	key := bst.Get([]byte("key"))
	if key == nil {
		t.Fatal("key is nil")
	}

	if string(key.Values[0]) != "value" {
		t.Fatal("value is not equal to value")
	}

	key = bst.Get([]byte("key44"))
	if key == nil {
		t.Fatal("key is nil")
	}

	if string(key.Values[0]) != "value 44" {
		t.Fatal("value is not equal to value 44")
	}

	if string(key.Values[1]) != "value 44 2" {
		t.Fatal("value is not equal to value 44 2")
	}

}

func TestBST_Remove(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))
	bst.Put([]byte("key"), []byte("value 2"))
	bst.Put([]byte("key"), []byte("value 3"))

	bst.Remove([]byte("key"), []byte("value 2"))

	key := bst.Get([]byte("key"))

	if len(key.Values) != 2 {
		t.Fatal("key.Values length is not 2")
	}

	if string(key.Values[0]) != "value" {
		t.Fatal("value is not equal to value")
	}

	if string(key.Values[1]) != "value 3" {
		t.Fatal("value is not equal to value 3")
	}
}

func TestBST_Delete(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))
	bst.Put([]byte("key33"), []byte("value 2"))
	bst.Put([]byte("key2"), []byte("value 3"))
	bst.Put([]byte("key3"), []byte("value 4"))

	bst.Delete([]byte("key2"))

	key := bst.Get([]byte("key2"))

	if key != nil {
		t.Fatal("key is not nil")
	}

	checkForKeys := []string{"key", "key33", "key3"}

	for _, k := range checkForKeys {
		key := bst.Get([]byte(k))
		if key == nil {
			t.Fatal("key is nil")
		}

	}
}

func TestBST_Range(t *testing.T) {
	bst := New()

	for i := 0; i < 100; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.Range([]byte("key10"), []byte("key20"))

	if len(keys) != 11 {
		t.Fatal("keys length is not 11")
	}

	expected := []string{
		"key10",
		"key11",
		"key12",
		"key13",
		"key14",
		"key15",
		"key16",
		"key17",
		"key18",
		"key19",
		"key20",
	}

	for i, key := range keys {
		if string(key.K) != expected[i] {
			t.Fatalf("expected %s, got %s", expected[i], string(key.K))
		}
	}
}

func TestBST_GreaterThan(t *testing.T) {
	bst := New()

	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.GreaterThan([]byte("key05"))

	for _, key := range keys {
		fmt.Println(string(key.K))
	}

	expect := []string{"key06", "key07", "key08", "key09"}

	for i, key := range keys {
		if string(key.K) != expect[i] {
			t.Fatalf("expected %s, got %s", expect[i], string(key.K))
		}
	}
}

func TestBST_GreaterThanEq(t *testing.T) {
	bst := New()

	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.GreaterThanEq([]byte("key05"))

	for _, key := range keys {
		fmt.Println(string(key.K))
	}

	expect := []string{"key05", "key06", "key07", "key08", "key09"}

	for i, key := range keys {
		if string(key.K) != expect[i] {
			t.Fatalf("expected %s, got %s", expect[i], string(key.K))
		}
	}
}

func TestBST_LessThan(t *testing.T) {
	bst := New()

	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.LessThan([]byte("key05"))

	for _, key := range keys {
		fmt.Println(string(key.K))
	}

	expect := []string{"key00", "key01", "key02", "key03", "key04"}

	for i, key := range keys {
		if string(key.K) != expect[i] {
			t.Fatalf("expected %s, got %s", expect[i], string(key.K))
		}
	}
}

func TestBST_LessThanEq(t *testing.T) {
	bst := New()

	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.LessThanEq([]byte("key05"))

	for _, key := range keys {
		fmt.Println(string(key.K))
	}

	expect := []string{"key00", "key01", "key02", "key03", "key04", "key05"}

	for i, key := range keys {
		if string(key.K) != expect[i] {
			t.Fatalf("expected %s, got %s", expect[i], string(key.K))
		}
	}
}

func TestBST_NGet(t *testing.T) {
	bst := New()

	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	keys := bst.NGet([]byte("key05"))

	expect := []string{"key00", "key01", "key02", "key03", "key04", "key06", "key07", "key08", "key09"}

	for i, key := range keys {
		if string(key.K) != expect[i] {
			t.Fatalf("expected %s, got %s", expect[i], string(key.K))
		}
	}
}

func TestBST_ConcurrentPut(t *testing.T) {
	bst := New()
	var wg sync.WaitGroup
	numGoroutines := 10
	keysPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < keysPerGoroutine; j++ {
				key := fmt.Sprintf("key%02d-%d", j, goroutineID)
				value := fmt.Sprintf("value-%d", goroutineID)
				bst.Put([]byte(key), []byte(value))
			}
		}(i)
	}

	wg.Wait()

	// Verify that all keys were inserted
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < keysPerGoroutine; j++ {
			key := fmt.Sprintf("key%02d-%d", j, i)
			val := bst.Get([]byte(key))
			if val == nil {
				t.Fatalf("Expected key %s not found", key)
			}
		}
	}
}

func TestBST_ConcurrentGet(t *testing.T) {
	bst := New()
	var wg sync.WaitGroup
	numGoroutines := 10
	keysPerGoroutine := 10

	// Pre-fill the tree
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < keysPerGoroutine; j++ {
			key := fmt.Sprintf("key%02d-%d", j, i)
			bst.Put([]byte(key), []byte(fmt.Sprintf("value-%d", i)))
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < keysPerGoroutine; j++ {
				key := fmt.Sprintf("key%02d-%d", j, goroutineID)
				val := bst.Get([]byte(key))
				if val == nil {
					t.Fatalf("Expected key %s not found", key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestBST_ConcurrentDelete(t *testing.T) {
	bst := New()
	var wg sync.WaitGroup
	numGoroutines := 10
	keysPerGoroutine := 10

	// Pre-fill the tree
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < keysPerGoroutine; j++ {
			key := fmt.Sprintf("key%02d-%d", j, i)
			bst.Put([]byte(key), []byte(fmt.Sprintf("value-%d", i)))
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < keysPerGoroutine; j++ {
				key := fmt.Sprintf("key%02d-%d", j, goroutineID)
				bst.Delete([]byte(key))
				// Verify that the key has been deleted
				val := bst.Get([]byte(key))
				if val != nil {
					t.Fatalf("Expected key %s to be deleted", key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestBST_DuplicateKeys(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value1"))
	bst.Put([]byte("key"), []byte("value2"))

	key := bst.Get([]byte("key"))
	if len(key.Values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(key.Values))
	}
}

func TestBST_EmptyKeyValue(t *testing.T) {
	bst := New()
	bst.Put([]byte(""), []byte("value"))
	key := bst.Get([]byte(""))
	if key == nil {
		t.Fatal("Expected to find empty key")
	}
}

func TestBST_SpecialCharacters(t *testing.T) {
	bst := New()
	specialKeys := []string{"key@#", "key space", "key:colon", "key#hash"}
	for _, k := range specialKeys {
		bst.Put([]byte(k), []byte("value"))
	}

	for _, k := range specialKeys {
		key := bst.Get([]byte(k))
		if key == nil {
			t.Fatalf("Expected to find key: %s", k)
		}
	}
}

func TestBST_ConcurrentDeletes(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bst.Delete([]byte("key"))
		}()
	}

	wg.Wait()
	if bst.Get([]byte("key")) != nil {
		t.Fatal("Key should have been deleted")
	}
}

func TestBST_GetAfterDelete(t *testing.T) {
	bst := New()
	bst.Put([]byte("key"), []byte("value"))
	bst.Delete([]byte("key"))

	if bst.Get([]byte("key")) != nil {
		t.Fatal("Expected key to be deleted")
	}
}

func TestBST_RangeNoMatches(t *testing.T) {
	bst := New()
	bst.Put([]byte("key01"), []byte("value"))
	bst.Put([]byte("key02"), []byte("value"))

	keys := bst.Range([]byte("key03"), []byte("key04"))
	if len(keys) != 0 {
		t.Fatal("Expected no keys in range")
	}
}

func TestBST_AllKeysRemoved(t *testing.T) {
	bst := New()
	for i := 0; i < 10; i++ {
		bst.Put([]byte(fmt.Sprintf("key%02d", i)), []byte("value"))
	}

	for i := 0; i < 10; i++ {
		bst.Delete([]byte(fmt.Sprintf("key%02d", i)))
	}

	if bst.Get([]byte("key00")) != nil {
		t.Fatal("Expected all keys to be deleted")
	}
}

// Benchmarking

var (
	numGoroutines    = flag.Int("numGoroutines", 10, "Number of goroutines to use for the benchmark")
	keysPerGoroutine = flag.Int("keysPerGoroutine", 1_000_000, "Number of keys per goroutine")
)

// BenchmarkBST_ConcurrentPut benchmarks the concurrent Put operation.
func BenchmarkBST_ConcurrentPut(b *testing.B) {
	flag.Parse() // Parse command line flags

	b.Run(fmt.Sprintf("ConcurrentPut-%dGoroutines-%dKeys", *numGoroutines, *keysPerGoroutine), func(b *testing.B) {
		b.ResetTimer()
		var wg sync.WaitGroup
		bst := New()

		for i := 0; i < *numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < *keysPerGoroutine; j++ {
					key := []byte(fmt.Sprintf("key%09d-%d", j, goroutineID))
					value := []byte(fmt.Sprintf("value-%d", goroutineID))
					bst.Put(key, value)
				}
			}(i)
		}

		wg.Wait()
	})
}

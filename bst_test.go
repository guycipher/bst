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

import "testing"

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

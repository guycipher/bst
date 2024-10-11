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
	"unsafe"
)

// BST is the binary search tree struct
type BST struct {
	Root unsafe.Pointer // Root of the binary search tree
}

// Node is a node within the binary search tree
type Node struct {
	Key   *Key           // Key of the node
	Left  unsafe.Pointer // Left node
	Right unsafe.Pointer // Right node
}

// Key is the key for the binary search tree
type Key struct {
	K      []byte   // Key value
	Values [][]byte // Values within the key
}

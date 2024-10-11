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
	"bytes"
	"sync"
	"sync/atomic"
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
	Latch  *sync.Mutex
}

// New creates a new BST
func New() *BST {
	return &BST{}
}

// Put adds a new key to BST or append value to existing key
func (bst *BST) Put(key, value []byte) {
	newNode := &Node{Key: &Key{K: key, Values: [][]byte{value}, Latch: &sync.Mutex{}}}
	for {
		root := atomic.LoadPointer(&bst.Root)
		if root == nil {
			if atomic.CompareAndSwapPointer(&bst.Root, nil, unsafe.Pointer(newNode)) {
				return
			}
		} else {
			if bst.put(root, newNode) {
				return
			}
		}
	}
}

func (bst *BST) put(rootPointer unsafe.Pointer, newNode *Node) bool {
	root := (*Node)(rootPointer)

	if bytes.Compare(newNode.Key.K, root.Key.K) < 0 {
		left := atomic.LoadPointer(&root.Left)
		if left == nil {
			if atomic.CompareAndSwapPointer(&root.Left, nil, unsafe.Pointer(newNode)) {
				return true
			}
		} else {
			return bst.put(left, newNode)
		}
	} else if bytes.Compare(newNode.Key.K, root.Key.K) > 0 {
		right := atomic.LoadPointer(&root.Right)
		if right == nil {
			if atomic.CompareAndSwapPointer(&root.Right, nil, unsafe.Pointer(newNode)) {
				return true
			}
		} else {
			return bst.put(right, newNode)
		}
	} else {
		// If the keys are equal, append the new value to the existing key's values

		root.Key.Latch.Lock()

		root.Key.Values = append(root.Key.Values, newNode.Key.Values[0])
		root.Key.Latch.Unlock()

		return true
	}
	return false
}

type NodePos int

const (
	Left NodePos = iota
	Right
	Root
)

// Print displays the BST values in-order
func (bst *BST) Print() {
	root := atomic.LoadPointer(&bst.Root)
	bst.print((*Node)(root), Root)
}

func (bst *BST) print(node *Node, pos NodePos) {
	if node == nil {
		return
	}
	bst.print((*Node)(node.Left), Left)
	switch pos {
	case Left:
		println("L: ", string(node.Key.K))
	case Right:
		println("R: ", string(node.Key.K))
	case Root:
		println("ROOT: ", string(node.Key.K))
	}
	bst.print((*Node)(node.Right), Right)
}

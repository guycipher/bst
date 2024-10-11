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

// put adds a new key to BST or append value to existing key
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

// Get retrieves a key from the BST
func (bst *BST) Get(key []byte) *Key {
	root := atomic.LoadPointer(&bst.Root)
	return bst.get((*Node)(root), key)
}

// get retrieves a key from the BST
func (bst *BST) get(node *Node, key []byte) *Key {
	if node == nil {
		return nil
	}

	if bytes.Compare(key, node.Key.K) < 0 {
		return bst.get((*Node)(node.Left), key)
	} else if bytes.Compare(key, node.Key.K) > 0 {
		return bst.get((*Node)(node.Right), key)
	}

	return node.Key
}

// Remove removes a value from a key
func (bst *BST) Remove(key, value []byte) {
	root := atomic.LoadPointer(&bst.Root)
	bst.remove((*Node)(root), key, value)
}

// remove removes a value from a key
func (bst *BST) remove(node *Node, key, value []byte) {
	if node == nil {
		return
	}

	if bytes.Compare(key, node.Key.K) < 0 {
		bst.remove((*Node)(node.Left), key, value)
	} else if bytes.Compare(key, node.Key.K) > 0 {
		bst.remove((*Node)(node.Right), key, value)
	} else {
		node.Key.Latch.Lock()
		for i, v := range node.Key.Values {
			if bytes.Compare(v, value) == 0 {
				node.Key.Values = append(node.Key.Values[:i], node.Key.Values[i+1:]...)
				break
			}
		}
		node.Key.Latch.Unlock()
	}
}

// Delete removes a key from the BST
func (bst *BST) Delete(key []byte) {
	root := (*Node)(atomic.LoadPointer(&bst.Root))
	newRoot := bst.delete(root, key)
	atomic.StorePointer(&bst.Root, unsafe.Pointer(newRoot))
}

// delete removes a node from the BST
func (bst *BST) delete(node *Node, key []byte) *Node {
	if node == nil {
		return nil
	}

	if bytes.Compare(key, node.Key.K) < 0 {
		node.Left = unsafe.Pointer(bst.delete((*Node)(node.Left), key))
	} else if bytes.Compare(key, node.Key.K) > 0 {
		node.Right = unsafe.Pointer(bst.delete((*Node)(node.Right), key))
	} else {
		// node with only one child or no child
		if node.Left == nil {
			return (*Node)(node.Right)
		} else if node.Right == nil {
			return (*Node)(node.Left)
		}

		// node with two children: get the inorder successor (smallest in the right subtree)
		minNode := bst.minValueNode((*Node)(node.Right))

		// copy the inorder successor's content to this node
		node.Key = minNode.Key

		// delete the inorder successor
		node.Right = unsafe.Pointer(bst.delete((*Node)(node.Right), minNode.Key.K))
	}
	return node
}

// minValueNode gets the node with minimum key value found in that tree. The tree argument is pointer to the root node of the tree.
func (bst *BST) minValueNode(node *Node) *Node {
	current := node

	// loop down to find the leftmost leaf
	for (*Node)(current.Left) != nil {
		current = (*Node)(current.Left)
	}
	return current
}

// Range retrieves all keys within a range
func (bst *BST) Range(start, end []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.rangeKeys((*Node)(root), start, end, &keys)
	return keys
}

// rangeKeys retrieves all keys within a range
func (bst *BST) rangeKeys(node *Node, start, end []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// If the current node's key is greater than the start key, then there might be keys in the left subtree that are in the range
	if bytes.Compare(node.Key.K, start) > 0 {
		bst.rangeKeys((*Node)(node.Left), start, end, keys)
	}

	// If the current node's key is within the range, add it to the keys slice
	if bytes.Compare(node.Key.K, start) >= 0 && bytes.Compare(node.Key.K, end) <= 0 {
		*keys = append(*keys, node.Key)
	}

	// If the current node's key is less than the end key, then there might be keys in the right subtree that are in the range
	if bytes.Compare(node.Key.K, end) < 0 {
		bst.rangeKeys((*Node)(node.Right), start, end, keys)
	}
}

// GreaterThan retrieves all keys greater than the specified key
func (bst *BST) GreaterThan(key []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.greaterThan((*Node)(root), key, &keys)
	return keys
}

// greaterThan is a helper function to find keys greater than the specified key
func (bst *BST) greaterThan(node *Node, key []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// If the current node's key is greater than the specified key,
	// we need to check the left subtree first (for potentially smaller keys)
	if bytes.Compare(node.Key.K, key) > 0 {
		bst.greaterThan((*Node)(node.Left), key, keys)

		// Since the current node's key is greater, add it to the keys slice
		*keys = append(*keys, node.Key)

		// Continue searching in the right subtree for more keys
		bst.greaterThan((*Node)(node.Right), key, keys)
	} else {
		// If the current node's key is not greater, only search in the right subtree
		bst.greaterThan((*Node)(node.Right), key, keys)
	}
}

// GreaterThanEq retrieves all keys greater than or equal to the specified key
func (bst *BST) GreaterThanEq(key []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.greaterThanEq((*Node)(root), key, &keys)
	return keys
}

// greaterThanEq is a helper function to find keys greater than or equal to the specified key
func (bst *BST) greaterThanEq(node *Node, key []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// If the current node's key is greater than or equal to the specified key,
	// we need to check the left subtree first (for potentially smaller keys)
	if bytes.Compare(node.Key.K, key) >= 0 {
		// Include the current node's key
		*keys = append(*keys, node.Key)

		// Continue searching in the left subtree for more keys
		bst.greaterThanEq((*Node)(node.Left), key, keys)

		// Search in the right subtree for additional greater keys
		bst.greaterThanEq((*Node)(node.Right), key, keys)
	} else {
		// If the current node's key is less than the specified key, only search in the right subtree
		bst.greaterThanEq((*Node)(node.Right), key, keys)
	}
}

// LessThan retrieves all keys less than the specified key
func (bst *BST) LessThan(key []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.lessThan((*Node)(root), key, &keys)
	return keys
}

// lessThan is a helper function to find keys less than the specified key
func (bst *BST) lessThan(node *Node, key []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// If the current node's key is less than the specified key,
	// we need to check the left subtree first (for potentially smaller keys)
	if bytes.Compare(node.Key.K, key) < 0 {
		*keys = append(*keys, node.Key)

		// Continue searching in the left subtree
		bst.lessThan((*Node)(node.Left), key, keys)

		// Search in the right subtree for more keys that might also be less
		bst.lessThan((*Node)(node.Right), key, keys)
	} else {
		// If the current node's key is not less, only search in the left subtree
		bst.lessThan((*Node)(node.Left), key, keys)
	}
}

// LessThanEq retrieves all keys less than or equal to the specified key
func (bst *BST) LessThanEq(key []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.lessThanEq((*Node)(root), key, &keys)
	return keys
}

// lessThanEq is a helper function to find keys less than or equal to the specified key
func (bst *BST) lessThanEq(node *Node, key []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// If the current node's key is less than or equal to the specified key,
	// we need to check the left subtree first (for potentially smaller keys)
	if bytes.Compare(node.Key.K, key) <= 0 {
		*keys = append(*keys, node.Key)

		// Continue searching in the left subtree
		bst.lessThanEq((*Node)(node.Left), key, keys)

		// Search in the right subtree for more keys that might also be less than or equal
		bst.lessThanEq((*Node)(node.Right), key, keys)
	} else {
		// If the current node's key is greater, only search in the left subtree
		bst.lessThanEq((*Node)(node.Left), key, keys)
	}
}

// NGet retrieves all keys except the specified key
func (bst *BST) NGet(key []byte) []*Key {
	var keys []*Key
	root := atomic.LoadPointer(&bst.Root)
	bst.nGet((*Node)(root), key, &keys)
	return keys
}

// nGet is a helper function to find all keys except the specified key
func (bst *BST) nGet(node *Node, key []byte, keys *[]*Key) {
	if node == nil {
		return
	}

	// Check the left subtree first
	bst.nGet((*Node)(node.Left), key, keys)

	// If the current node's key does not match the specified key, add it to the keys slice
	if bytes.Compare(node.Key.K, key) != 0 {
		*keys = append(*keys, node.Key)
	}

	// Check the right subtree
	bst.nGet((*Node)(node.Right), key, keys)
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

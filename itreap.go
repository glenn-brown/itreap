// Copyright (c) 2012 by Glenn Brown.  All rights reserved.  See LICENSE.

// Package itreap implements an immutable ordered list.  Because the
// list is immutable, the Insert() and Remove() operations do not
// modify the original list, but return a new list with the node
// inserted or removed in O(log(N)) time where N is the
// number of nodes in the tree.
//
package itreap

//
// A treap is simultaneously a tree and a heap. Each time a value is
// inserted, its node is assigned a random priority.  Tree nodes are
// sorted by value, and the heap has higher priority values nearer the
// root.  This keeps the tree balanced regardless which values are inserted.

import (
	"fmt"
	"github.com/glenn-brown/ordinal"
	"math/rand"
)

// Type T is an immutable ordered list.
//
type T struct {
	count       int
	priority    int32
	value       interface{}
	score       float64
	left, right *T
}

// Return nil, the empty immutable list.
//
func New() *T { return nil }

// Move sorted-but-misprioritized node t in sorted tree t down to its
// appropriate heap level in heap t.  t.left and t.right are valid
// treaps.
//
func (t *T) prioritize() *T {
	if nil == t {
		return t
	}
	left, right := t.left, t.right
	if nil == left || left.priority <= t.priority {
		if nil == right || right.priority <= t.priority {
			return t
		}
		goto right
	}
	if nil == right || right.priority <= t.priority || left.priority > right.priority {
		return &T{
			t.count,
			left.priority,
			left.value,
			left.score,
			left.left,
			(&T{1 + sum(left.right, t.right),
				t.priority,
				t.value,
				t.score,
				left.right,
				t.right}).prioritize()}
	}
right:
	return &T{
		t.count,
		right.priority,
		right.value,
		right.score,
		(&T{1 + sum(t.left, right.left), t.priority, t.value,
			t.score, t.left, right.left}).prioritize(),
		right.right}
}

// Contains returns true iff the tree contains the specified value, in O(log(N)) time.
//
func (a *T) Contains(value interface{}) bool {
	if a == nil {
		return false
	}
	lessFn, s := ordinal.FnScore(value)
	for {
		switch {
		case a == nil:
			return false
		case s < a.score:
			a = a.left
			continue
		case a.score < s:
			a = a.right
			continue
		case lessFn(value, a.value):
			a = a.left
			continue
		case lessFn(a.value, value):
			a = a.right
			continue
		default:
			return true
		}
	}
	panic("never")
}

// Insert returns a new tree like the original, but with the value inserted, in O(log(N)) time.
//
func (t *T) Insert(value interface{}) *T {
	less, score := ordinal.FnScore(value)
	nu := &T{1, rand.Int31(), value, score, nil, nil}
	return t.insert(nu, less)
}

// Return a new immutable treap like treap t, but with node nu inserted, in O(log(N)) time.
//
func (t *T) insert(nu *T, less func(a, b interface{}) bool) *T {
	if nil == t {
		return nu
	}

	// Insert on left if less than root, and right if greater, taking care to
	// handle score cases first for performance.

	if nu.score < t.score {
		goto left
	}
	if t.score < nu.score || less(t.value, nu.value) {
		right := t.right.insert(nu, less)
		if right.priority > t.priority {
			// Rotate left, replacing t.right with right.
			return &T{
				t.count + 1, right.priority, right.value, right.score,
				&T{1 + sum(t.left, right.left), t.priority, t.value, t.score, t.left, right.left},
				right.right}
		}
		return &T{t.count + 1, t.priority, t.value, t.score, t.left, right}
	}
left:
	left := t.left.insert(nu, less)
	if left.priority > t.priority {
		// Rotate right, replacing t.left with left.
		return &T{
			t.count + 1, left.priority, left.value, left.score, left.left,
			&T{1 + sum(left.right, t.right), t.priority, t.value, t.score,
				left.right, t.right}}
	}
	return &T{t.count + 1, t.priority, t.value, t.score, left, t.right}
}

// Remove returns a new treap like the original, but with the value removed, in O(log(N)) time.
// If there is no matching value to remove, the original tree is returned.
// If there are multiple matching values, only one is removed.
//
func (t *T) Remove(value interface{}) *T {
	less, score := ordinal.FnScore(value)
	rv, ok := t.remove(value, score, less)
	if !ok {
		return t
	}
	return rv
}

func (t *T) remove(value interface{}, score float64, less func(a, b interface{}) bool) (*T, bool) {
	if nil == t {
		return nil, false
	}
	if score < t.score {
		goto left
	}
	if t.score < score || less(t.value, value) {
		right, ok := t.right.remove(value, score, less)
		return &T{t.count - 1, t.priority, t.value, t.score, t.left, right}, ok
	}
	if !less(value, t.value) {
		return t.removeNode(), true
	}
left:
	left, ok := t.left.remove(value, score, less)
	return &T{t.count - 1, t.priority, t.value, t.score, left, t.right}, ok
}

func (t *T) removeNode() *T {
	left, right := t.left, t.right
	if nil == left {
		return right
	}
	if nil == right {
		return left
	}
	// Find and remove the successor node.
	n, right := right.removeLeftmost()
	// Repace the top (removed) node with the successor, and restore priority.
	return (&T{t.count - 1, n.priority, n.value, n.score, left, right}).prioritize()
}

func (t *T) removeLeftmost() (left *T, after *T) {
	if nil == t {
		return nil, t
	}
	if nil == t.left {
		return t, t.right
	}
	n, left := t.left.removeLeftmost()
	return n, &T{t.count - 1, t.priority, t.value, t.score, left, t.right}
}

func (t *T) removeRightmost() (right *T, after *T) {
	if nil == t {
		return nil, t
	}
	if nil == t.right {
		return t, t.left
	}
	n, right := t.right.removeRightmost()
	return n, &T{t.count - 1, t.priority, t.value, t.score, t.left, right}
}

// Len returns the number of values in the list.
//
func (t *T) Len() int {
	if nil == t {
		return 0
	}
	return t.count
}

// RemoveN removes the nth element from the list, returning the
// modified list and removed value.  Use t.RemoveN(0) to pop the first (least)
// value and t.RemoveN(t.Len()-1) to remove the last (greatest).
//
func (t *T) RemoveN(n int) (nu *T, val interface{}) {
	if nil == t || n < 0 || t.count <= n {
		return nil, nil
	}
	lcount := 0
	if nil != t.left {
		lcount = t.left.count
	}
	if n < lcount {
		left, val := t.left.RemoveN(n)
		return &T{t.count - 1, t.priority, t.value, t.score, left, t.right}, val
	}
	if n > lcount {
		right, val := t.right.RemoveN(n - lcount - 1)
		return &T{t.count - 1, t.priority, t.value, t.score, t.left, right}, val
	}
	return t.removeNode(), t.value
}

// Return the value at position n in the list.  The index n must be in the interval
// [0,t.Len()).
//
func (t *T) GetN(n int) (value interface{}) {
	if nil == t {
		return nil
	}
	lcount := 0
	if nil != t.left {
		lcount = t.left.count
	}
	if n < lcount {
		return t.left.GetN(n)
	}
	if lcount < n {
		return t.right.GetN(n - lcount - 1)
	}
	return t.value
}

// Return a string representation of the immutable treap.
//
func (t *T) String() string {
	if nil == t {
		return ""
	}
	left, right := t.left, t.right
	if nil == left && nil == right {
		return fmt.Sprintf("%v", t.value)
	}
	if nil == left {
		return fmt.Sprintf("%v %v", t.value, right)
	}
	if nil == right {
		return fmt.Sprintf("%v %v", left, t.value)
	}
	return fmt.Sprintf("%v %v %v", left, t.value, right)
}

func sum(a, b *T) (count int) {
	if nil != a {
		count += a.count
	}
	if nil != b {
		count += b.count
	}
	return count
}

// Copyright (c) 2012, Glenn Brown.  All rights reserved.  See LICENSE.

package itreap

import "fmt"

// For any old type:
type MyType struct{ a, b int }

// Implement the SlowKey interface:
func (a *MyType) Less(b interface{}) bool {
	// For example, sort by the sum of the elements in the struct:
	mb := b.(*MyType)
	return (a.a + a.b) < (mb.a + mb.b)
}

// Any type implementing the SlowKey interface can be used as a key.
//
func ExampleSlowKey() {
	keys := []MyType{{1, 2}, {5, 6}, {3, 4}}
	fmt.Print(New().Insert(&keys[0]).Insert(&keys[1]).Insert(&keys[2]))
	// Output: &{1 2} &{3 4} &{5 6}
}

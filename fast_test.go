// Copyright (c) 2012, Glenn Brown.  All rights reserved.  See LICENSE.

package itreap

import "fmt"

// For any old type:
type FastType struct{ a, b int }

// Implement the SlowKey interface:
func (a *FastType) Less(b interface{}) bool {
	// For example, sort by the sum of the elements in the struct:
	mb := b.(*FastType)
	return (a.a + a.b) < (mb.a + mb.b)
}
func (*FastType) Score(i interface{}) float64 {
	// Score(i) increase monotonically with increasing key value.
	m := i.(*FastType)
	return float64(m.a + m.b)
}

// Any type implementing the FastKey interface can be used as a key.
//
func ExampleFastKey() {
	keys := []FastType{{1, 2}, {5, 6}, {3, 4}}
	fmt.Print(New().Insert(&keys[0]).Insert(&keys[1]).Insert(&keys[2]))
	// Output: &{1 2} &{3 4} &{5 6}
}

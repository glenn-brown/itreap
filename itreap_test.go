package itreap

import (
	"math/rand"
	"testing"
)

func itreap(n int) *T {
	a := rand.Perm(n)
	rv := New()
	for _, v := range a {
		rv = rv.Insert(v)
	}
	return rv
}

func (i *T) verifyCounts(t *testing.T) {
	if nil == i {
		return
	}
	s := sum(i.left, i.right)
	if i.count != 1+s {
		t.Errorf("%v != %v", i.count, s)
	}
	i.left.verifyCounts(t)
	i.left.verifyCounts(t)
}

func TestT_Insert(t *testing.T) {
	t.Parallel()
	a := rand.Perm(6)
	i := New()
	for _, v := range a {
		before := i.String()
		nu := i.Insert(v)
		t.Log(i)
		if nil == nu {
			t.Error("Insert returned nil tree.")
		}
		after := i.String()
		if before != after {
			t.Error(before + " != " + after)
		}
		next := nu.String()
		if before == next {
			t.Error(before + " == " + next)
		}
		i = nu
		i.verifyCounts(t)
	}
	s := i.String()
	x := "0 1 2 3 4 5"
	if s != x {
		t.Error(s + " != " + x)
	}
}

func TestT_Remove(t *testing.T) {
	t.Parallel()

	// Remove entries one at a time, confirming that an entry is removed each time.
	// If a wrong entry is removed, a later remove will fail.

	i := itreap(100)
	rm := rand.Perm(100)
	for _, v := range rm {
		before := i.String()
		t.Log(before)
		nu := i.Remove(v)
		after := i.String()
		if before != after {
			t.Error(before + " != " + after)
		}
		next := nu.String()
		if before == next {
			t.Error(before + " == " + next)
		}
		i = nu
		i.verifyCounts(t)
	}
	final := i.String()
	if final != "" {
		t.Error(final + " != ")
	}

	// Remove first/middle/last entries and check that the right one was removed.

	i = itreap(11)
	i = i.Remove(0)
	if nil == i {
		t.Error("Remove(0)")
	}
	i = i.Remove(5)
	if nil == i {
		t.Error("Remove(5)")
	}
	i = i.Remove(10)
	if nil == i {
		t.Error("Remove(10)")
	}
	s := i.String()
	x := "1 2 3 4 6 7 8 9"
	if s != x {
		t.Error(s + " != " + x)
	}
}

func TestT_RemoveN(s *testing.T) {
	s.Parallel()
	t := itreap(100)
	for i := 100; i > 0; i-- {
		if t.Len() != i {
			s.Error(t.Len())
		}
		var val interface{}
		t, val = t.RemoveN(rand.Intn(i))
		if nil == val {
			s.Error(val)
		}
	}
	if t.Len() != 0 {
		s.Error("t.Len() != 0")
	}
}

func TestT_GetN(s *testing.T) {
	s.Parallel()
	t := itreap(100)
	for i := t.Len() - 1; i >= 0; i-- {
		g := t.GetN(i).(int)
		if i != g {
			s.Error(i, " != ", g)
		}
	}
}

func BenchmarkT_Contains(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	p := rand.Perm(b.N)
	b.StartTimer()
	for _, v := range p {
		t.Contains(v)
	}
}

func BenchmarkT_GetN_first(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t.GetN(0)
	}
}

func BenchmarkT_GetN_last(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t.GetN(t.Len() - 1)
	}
}

func BenchmarkT_GetN_random(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	p := rand.Perm(b.N)
	b.StartTimer()
	for _, v := range p {
		t.GetN(v)
	}
}

func BenchmarkT_Insert(b *testing.B) {
	b.StopTimer()
	in := rand.Perm(b.N)
	t := New()
	b.StartTimer()
	for _, v := range in {
		t = t.Insert(v)
	}
}

func BenchmarkT_Remove(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	out := rand.Perm(b.N)
	b.StartTimer()
	for _, v := range out {
		t = t.Remove(v)
	}
}

func BenchmarkT_RemoveN_first(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t, _ = t.RemoveN(0)
	}
}

func BenchmarkT_RemoveN_last(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t, _ = t.RemoveN(t.Len() - 1)
	}
}

func BenchmarkT_RemoveN_mid(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t, _ = t.RemoveN(i / 2)
	}
}

func BenchmarkT_RemoveN_random(b *testing.B) {
	b.StopTimer()
	t := itreap(b.N)
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		t, _ = t.RemoveN(rand.Intn(i))
	}
}

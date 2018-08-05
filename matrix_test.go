package lev

import "testing"

func TestMatrixIndex(t *testing.T) {
	var m matrix
	m.reset(3, 5)
	var v int
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			m.set(i, j, v)
			v++
		}
	}
	v = 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			if got := m.at(i, j); got != v {
				t.Fatalf("expected %d; got %d", v, got)
			}
			v++
		}
	}
}

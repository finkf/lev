package lev

import (
	"fmt"
	"testing"
)

func TestWLev(t *testing.T) {
	tests := []struct {
		a, b []string
		want int
	}{
		{[]string{"a", "b", "c"}, []string{"a", "c"}, 1},
		{[]string{"foo", "bar", "baz"}, []string{"foo", "baz", "bar"}, 2},
		{[]string{}, []string{}, 0},
		{[]string{}, []string{"abc"}, 3},
		{[]string{"xyz"}, []string{"aba"}, 3},
	}
	var lev Lev
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v-%v", tc.a, tc.b), func(t *testing.T) {
			var wlev WLev
			a, b := StringArray(&lev, tc.a...), StringArray(&lev, tc.b...)
			if got := wlev.EditDistance(a, b); got != tc.want {
				t.Fatalf("expected edit distance %d; got %d", tc.want, got)
			}
		})
	}
}

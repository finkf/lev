package lev_test

import (
	"testing"

	"github.com/finkf/lev"
)

func TestLevCalculate(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"abc", "abb", 1},
		{"abc", "bbc", 1},
		{"abb", "abc", 1},
		{"bbc", "abc", 1},
		{"für", "fur", 1},
		{"Bäume", "Träume", 2},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			got := l.EditDistance(tc.s1, tc.s2)
			if got != tc.want {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
		})
	}
}

func TestLevBacktrace(t *testing.T) {
	tests := []struct {
		s1, s2, want string
	}{
		{"", "", ""},
		{"abc", "abc", "|||"},
		{"abd", "abc", "||#"},
		{"ac", "abc", "|+|"},
		{"abc", "ac", "|-|"},
		{"xxabc", "abc", "--|||"},
		{"abc", "xxabc", "++|||"},
		{"xabyx", "abc", "-||#-"},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			_, b := l.Backtrace(tc.s1, tc.s2)
			if got := b.String(); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestBacktraceValidate(t *testing.T) {
	tests := []struct {
		backtrace string
		ok        bool
	}{
		{"||#-#+", true},
		{"||#-xx+", false},
		{"", true},
	}
	for _, tc := range tests {
		t.Run(tc.backtrace, func(t *testing.T) {
			err := lev.Backtrace(tc.backtrace).Validate()
			if tc.ok && err != nil {
				t.Fatalf("got error %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("unexpected error %v", err)
			}
		})
	}
}
func TestLevString(t *testing.T) {
	tests := []struct {
		s1, s2, want string
	}{
		{"ab", "ab", `_ _ a b
_ 0 1 2
a 1 0 1
b 2 1 0
`},
		{"ac", "ab", `_ _ a b
_ 0 1 2
a 1 0 1
c 2 1 1
`},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			_ = l.EditDistance(tc.s1, tc.s2)
			got := l.String()
			if got != tc.want {
				t.Fatalf("expected:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}
}

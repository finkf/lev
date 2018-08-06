package lev_test

import (
	"testing"

	"github.com/finkf/lev"
)

func TestCalculate(t *testing.T) {
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

func TestTrace(t *testing.T) {
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
			_, b := l.Trace(tc.s1, tc.s2)
			if got := b.String(); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestTraceValidate(t *testing.T) {
	tests := []struct {
		trace string
		ok    bool
	}{
		{"||#-#+", true},
		{"||#-xx+", false},
		{"", true},
	}
	for _, tc := range tests {
		t.Run(tc.trace, func(t *testing.T) {
			err := lev.Trace(tc.trace).Validate()
			if tc.ok && err != nil {
				t.Fatalf("got error %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("unexpected error %v", err)
			}
		})
	}
}

func TestAlignment(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   lev.Alignment
	}{
		{"", "", lev.Alignment{}},
		{"abc", "abc", lev.Alignment{
			Trace: "|||", S1: "abc", S2: "abc"}},
		{"ab", "abc", lev.Alignment{
			Trace: "||+", S1: "ab~", S2: "abc", Distance: 1}},
		{"abc", "ab", lev.Alignment{
			Trace: "||-", S1: "abc", S2: "ab~", Distance: 1}},
		{"abc", "abd", lev.Alignment{
			Trace: "||#", S1: "abc", S2: "abd", Distance: 1}},
		{"", "abc", lev.Alignment{
			Trace: "+++", S1: "~~~", S2: "abc", Distance: 3}},
		{"abc", "", lev.Alignment{
			Trace: "---", S1: "abc", S2: "~~~", Distance: 3}},
		{"abc", "xyz", lev.Alignment{
			Trace: "###", S1: "abc", S2: "xyz", Distance: 3}},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			got, err := l.Alignment(l.Trace(tc.s1, tc.s2))
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestInvalidAlignment(t *testing.T) {
	tests := []struct {
		s1, s2 string
		trace  lev.Trace
	}{
		{"abc", "abc", lev.Trace("||x")},
		{"abc", "abc", lev.Trace("||||")},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			d := l.EditDistance(tc.s1, tc.s2)
			_, err := l.Alignment(d, tc.trace)
			if err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestString(t *testing.T) {
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

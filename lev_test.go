package lev_test

import (
	"reflect"
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

func newAlignment(t, s1, s2 string, d int) lev.Alignment {
	return lev.Alignment{
		Trace:    []byte(t),
		S1:       []rune(s1),
		S2:       []rune(s2),
		Distance: d,
	}
}

func TestAlignment(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   lev.Alignment
	}{
		{"", "", newAlignment("", "", "", 0)},
		{"abc", "abc", newAlignment("|||", "abc", "abc", 0)},
		{"ab", "abc", newAlignment("||+", "ab~", "abc", 1)},
		{"abc", "ab", newAlignment("||-", "abc", "ab~", 1)},
		{"abc", "abd", newAlignment("||#", "abc", "abd", 1)},
		{"", "abc", newAlignment("+++", "~~~", "abc", 3)},
		{"abc", "", newAlignment("---", "abc", "~~~", 3)},
		{"abc", "xyz", newAlignment("###", "abc", "xyz", 3)},
		{"file://a.txt", "Der alte Mann", newAlignment(
			"##+++|+|##-|##--", "fi~~~l~e://a.txt", "Der alte M~ann~~", 13)},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			got, err := l.Alignment(l.Trace(tc.s1, tc.s2))
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %q; got %q", tc.want, got)
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
		{"abc", "abc", lev.Trace("||--")},
		{"a", "abc", lev.Trace("++++++")},
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

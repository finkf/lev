package lev

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Lev holds two aligned strings and the cost matrix.
type Lev struct {
	matrix
	s1, s2 []rune
}

// EditDistance calculates the minimal edit distance
// between the two given strings.
func (l *Lev) EditDistance(s1, s2 string) int {
	m, n := l.init(s1, s2)
	for i := 0; i < m+1; i++ {
		l.set(i, 0, i)
	}
	for i := 0; i < n+1; i++ {
		l.set(0, i, i)
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			v, _, _, _ := l.argMin(i, j, l.cost)
			l.set(i, j, v)
		}
	}
	return l.at(len(l.s1), len(l.s2))
}

func (l *Lev) cost(op byte, i, j int) int {
	switch op {
	case Del, Ins:
		return 1
	default:
		if l.s1[i-1] == l.s2[j-1] {
			return 0
		}
		return 1
	}
}

// String returns the matrix format of the last
// calculated edit distance.
func (l *Lev) String() string {
	const (
	// flags = tabwriter.AlignRight | tabwriter.Debug
	)
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	// write header
	fmt.Fprintf(w, "_\t_")
	for i := 0; i < len(l.s2); i++ {
		fmt.Fprintf(w, "\t%c", l.s2[i])
	}
	fmt.Fprint(w, "\n")
	// write lines
	for i := 0; i < len(l.s1)+1; i++ {
		if i > 0 {
			fmt.Fprintf(w, "%c", l.s1[i-1])
		} else {
			fmt.Fprintf(w, "_")
		}
		for j := 0; j < len(l.s2)+1; j++ {
			fmt.Fprintf(w, "\t%d", l.at(i, j))
		}
		fmt.Fprint(w, "\n")
	}
	w.Flush()
	return buf.String()
}

func (l *Lev) init(s1, s2 string) (int, int) {
	l.s1 = []rune(s1)
	l.s2 = []rune(s2)
	m := len(l.s1)
	n := len(l.s2)
	l.reset(m+1, n+1)
	return m, n
}

// Trace defines an array of edit operations.
type Trace []byte

const (
	// Del marks a deletion in s2.
	Del = byte('-')
	// Sub marks a substitution.
	Sub = byte('#')
	// Ins marks a insertion in s1.
	Ins = byte('+')
	// Nop marks a non-edit operation.
	Nop = byte('|')
	// Mis marks a missing character in s1 or s2.
	Mis = byte('~')
)

// Trace returns the edit distance between the two given strings
// and the trace of the according edit operations.
func (l *Lev) Trace(s1, s2 string) (int, Trace) {
	d := l.EditDistance(s1, s2)
	return d, l.trace(l.cost)
}

// Validate returns an error if the trace is not valid.
func (b Trace) Validate() error {
	for _, op := range b {
		switch op {
		case Del, Ins, Sub, Nop:
		default:
			return fmt.Errorf("invalid edit operation %c", op)
		}
	}
	return nil
}

func (b Trace) String() string {
	return string(b)
}

func max(m, n int) int {
	if m > n {
		return m
	}
	return n
}

// Alignment captures the alignment of two strings with
// the accoriding trace of edit operations.
type Alignment struct {
	S1, S2   []rune
	Trace    Trace
	Distance int
}

// NewAlignment creates a new Alignment from two and a given trace.
// The error distance is calculated using the trace directly.
func NewAlignment(s1, s2, trace string) (Alignment, error) {
	a := Alignment{S1: []rune(s1), S2: []rune(s2), Trace: []byte(trace)}
	if len(a.S1) != len(a.Trace) || len(a.S2) != len(a.Trace) {
		return Alignment{},
			fmt.Errorf("trace and/or string lengths do not match")
	}
	for _, op := range a.Trace {
		switch op {
		case Del, Ins, Sub:
			a.Distance++
		case Nop:
		default:
			return Alignment{}, fmt.Errorf("invalid trace: %s", a.Trace)
		}
	}
	return a, nil
}

// Alignment returns the given alignment strings and the according
// trace.
func (l *Lev) Alignment(d int, b Trace) (Alignment, error) {
	a := Alignment{Distance: d, Trace: b}
	i, j := 0, 0
	n := max(len(l.s1), len(l.s2))
	a.S1 = make([]rune, 0, n)
	a.S2 = make([]rune, 0, n)
	for _, c := range b {
		switch c {
		case Nop, Sub:
			if i >= len(l.s1) || j >= len(l.s2) {
				return l.alignmentError(b)
			}
			a.S1 = append(a.S1, l.s1[i])
			a.S2 = append(a.S2, l.s2[j])
			i, j = i+1, j+1
		case Ins:
			if j >= len(l.s2) {
				return l.alignmentError(b)
			}
			a.S1 = append(a.S1, rune(Mis))
			a.S2 = append(a.S2, l.s2[j])
			j++
		case Del:
			if i >= len(l.s1) {
				return l.alignmentError(b)
			}
			a.S1 = append(a.S1, l.s1[i])
			a.S2 = append(a.S2, rune(Mis))
			i++
		default:
			return l.alignmentError(b)
		}
	}
	return a, nil
}

func (l *Lev) alignmentError(b Trace) (Alignment, error) {
	var a Alignment
	return a, fmt.Errorf("align %q, %q: %q",
		string(l.s1), string(l.s2), b)
}

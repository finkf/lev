package lev

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"unicode/utf8"
)

// Lev holds two aligned strings and the weight matrix.
type Lev struct {
	matrix
	s1, s2 []rune
}

// EditDistance calculates the minimal edit distance
// between the two given strings.
func (l *Lev) EditDistance(s1, s2 string) int {
	l.s1 = toRunes(s1)
	l.s2 = toRunes(s2)
	m := len(l.s1)
	n := len(l.s2)
	l.reset(m+1, n+1)
	for i := 0; i < m+1; i++ {
		l.set(i, 0, i)
	}
	for i := 0; i < n+1; i++ {
		l.set(0, i, i)
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			if l.s1[i-1] == l.s2[j-1] {
				l.set(i, j, l.at(i-1, j-1))
			} else {
				v, _, _, _ := l.argMin(i, j)
				l.set(i, j, v)
			}
		}
	}
	// m = len(l.ws1) + 1, n = len(l.ws2) + 1
	return l.at(len(l.s1), len(l.s2))
}

func (l *Lev) argMin(i, j int) (v, ii, jj int, op byte) {
	// no deletion or substitution possible
	if i < 1 {
		return l.at(i, j-1) + 1, i, j - 1, Ins
	}
	// no insertion or substitution possible
	if j < 1 {
		return l.at(i-1, j) + 1, i - 1, j, Del
	}
	if l.s1[i-1] == l.s2[j-1] {
		return l.at(i-1, j-1), i - 1, j - 1, Nop
	}
	sub := l.at(i-1, j-1)
	ins := l.at(i, j-1)
	del := l.at(i-1, j)
	if sub < ins {
		if sub < del {
			return sub + 1, i - 1, j - 1, Sub
		}
		return del + 1, i - 1, j, Del
	}
	if ins < del {
		return ins + 1, i, j - 1, Ins
	}
	return del + 1, i - 1, j, Del
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

// toRunes converts a string to an array of runes.
func toRunes(str string) []rune {
	n := utf8.RuneCountInString(str)
	runes := make([]rune, 0, n)
	for _, r := range str {
		runes = append(runes, r)
	}
	return runes
}

// // EditOperation defines the type of an edit operation.
// type EditOperation byte

// Backtrace defines an array of edit operations.
type Backtrace []byte

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

// Backtrace returns the edit distance between the two given strings
// and the backtrace of the according edit operations.
func (l *Lev) Backtrace(s1, s2 string) (int, Backtrace) {
	d := l.EditDistance(s1, s2)
	return d, l.calculateBacktrace()
}

func (l *Lev) calculateBacktrace() Backtrace {
	length := max(len(l.s1), len(l.s2))
	b := make(Backtrace, 0, length)
	// m = len(l.ws1) + 1, n = len(l.ws2) + 1
	for i, j := len(l.s1), len(l.s2); i > 0 || j > 0; {
		_, ii, jj, op := l.argMin(i, j)
		b = append(b, op)
		i = ii
		j = jj
	}
	return b.reverse()
}

// Validate returns an error if the backtrace is not valid.
func (b Backtrace) Validate() error {
	for _, op := range b {
		switch op {
		case Del, Ins, Sub, Nop:
		default:
			return fmt.Errorf("invalid edit operation %c", op)
		}
	}
	return nil
}

func (b Backtrace) String() string {
	return string(b)
}

func (b Backtrace) reverse() Backtrace {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func max(m, n int) int {
	if m > n {
		return m
	}
	return n
}

// Alignment captures the alignment of two strings with
// the accoriding backtrace of edit operations.
type Alignment struct {
	S1, S2, Backtrace string
	Distance          int
}

// Alignment returns the given alignment strings and the according
// backtrace.
func (l *Lev) Alignment(d int, b Backtrace) (Alignment, error) {
	a := Alignment{Distance: d}
	if err := b.Validate(); err != nil {
		return a, err
	}
	if len(l.s1) != len(b) && len(l.s2) != len(b) {
		return a, fmt.Errorf("invalid backtrace")
	}
	var b1, b2 bytes.Buffer
	i, j := 0, 0
	for _, c := range b {
		switch c {
		case Nop, Sub:
			b1.WriteRune(l.s1[i])
			b2.WriteRune(l.s2[j])
			i, j = i+1, j+1
		case Ins:
			b1.WriteByte(Mis)
			b2.WriteRune(l.s2[j])
			j++
		case Del:
			b1.WriteRune(l.s1[i])
			b2.WriteByte(Mis)
			i++
		}
	}
	a.S1, a.S2, a.Backtrace = b1.String(), b2.String(), b.String()
	return a, nil
}

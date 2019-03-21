package lev

// Array defines the interface to calculate the costs for the
// levenshtein distance and wagner-fisher algorithm.
type Array interface {
	// Return the length of the array
	Len() int
	// Return the cost between self[i] with a[j].  If `a` is nil this
	// method should return the cost for an insertion or deletion at
	// `i` (`j` is set to -1 in theses cases).
	Cost(a Array, i int, j int) int
}

// WLev holds two aligned Arrays and the cost matrix.
type WLev struct {
	matrix
	a, b Array
}

// EditDistance returns the cost between the two given arrays.
func (l *WLev) EditDistance(a, b Array) int {
	m, n := l.init(a, b)
	for i := 1; i < m+1; i++ {
		l.set(i, 0, a.Cost(nil, i-1, -1)+l.at(i-1, 0))
	}
	for i := 1; i < n+1; i++ {
		l.set(0, i, b.Cost(nil, i-1, -1)+l.at(0, i-1))
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			v, _, _, _ := l.argMin(i, j, l.cost)
			l.set(i, j, v)
		}
	}
	return l.at(l.a.Len(), l.b.Len())
}

// Trace returns the edit distance between the two given strings
// and the trace of the according edit operations.
func (l *WLev) Trace(a, b Array) (int, Trace) {
	d := l.EditDistance(a, b)
	return d, l.trace(l.cost)
}

func (l *WLev) init(a, b Array) (int, int) {
	l.a = a
	l.b = b
	m := a.Len()
	n := b.Len()
	l.reset(m+1, n+1)
	return m, n
}

func (l *WLev) cost(op byte, i, j int) int {
	switch op {
	case Del:
		return l.a.Cost(nil, i-1, -1)
	case Ins:
		return l.b.Cost(nil, j-1, -1)
	default:
		return l.a.Cost(l.b, i-1, j-1)
	}
}

// NewStringArray returns an array of words (or sentences...).  The
// array can be used to align words using Wagner-Fischer.  Costs
// calculated by the Array are based on the Levenshtein-Distance
// between the Array's tokens.
func NewStringArray(l *Lev, args ...string) StringArray {
	return StringArray{l, args}
}

// StringArray holds a list of strings.  It implements the Array
// interface.
type StringArray struct {
	lev  *Lev
	strs []string
}

// Len returns the length of this string Array.
func (s StringArray) Len() int {
	return len(s.strs)
}

// Cost returns the costs between entries in the StringArray using the
// Levenshtein-Distance.
func (s StringArray) Cost(o Array, i, j int) int {
	if o == nil {
		return len(s.strs[i])
	}
	a := o.(StringArray)
	return s.lev.EditDistance(s.strs[i], a.strs[j])
}

// At returns the string at position `i`.
func (s StringArray) At(i int) string {
	return s.strs[i]
}

var _ Array = StringArray{}

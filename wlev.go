package lev

// Array defines the interface to calculate the weighted levenshtein
// distance and wagner-fisher algorithm.
type Array interface {
	// return the length of the array
	Len() int
	// return the weight between self[i] with a[j].  If a is nil the
	// weight between self and the empty word.
	Weight(a Array, i int, j int) int
}

// WLev holds two aligned Arrays and the weight matrix.
type WLev struct {
	matrix
	a, b Array
}

// EditDistance returns the weighted edit distance between the two
// given arrays.
func (l WLev) EditDistance(a, b Array) int {
	m, n := l.init(a, b)
	for i := 1; i < m+1; i++ {
		l.set(i, 0, a.Weight(nil, i-1, 0))
	}
	for i := 1; i < n+1; i++ {
		l.set(0, i, b.Weight(nil, i-1, 0))
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			w := l.a.Weight(l.b, i-1, j-1)
			v, _, _, _ := l.argMin(i, j, w)
			l.set(i, j, v)
		}
	}
	return l.at(l.a.Len(), l.b.Len())
}

func (l *WLev) init(a, b Array) (int, int) {
	l.a = a
	l.b = b
	m := a.Len()
	n := b.Len()
	l.reset(m+1, n+1)
	return m, n
}

func StringArray(l *Lev, args ...string) Array {
	return stringa{l, args}
}

type stringa struct {
	lev  *Lev
	strs []string
}

func (s stringa) Len() int {
	return len(s.strs)
}

func (s stringa) Weight(o Array, i, j int) int {
	if o == nil {
		return len(s.strs[i])
	}
	a := o.(stringa)
	return s.lev.EditDistance(s.strs[i], a.strs[j])
}

var _ Array = stringa{}

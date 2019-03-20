package lev

// Array defines the interface to calculate the costs for the
// levenshtein distance and wagner-fisher algorithm.
type Array interface {
	// Return the length of the array
	Len() int
	// Return the cost between self[i] with a[j].  If `a` is nil this
	// method should return the cost for an insertion or deletion at
	// `i`.
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
		l.set(i, 0, a.Cost(nil, i-1, -1))
	}
	for i := 1; i < n+1; i++ {
		l.set(0, i, b.Cost(nil, i-1, -1))
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			v, _, _, _ := l.argMin(i, j, l.cost)
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

func (s stringa) Cost(o Array, i, j int) int {
	if o == nil {
		return len(s.strs[i])
	}
	a := o.(stringa)
	return s.lev.EditDistance(s.strs[i], a.strs[j])
}

var _ Array = stringa{}

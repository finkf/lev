package lev

// Array defines the interface to calculate the weighted levenshtein
// distance and wagner-fisher algorithm.
type Array interface {
	Len() int                         // return the length of the array
	Weight(a Array, i int, j int) int // return the weight between self[i] with a[j]
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
	for i := 0; i < m+1; i++ {
		l.set(i, 0, i)
	}
	for i := 0; i < n+1; i++ {
		l.set(0, i, i)
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			w := l.A.Weight(l.b, i, j)
			if w == 0 {
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

func (l *WLev) init(a, b Array) (int, int) {
	l.a = a
	l.b = b
	m := a.Len()
	n := b.Len()
	l.reset(m+1, n+1)
	return m, n
}

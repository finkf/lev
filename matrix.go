package lev

type matrix struct {
	table []int
	m, n  int
}

func (ma *matrix) reset(m, n int) {
	l := m * n
	if l > len(ma.table) {
		ma.table = make([]int, l)
	}
	ma.m = m
	ma.n = n
}

func (ma *matrix) index(i, j int) int {
	return i*ma.n + j
}

func (ma *matrix) at(i, j int) int {
	return ma.table[ma.index(i, j)]
}

func (ma *matrix) set(i, j int, v int) {
	ma.table[ma.index(i, j)] = v
}

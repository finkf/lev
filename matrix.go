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

func (ma *matrix) argMin(i, j int, c func(byte, int, int) int) (v, ii, jj int, op byte) {
	// no deletion or substitution possible
	if i < 1 {
		return ma.at(i, j-1) + c(Ins, -1, j), i, j - 1, Ins
	}
	// no insertion or substitution possible
	if j < 1 {
		return ma.at(i-1, j) + c(Del, i, -1), i - 1, j, Del
	}
	csub := c(Sub, i, j)
	sub := ma.at(i-1, j-1) + csub
	ins := ma.at(i, j-1) + c(Ins, -1, j)
	del := ma.at(i-1, j) + c(Del, i, -1)
	if sub <= ins {
		if sub <= del {
			if csub == 0 {
				return sub, i - 1, j - 1, Nop
			}
			return sub, i - 1, j - 1, Sub
		}
		return del, i - 1, j, Del
	}
	if ins < del {
		return ins, i, j - 1, Ins
	}
	return del, i - 1, j, Del
}

func (ma *matrix) trace(c func(byte, int, int) int) Trace {
	length := max(ma.m-1, ma.n-1)
	trace := make(Trace, 0, length)
	for i, j := ma.m-1, ma.n-1; i > 0 || j > 0; {
		_, ii, jj, op := ma.argMin(i, j, c)
		trace = append(trace, op) //trace[k-1] = op
		i = ii
		j = jj
	}
	for i, j := 0, len(trace)-1; i < j; i, j = i+1, j-1 {
		trace[i], trace[j] = trace[j], trace[i]
	}
	return trace
}

func max(m, n int) int {
	if m > n {
		return m
	}
	return n
}

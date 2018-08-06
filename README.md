![build status](https://travis-ci.org/finkf/lev.svg?branch=master)
# lev
Simple [go](https://golang.org) package to calculate the
[Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
between two strings.

## Usage examples
### Calculate edit distance between two strings
```golang
package main

import(
	"fmt"
	"github.com/finkf/lev"
)

func main() {
	var l lev.Lev
	d := l.EditDistance("abc", "abd")
	fmt.Printf("distance between %q and %q = %d", "abc", "abd", d)
}
```

### Calculate the trace of two strings
```golang
package main

import(
	"fmt"
	"github.com/finkf/lev"
)

func main() {
	var l lev.Lev
	_, b := l.Trace("abc", "abd")
	fmt.Printf("trace of %q and %q = %s", "abc", "abd", b)
}
```

### Calculate the alignment of two strings
```golang
package main

import(
	"fmt"
	"github.com/finkf/lev"
)

func main() {
	var l lev.Lev
	a, _ := l.Alignment(l.Trace("abc", "abd"))
	fmt.Printf("%s\n%s\n%s\n", a.S1, a.Trace, a.S2)
}
```

![build status](https://travis-ci.org/finkf/lev.svg?branch=master)
# lev
Simple [go](https://golang.org) package to calculate the
[Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
between two strings.

## Usage
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
### Calculate the backtrace of two strings
```golang
package main

import(
	"fmt"
	"github.com/finkf/lev"
)

func main() {
	var l lev.Lev
	_, b := l.Backtrace("abc", "abd")
	fmt.Printf("backtrace of %q and %q = %s", "abc", "abd", b)
}
```

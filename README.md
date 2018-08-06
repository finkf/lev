![build status](https://travis-ci.org/finkf/lev.svg?branch=master)
# lev
Simple [go](https://golang.org) package to calculate the
[Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
between two strings.

## Usage
```golang
package main

import("github.com/finkf/lev")

func main() {
	var l lev.Lev
	d := l.EditDistance("abc", "abd")
	fmt.Printf("distance between %q and %q = %d", "abc", "abd", d)
}
```

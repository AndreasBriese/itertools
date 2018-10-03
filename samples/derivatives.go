// derivatives.go
//

package main

import (
	"fmt"
	"github.com/AndreasBriese/itertools"
	// "sync"
	// "time"
)

const SAMPLELEN = 10

func main() {
	// initialize
	f1 := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f1seq := itertools.ToIterFloat64(f1).Map(func(x float64) float64 { return x * x })
	f2seq := itertools.ToIterFloat64(f1).Map(func(x float64) float64 { return 2 * x * x })

	f1f2 := itertools.ZipToIterFloat64(f1seq.List(), f2seq.List())
	f1f2 = f1f2.PairOp(func(a, b float64) float64 { fmt.Println(a, b); return (b - a) }, 2)

	fmt.Println(f1f2.List())
	// Output: [0 1 4 9 16 25 36 49 64 81]

}

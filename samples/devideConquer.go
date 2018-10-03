// divideConquer.go
// map - filter - reduce
//

package main

import (
	"fmt"
	iter "github.com/AndreasBriese/itertools"
	"sync"
	"time"
)

const SAMPLELEN = 100000

func main() {
	// initialize int sequences
	l1 := make([]int, SAMPLELEN)
	copy(l1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29})
	copy(l1[30:], []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29})
	for i := 60; i < SAMPLELEN; i += 30 {
		copy(l1[i:], iter.ToIterInt(l1[i-30:i]).MapInto(func(x int) int { return x + 30 }).ToList())
	}
	l2 := make([]int, SAMPLELEN)
	copy(l2, l1)

	mapFn := func(i int) int {
		// something costly
		return (i + i<<1 + (i*i)<<1) / (1 + i>>1)
	}

	// map - filter - reduce with 4 goroutines
	wg := new(sync.WaitGroup)
	tms := make([]float64, 100)
	erg := 0
	for i := 0; i < 100; i++ {
		seq := iter.ToIterInt(l1)
		st := time.Now()
		erg = 0
		chF := make(chan *iter.IterableInt)
		chR := make(chan *iter.IterableInt)
		for _, seqU := range seq.Tee(4) {
			go func(seqU *iter.IterableInt) {
				chF <- seqU.MapInto(mapFn)
			}(seqU)
			go func(chf chan *iter.IterableInt) {
				ze := <-chf
				chR <- ze.Filter(func(i int) bool { return i&1 == 0 })
			}(chF)
			wg.Add(1)
			go func(chr chan *iter.IterableInt) {
				ze := <-chr
				erg += ze.Reduce(func(x, y int) int { return x + y })
				wg.Done()
			}(chR)
		}
		wg.Wait()
		tms[i] = float64(time.Since(st))
		copy(l1, l2) // revert changes
		close(chF)
		close(chR)
	}
	fmt.Println("map - filter - reduce with 4 goroutines",
		"\nresult:", erg,
		"\ntiming (100 rep):", iter.ToIterFloat64(tms).Reduce(func(x, y float64) float64 { return x + y })/100, "ns/op length", SAMPLELEN)

	// map - filter - reduce with chaining
	tms = make([]float64, 100)
	erg = 0
	for i := 0; i < 100; i++ {
		seq := iter.ToIterInt(l1)
		st := time.Now()
		erg = seq.Map(
			mapFn,
		).Filter(
			func(i int) bool { return i&1 == 0 },
		).Reduce(
			func(x, y int) int { return x + y })
		tms[i] = float64(time.Since(st))
		copy(l1, l2) // revert changes
	}
	fmt.Println("map - filter - reduce with chaining",
		"\nresult:", erg,
		"\ntiming (100 rep):", iter.ToIterFloat64(tms).Reduce(func(x, y float64) float64 { return x + y })/100, "ns/op length", SAMPLELEN)

}

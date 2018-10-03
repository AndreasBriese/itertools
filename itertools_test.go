// go package itertools
//
// The MIT License (MIT)
// Copyright (c) 2018 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package itertools

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

const CHARS = "!abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890#"

const SAMPLELEN = 1000000

var (
	prng         = rand.New(rand.NewSource(int64(12345))) //time.Now().Nanosecond())))
	l1, l2       []int
	li1, li2     []int64
	f1, f2       []float64
	s1, s2       []string
	lsum, fsum   = 0, 0.0
	lmean, fmean = 0, 0.0
	wg           = new(sync.WaitGroup)
)

func init() {
	// init testseq
	l1, l2 = make([]int, SAMPLELEN), make([]int, SAMPLELEN)
	li1, li2 = make([]int64, SAMPLELEN), make([]int64, SAMPLELEN)
	f1, f2 = make([]float64, SAMPLELEN), make([]float64, SAMPLELEN)
	s1, s2 = make([]string, SAMPLELEN), make([]string, SAMPLELEN)
	lsum, fsum = 0, 0.0 //int, float64
	for i := range l1 {
		r := prng.Intn(300)
		l1[i], l2[i] = r, r
		li1[i], li2[i] = int64(r), int64(r)
		lsum += r
		f := prng.Float64()
		f1[i], f2[i] = f, f
		fsum += f
		c := string(CHARS[prng.Intn(len(CHARS))]) + string(CHARS[prng.Intn(len(CHARS))]) + string(CHARS[prng.Intn(len(CHARS))])
		s1[i], s2[i] = c, c
	}
	lmean = lsum / SAMPLELEN
	fmean = fsum / SAMPLELEN
	// fmt.Println(l1, l2, s1, s2)
}

// itertools int
func TestSeqEqualIterInt(t *testing.T) {
	seq := ToIterInt(l1)

	seq.Destroy()
	if seq.Len != 0 {
		t.Errorf("%#v had not been destroyed", seq)
	}

	seq = ToIterInt(l1)

	if l2[len(l2)-1] != seq.Last() ||
		l1[0] != seq.First() {
		t.Errorf("element unequal %v", seq.This())
	}

	seq.Reset()
	for i, v := range l2 {
		if v != seq.Next() {
			t.Errorf("element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}

	seq.Reset()
	ll1 := seq.List()
	for i, v := range l2 {
		if v != ll1[i] {
			t.Errorf("element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}

	seq.Reset()
	ll1 = seq.List()
	for i, v := range l2 {
		if v != ll1[i] {
			t.Errorf("element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}

	// check behaviour overlength
	seq.Next() // ->iterLen
	seq.This()

	seq.ToEnd()
	for i := range l2 {
		if ii, v := len(l2)-1-i, l2[len(l2)-1-i]; v != seq.Back() {
			t.Errorf("element %3d unequal: is %v != should %v", ii, seq.This(), v)
		}
	}
	// check behaviour underlength
	seq.Back() // --> -1
	seq.This()

	circ := seq.Cycle()
	for i := 0; i < 200; i++ {
		if v := l2[i%len(l2)]; v != circ() {
			t.Errorf("Cycle: element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}

	i := 0
	for step, v, ex := seq.MapNext(func(elem int) int { return elem * 4 }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		if v2 := func(elem int) int { return elem * 4 }(l2[i]); v2 != v {
			t.Errorf("MapNext: element %3d unequal: is %v != should %v", i, seq.This(), v2)
		}
		i++
	}

	seq = ToIterInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	last := -999
	for step, v, ex := seq.DoubleCompNext(func(x, y int) bool { return x == y }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		if v == last {
			t.Errorf("DoubleCompNext != : element %3d equal %v ", v, last)
		}
		last = v
	}

	seq2 := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	i = 0
	for step, v, ex := seq2.FilterNext(func(elem int) bool { return elem < 5 }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		if v > 4 {
			t.Errorf("FilterNext: element %3d > 4: is %v ", i, seq2.This())
		}
		i++
	}

	seq = ToIterInt(l1)
	if e := seq.Reduce(func(x, y int) int { return x + y }); e != lsum {
		t.Errorf("Reduce: sum is %v ; should be %v", e, lsum)
	}
}

func XX_ExampleIterIntTeeGoroutines() {
	tms := make([]float64, 100)
	for i := 0; i < 100; i++ {
		seq := ToIterInt(l1)
		st := time.Now()
		e := seq.MapInto(func(i int) int { return 2 * i }).Filter(func(i int) bool { return i > 100 }).Reduce(func(x, y int) int { return x + y })
		// fmt.Println(time.Since(st), e)
		tms[i] = float64(time.Since(st))
		copy(l1, l2)
		e = e
	}
	fmt.Println(ToIterFloat64(tms).Reduce(func(x, y float64) float64 { return x + y }) / 100)

	tms = make([]float64, 100)
	for i := 0; i < 100; i++ {
		seq := ToIterInt(l1)
		st := time.Now()
		e := 0
		for _, seqU := range seq.Tee(8) {
			e += seqU.MapInto(func(i int) int { return 2 * i }).Filter(func(i int) bool { return i > 100 }).Reduce(func(x, y int) int { return x + y })
		}
		e = e
		tms[i] = float64(time.Since(st))
		// fmt.Println(time.Since(st), e)
		copy(l1, l2)
	}
	fmt.Println(ToIterFloat64(tms).Reduce(func(x, y float64) float64 { return x + y }) / 100)

	tms = make([]float64, 100)
	for i := 0; i < 100; i++ {

		seq := ToIterInt(l1)
		st := time.Now()
		e := 0
		chF := make(chan *IterableInt)
		chR := make(chan *IterableInt)
		defer close(chF)
		defer close(chR)

		for _, seqU := range seq.Tee(6) {
			go func(seqU *IterableInt) {
				chF <- seqU.MapInto(func(i int) int { return 2 * i }) // .Filter(func(i int) bool { return i > 100 })
			}(seqU)
			go func() {
				ze := <-chF
				chR <- ze.Filter(func(i int) bool { return i > 100 })
			}()
			wg.Add(1)
			go func() {
				ze := <-chR
				e += ze.Reduce(func(x, y int) int { return x + y })
				wg.Done()
			}()
		}
		wg.Wait()
		// fmt.Println(time.Since(st), e)
		tms[i] = float64(time.Since(st))
		copy(l1, l2)
	}
	fmt.Println(ToIterFloat64(tms).Reduce(func(x, y float64) float64 { return x + y }) / 100)
	// Output: 0
}

func ExampleDoubleComp() {
	seq := ToIterInt([]int{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	fmt.Println(seq.DoubleComp(func(x, y int) bool { return x != y }).List())
	// Output: [1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0]
}

func ExampleDoubleCompNext() {
	seq := ToIterInt([]int{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	for step, v, ex := seq.DoubleCompNext(func(x, y int) bool { return x != y }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	// Output: 1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0
}

func Example_Int_Map_vs_MapInto() {
	l := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	seq := ToIterInt(l)
	fmt.Printf("Result: %v slice l: %v\n", seq.Map(func(x int) int { return 2 * x }).List(), l)
	fmt.Printf("Result: %v slice l: %v\n", seq.MapInto(func(x int) int { return 2 * x }).List(), l)
	// Output: Result: [0 2 4 6 8 10 12 14 16 18] slice l: [0 1 2 3 4 5 6 7 8 9]
	// Result: [0 2 4 6 8 10 12 14 16 18] slice l: [0 2 4 6 8 10 12 14 16 18]
}

func ExampleDoubleOp() {
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(seq.DoubleOp(func(x, y int) int { return x + y }).List())
	// Output: [1 3 5 7 9 11 13 15 17]
}

func ExampleDoubleOpNext() {
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	for step, v, ex := seq.DoubleOpNext(func(x, y int) int { return x + y }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println(seq.List())
	// Output: 1 3 5 7 9 11 13 15 17 [0 1 2 3 4 5 6 7 8 9]
}

func ExampleDoubleCompF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	fmt.Println(seq.DoubleComp(func(x, y float64) bool { return x != y }).List())
	// Output: [1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0]
}

func ExampleDoubleCompNextF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	for step, v, ex := seq.DoubleCompNext(func(x, y float64) bool { return x != y }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	// Output: 1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0
}

func ExamplePairOp_Float64() {
	// initialize
	f1 := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f1seq := ToIterFloat64(f1).Map(func(x float64) float64 { return x * x })
	f2seq := ToIterFloat64(f1).Map(func(x float64) float64 { return 2 * x * x })

	f1f2 := ZipToIterFloat64(f1seq.List(), f2seq.List())
	f1f2 = f1f2.PairOp(func(a, b float64) float64 { return (b - a) }, 2) // fmt.Println(a, b);

	fmt.Println(f1f2.List())
	// Output: [0 1 4 9 16 25 36 49 64 81]
}

func ExamplePairOpF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(seq.PairOp(func(x, y float64) float64 { return x + y }).List())
	// Output: [1 5 9 13 17]
}

func ExamplePairOpNextF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	for step, v, ex := seq.PairOpNext(func(x, y float64) float64 { return x + y }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println(seq.List())
	// Output: 1 5 9 13 17 [0 1 2 3 4 5 6 7 8 9]
}

func ExampleDoubleOpF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(seq.DoubleOp(func(x, y float64) float64 { return x + y }).List())
	// Output: [1 3 5 7 9 11 13 15 17]
}

func ExampleDoubleOpNextF64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	for step, v, ex := seq.DoubleOpNext(func(x, y float64) float64 { return x + y }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println(seq.List())
	// Output: 1 3 5 7 9 11 13 15 17 [0 1 2 3 4 5 6 7 8 9]
}

func ExampleDoubleCompIf() {
	seq := ToIterIf([]float64{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	fmt.Println(seq.DoubleComp(func(x, y interface{}) bool { return x.(float64) != y.(float64) }).List())
	// Output: [1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0]
}

func ExampleDoubleCompNextIf() {
	seq := ToIterIf([]float64{0, 1, 2, 2, 3, 4, 4, 5, 6, 7, 8, 9, 9, 8, 8, 7, 6, 6, 5, 5, 4, 4, 3, 2, 2, 2, 1, 0})
	for step, v, ex := seq.DoubleCompNext(func(x, y interface{}) bool { return x.(float64) != y.(float64) }), interface{}(0.0), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	// Output: 1 2 3 4 5 6 7 8 9 8 7 6 5 4 3 2 1 0
}

func ExampleDoubleOpIf() {
	seq := ToIterIf([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(seq.DoubleOp(func(x, y interface{}) interface{} { return x.(float64) + y.(float64) }).List())
	// Output: [1 3 5 7 9 11 13 15 17]
}

func ExampleDoubleOpNextIf() {
	seq := ToIterIf([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	for step, v, ex := seq.DoubleOpNext(func(x, y interface{}) interface{} { return x.(float64) + y.(float64) }), interface{}(0.0), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println(seq.List())
	// Output: 1 3 5 7 9 11 13 15 17 [0 1 2 3 4 5 6 7 8 9]
}

func ExampleIterChain() {
	seqI := ChainToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	seqF := ChainToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	// seqS := ChainToIterIf([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"A", "B", "C", "D", "E", "F", "G"})
	fmt.Println(seqI.List())
	fmt.Println(seqF.List())
	// fmt.Println(seqS.List())
	// Output: [1 3 5 7 9 11 13 15 17]
}

func ExampleIterIntTee() {
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	seqTeed2 := seq.Tee(2)
	fmt.Print(seqTeed2[0].List())
	fmt.Println(seqTeed2[1].List())
	seqTeed2 = seq.Tee(3)
	fmt.Print(seqTeed2[0].List())
	fmt.Print(seqTeed2[1].List())
	fmt.Println(seqTeed2[2].List())
	// Output: [0 1 2 3 4 5 6 7 8 9][9 8 7 6 5 4 3 2 1 0]
	// [0 1 2 3 4 5 6][7 8 9 9 8 7 6][5 4 3 2 1 0]
}

func ExampleIterIntFilterNext() {
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	for step, v, ex := seq.FilterNext(func(elem int) bool { return elem < 5 }), 0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v) // should print 0 1 2 3 4 4 3 2 1 0
	}
	seq1 := seq.Filter(func(elem int) bool { return elem < 5 })
	for i := 0; i < seq1.Len; i++ {
		fmt.Printf("%v ", seq1.Next())
	}
	list := seq.Map(func(elem int) int { return elem * 2 }).List()
	fmt.Println(list)
	// Output: 0 1 2 3 4 4 3 2 1 0 0 1 2 3 4 4 3 2 1 0 [0 2 4 6 8 10 12 14 16 18 18 16 14 12 10 8 6 4 2 0]
}

func ExampleToListIterIntFilterNext() {
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	list := seq.Filter(func(elem int) bool { return elem < 5 }).List()
	fmt.Printf("%v ", list) // should print [0 1 2 3 4 4 3 2 1 0]
	seq = ToIterInt(list)
	list2 := seq.Map(func(elem int) int { return elem * 4 }).List()
	fmt.Println(list2)
	// Output: [0 1 2 3 4 4 3 2 1 0] [0 4 8 12 16 16 12 8 4 0]
}

func ExampleToListIterInt_MapInto_FilterNext_Reduce() {
	// make iterable
	seq := ToIterInt([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	// MapNext
	seqFE := seq.MapInto(func(i int) int { return i * 2 })
	fmt.Println("MapNext: '*2' :",
		seqFE.List())

	// FilterNext
	seqF := seqFE.Filter(func(i int) bool { return i > 5 })
	fmt.Println("FilterNext '>5' :",
		seqF.List())

	// Reduce
	e := seqF.Reduce(func(x, y int) int { return x + y })
	fmt.Println("Reduce sum: ", e)

	// Output:
	// MapNext: '*2' : [0 2 4 6 8 10 12 14 16 18]
	// FilterNext '>5' : [6 8 10 12 14 16 18]
	// Reduce sum:  84
}

// itertools float64
func TestSeqEqualFloat64(t *testing.T) {
	seq := ToIterFloat64(f1)

	if f2[len(f2)-1] != seq.Last() ||
		f1[0] != seq.First() {
		t.Errorf("element unequal %v", seq.This())
	}

	seq.Reset()
	for i, v := range f2 {
		if v != seq.Next() {
			t.Errorf("element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}
	// check behaviour overlength
	seq.Next() // ->iterLen
	seq.This()

	seq.ToEnd()
	for i := range f2 {
		if ii, v := len(f2)-1-i, f2[len(f2)-1-i]; v != seq.Back() {
			t.Errorf("element %3d unequal: is %v != should %v", ii, seq.This(), v)
		}
	}
	// check behaviour underlength
	seq.Back() // --> -1
	seq.This()

	circ := seq.Cycle()
	for i := 0; i < 200; i++ {
		if v := f2[i%len(f2)]; v != circ() {
			t.Errorf("Cycle: element %3d unequal: is %v != should %v", i, seq.This(), v)
		}
	}

	seq.Reset()
	//	st := time.Now()
	s := seq.Reduce(func(x, y float64) float64 { return x + y })
	//	fmt.Println(time.Since(st))

	if s != fsum {
		t.Errorf("Reduce: %v should be %v", s, lsum)
	}

	i := 0
	for step, v, ex := seq.MapNext(func(elem float64) float64 { return elem * 4 }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		if v2 := func(elem float64) float64 { return elem * 4 }(f2[i]); v2 != v {
			t.Errorf("MapNext: element %3d unequal: is %v != should %v", i, seq.This(), v2)
		}
		i++
	}

	seq2 := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	i = 0
	for step, v, ex := seq2.FilterNext(func(elem float64) bool { return elem < 5 }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		if v > 4 {
			t.Errorf("FilterNext: element %3d > 4: is %v ", i, seq2.This())
		}
		i++
	}

	seq = ToIterFloat64(f1)
	sum := 0.0
	for _, v := range f1 {
		sum += v
	}
	if e := seq.Reduce(func(x, y float64) float64 { return x + y }); e != sum {
		t.Errorf("Reduce: sum is %v ; should be %v", e, sum)
	}
}

func ExampleIterFilterNextFloat64() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	for step, v, ex := seq.FilterNext(func(elem float64) bool { return elem < 5 }), 0.0, false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v) // should print 0 1 2 3 4 4 3 2 1 0
	}
	seq1 := seq.Filter(func(elem float64) bool { return elem < 5 })
	for i := 0; i < seq1.Len; i++ {
		fmt.Printf("%v ", seq1.Next())
	}
	// Output: 0 1 2 3 4 4 3 2 1 0 0 1 2 3 4 4 3 2 1 0
}

func ExampleToListIterFloat64FilterNext() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	list := seq.Filter(func(elem float64) bool { return elem < 5 }).List()
	fmt.Printf("%v ", list) // should print [0 1 2 3 4 4 3 2 1 0]
	seq = ToIterFloat64(list)
	list2 := seq.Map(func(elem float64) float64 { return elem * 4 }).List()
	fmt.Println(list2)
	// Output: [0 1 2 3 4 4 3 2 1 0] [0 4 8 12 16 16 12 8 4 0]
}

func ExampleIterFloat64Tee() {
	seq := ToIterFloat64([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	seqTeed2 := seq.Tee(2)
	fmt.Print(seqTeed2[0].List())
	fmt.Println(seqTeed2[1].List())
	seqTeed2 = seq.Tee(3)
	fmt.Print(seqTeed2[0].List())
	fmt.Print(seqTeed2[1].List())
	fmt.Println(seqTeed2[2].List())
	// Output: [0 1 2 3 4 5 6 7 8 9][9 8 7 6 5 4 3 2 1 0]
	// [0 1 2 3 4 5 6][7 8 9 9 8 7 6][5 4 3 2 1 0]
}

// itertools interface
func TestSeqEqualIf(t *testing.T) {
	seqIf := ToIterIf(s1)

	if s2[len(s2)-1] != seqIf.Last() ||
		s1[0] != seqIf.First() {
		t.Errorf("element unequal %v", seqIf.This())
	}

	seqIf.Reset()
	for i, v := range s2 {
		if v != seqIf.Next() {
			t.Errorf("element %3d unequal: is %v != should %v", i, seqIf.This(), v)
		}
	}
	// check behaviour overlength
	seqIf.Next() // ->iterLen
	seqIf.This()

	seqIf.ToEnd()
	for i := range s2 {
		if ii, v := len(s2)-1-i, s2[len(s2)-1-i]; v != seqIf.Back() {
			t.Errorf("element %3d unequal: is %v != should %v", ii, seqIf.This(), v)
		}
	}
	// check behaviour underlength
	seqIf.Back() // --> -1
	seqIf.This()

	circIf := seqIf.Cycle()
	for i := 0; i < 200; i++ {
		if v := s2[i%len(s2)]; v != circIf() {
			t.Errorf("Cycle: element %3d unequal: is %v != should %v", i, seqIf.This(), v)
		}
	}

	mut := func(elem interface{}) interface{} {
		return fmt.Sprintf("--> %v ", elem)
	}
	i := 0
	for it, v, ex := seqIf.MapNext(mut), interface{}(""), false; ; {
		v, ex = it()
		if ex {
			break
		}
		if v2 := mut(s2[i]); v2 != v {
			t.Errorf("MapNext: element %3d unequal: is %v != should %v", i, seqIf.This(), v2)
		}
		i++
	}

	seqIf2 := ToIterIf([]string{"a", "b", "c", "A", "B", "C", "d", "D"})
	condFn := func(elem interface{}) bool {
		return elem != strings.ToUpper(elem.(string))
	}
	i = 0
	for step, v, ex := seqIf2.FilterNext(condFn), interface{}(""), false; ; {
		v, ex = step()
		if ex {
			break
		}
		if !condFn(v) {
			t.Errorf("FilterNext: element %3d > 4: is %v ", i, seqIf2.This())
		}
		i++
	}

}

func ExampleIterInt_Chaining() {
	seq := ChainToIterInt([]int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9})
	fmt.Println(seq.List())
	seq2 := ChainToIterInt([]int{1, 2, 3}, []int{4, 5, 6})
	fmt.Println(ChainIterInt(seq, seq2).List())
	// Output: [1 2 3 4 5 6 7 8 9]
	// [1 2 3 4 5 6 7 8 9 1 2 3 4 5 6]
}

func ExampleIter_All_Any_Where() {
	seq := ToIterFloat64([]float64{1, 2, 3})
	if seq.Any(3) {
		fmt.Printf("equal 3 at index: %v\n", seq.Index())
	}
	if !seq.All(1) {
		fmt.Printf("%v at index: %v is not equal 1\n", seq.This(), seq.Index())
	}
	fmt.Println(seq.Any(4))
	fmt.Printf("Condition 2x<5 is met at indices %v\n\n", seq.Where(func(x float64) bool { return x*2 < 5 }).List())
	seqIf := ToIterIf([]float64{1, 2, 3})
	if seqIf.Any(float64(3)) {
		fmt.Printf("equal 3 at index: %v\n", seqIf.Index())
	}
	if !seqIf.All(float64(1)) {
		fmt.Printf("%v at index: %v is not equal 1\n", seqIf.This(), seqIf.Index())
	}
	fmt.Printf("Condition 2x<5 is met at indices %v", seqIf.Where(func(x interface{}) bool { return x.(float64)*2 < 5 }).List())
	// Output: equal 3 at index: 2
	// 2 at index: 1 is not equal 1
	// false
	// Condition 2x<5 is met at indices [0 1]
	//
	// equal 3 at index: 2
	// 2 at index: 1 is not equal 1
	// Condition 2x<5 is met at indices [0 1]
}

func ExampleIterIfFilterNext() {
	seq := ToIterIf([]string{"a", "b", "c", "A", "B", "C", "d", "D"})
	condFn := func(elem interface{}) bool {
		return elem == strings.ToLower(elem.(string))
	}
	for step, v, ex := seq.FilterNext(condFn), interface{}(""), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v) // should print a b c d
	}
	seq1 := seq.Filter(condFn)
	for i := 0; i < seq1.Len; i++ {
		fmt.Printf("%v ", seq1.Next())
	}
	// Output: a b c d a b c d
}

func ExampleZipToIterIf() {
	zipIf := ZipToIterIf(
		[]string{"a", "b", "c", "d", "e"},
		[]string{"A", "B", "C", "D", "E"},
	)
	fn := func(elem interface{}) interface{} {
		return elem == strings.ToLower(elem.(string))
	}
	for step, v, ex := zipIf.MapNext(fn), interface{}(false), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v) // should print true false true false true false true false true false
	}
	fmt.Println()
	// Output: true false true false true false true false true false
}

func ExampleZipToIterFloat64() {
	zip := ZipToIterFloat64(
		[]float64{0, 1, 2, 3, 4, 5},
		[]float64{0, 1, 4, 9, 16, 25},
	)
	fn := func(a, b float64) float64 {
		return a*a - b
	}
	fmt.Println(zip.PairOp(fn).List())
	// Output: [0 0 0 0 0 0]
}

func ExampleReduceIterInt() {
	seq := ToIterInt([]int{1, 2, 3, 4})
	fmt.Println(seq.Reduce(func(x, y int) int { return x + y }))
	// Output: 10
}

func ExampleMMapIterInt() {
	mapFn := func(iterVal []int) int {
		return iterVal[0] + iterVal[1] + iterVal[2]
	}
	mapp3 := MMapIterInt(mapFn, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, []int{1, 2, 3, 4})
	for i := 0; i < len([]int{1, 2, 3, 4}); i++ {
		fmt.Printf("%v:%v ", i, mapp3())
	}
	// Output: 0:3 1:6 2:9 3:12
}

func ExampleMMapToIterInt() {
	mapFn := func(iterVal []int) int {
		return iterVal[0] + iterVal[1] + iterVal[2]
	}
	mapp3 := MMapToIterInt(mapFn, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}).Map(func(x int) int { return 2 * x })
	for i := 0; i < len([]int{1, 2, 3, 4}); i++ {
		fmt.Printf("%v:%v ", i, mapp3.Next())
	}
	// Output: 0:6 1:12 2:18 3:24
}

func ExampleMMapIterIf() {
	mapFn := func(iterVal []interface{}) interface{} {
		return iterVal[0].(uint8) + iterVal[1].(uint8) + iterVal[2].(uint8)
	}
	ls := []interface{}{uint8(1), uint8(2), uint8(3), uint8(4)}
	mapsum3 := MMapIterIf(mapFn, ls, ls, ls)
	for i := 0; i < len(ls); i++ {
		fmt.Printf("%v:%v ", i, mapsum3())
	}
	// Output: 0:3 1:6 2:9 3:12
}

func ExampleMapIfMixed() {
	mapFn := func(iterVal []interface{}) interface{} {
		return int(iterVal[0].(uint8)) + int(iterVal[1].(int)) + int(iterVal[2].(uint8))
	}
	ls := []interface{}{uint8(1), uint8(2), uint8(3), uint8(4)}
	mapsum3 := MMapIterIf(mapFn, ls, []interface{}{1, 2, 3, 4}, ls)
	for i := 0; i < len(ls); i++ {
		fmt.Printf("%v:%v ", i, mapsum3().(int))
	}
	// Output: 0:3 1:6 2:9 3:12
}

func ExampleToListIterIfFilterNext() {
	seq := ToIterIf([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	list := seq.Filter(func(elem interface{}) bool { return elem.(float64) < 5 }).List()
	fmt.Printf("%v ", list) // should print [0 1 2 3 4 4 3 2 1 0]
	seq = ToIterIf(list)
	list2 := seq.Map(func(elem interface{}) interface{} { return elem.(float64) * 4 }).List()
	fmt.Println(list2)
	// Output: [0 1 2 3 4 4 3 2 1 0] [0 4 8 12 16 16 12 8 4 0]
}

func ExampleIterIfTee() {
	seq := ToIterIf([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	seqTeed2 := seq.Tee(2)
	fmt.Print(seqTeed2[0].List())
	fmt.Println(seqTeed2[1].List())
	seqTeed2 = seq.Tee(3)
	fmt.Print(seqTeed2[0].List())
	fmt.Print(seqTeed2[1].List())
	fmt.Println(seqTeed2[2].List())
	// Output: [0 1 2 3 4 5 6 7 8 9][9 8 7 6 5 4 3 2 1 0]
	// [0 1 2 3 4 5 6][7 8 9 9 8 7 6][5 4 3 2 1 0]
}

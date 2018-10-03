//iteratools_BM_test.go

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
	// "math/rand"
	"strings"
	"testing"
	// "time"
	// "github.com/yanatan16/itertools"
)

// var (
// 	prng         = rand.New(rand.NewSource(int64(12345))) //time.Now().Nanosecond())))
// 	l1, l2       []int
// 	f1, f2       []float64
// 	s1, s2       []string
// 	lsum, fsum   = 0, 0.0
// 	lmean, fmean = 0, 0.0
// )

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

/*

--------- Benchmarking

*/

func Benchmark_IterInt_ToIterInt(b *testing.B) {
	// b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq := ToIterInt(l1)
		seq = seq
	}
}

/*
func Benchmark_Int_SumLoopOverIntSlice(b *testing.B) {
	// b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		for i := range l1 {
			sum += l1[i]
		}
		b.StopTimer()
		if sum != lsum {
			fmt.Println(lsum, sum)
		}
	}
}

func Benchmark_IterInt_SumNextOverIterable(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		seq.Reset()
		for i := 0; i < SAMPLELEN; i++ {
			sum += seq.Next()
		}
		b.StopTimer()
		if sum != lsum {
			fmt.Println(lsum, sum)
		}
	}
}

func Benchmark_Yanatan_SumNextOverIterable(b *testing.B) {
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seqCh := itertools.Int64(li1...)
		sum := int64(0)
		b.StartTimer()
		for i := 0; i < SAMPLELEN; i++ {
			sum += (<-seqCh).(int64)
		}
		b.StopTimer()
		if sum != int64(lsum) {
			fmt.Println(lsum, sum)
		}
	}
}


func Benchmark_Int_MapLoopOverIntSlice(b *testing.B) {
	fn := func(x int) int { return 2 * x }
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		for i := range l1 {
			sum += fn(l1[i])
		}
		b.StopTimer()
		if sum != 2*lsum {
			fmt.Println(2*lsum, sum)
		}
	}
}

func Benchmark_IterInt_MapOverIterable(b *testing.B) {
	fn := func(x int) int { return 2 * x }
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		ssq := seq.Map(fn)
		b.StopTimer()
		for _, v := range ssq.List() {
			sum += v
		}
		if sum != 2*lsum {
			fmt.Println(2*lsum, sum)
		}
	}
}

func Benchmark_IterInt_MapIntoOverIterable(b *testing.B) {
	fn := func(x int) int { return 2 * x }
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		ssq := seq.MapInto(fn)
		b.StopTimer()
		for _, v := range ssq.List() {
			sum += v
		}
		if sum != 2*lsum {
			fmt.Println(2*lsum, sum)
		}
		copy(l1, l2)
	}

}

func Benchmark_IterInt_MapNextOverIterable(b *testing.B) {
	fn := func(x int) int { return 2 * x }
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		sum := 0
		b.StartTimer()
		seq.Reset()
		for i, mappr := 0, seq.MapNext(fn); i < SAMPLELEN; i++ {
			v, _ := mappr()
			sum += v
		}
		b.StopTimer()
		if sum != 2*lsum {
			fmt.Println(2*lsum, sum)
		}
	}
}

*/
func Benchmark_IterInt_ToList(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		l := seq.ToList()
		l = l
	}
}

func Benchmark_IterInt_List(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		l := seq.List()
		l = l
	}
}

/*
func Benchmark_Yanatan_MapNextOverIterable(b *testing.B) {
	fn := func(x interface{}) interface{} { return int64(2) * x.(int64) }
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seqCh := itertools.Int64(li1...)
		mppd := itertools.Map(fn, seqCh)
		sum := int64(0)
		b.StartTimer()
		for i := 0; i < SAMPLELEN; i++ {
			sum += (<-mppd).(int64)
		}
		b.StopTimer()
		if sum != int64(2*lsum) {
			fmt.Println(lsum, sum)
		}
	}
}
*/

func Benchmark_IterInt_Map(b *testing.B) {
	fn := func(elem int) int {
		return 2 * elem
	}
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Map(fn)
		v = v
	}
}

func Benchmark_IterInt_MapInto(b *testing.B) {
	fn := func(elem int) int {
		return 2 * elem
	}
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		v := seq.MapInto(fn)
		v = v
		b.StopTimer()
		// restore l1
		copy(l1, l2)
	}
}

func Benchmark_IterInt_MapNext(b *testing.B) {
	fn := func(elem int) int {
		return 10.0 * elem
	}
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.MapNext(fn), 0, false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterInt_Filter(b *testing.B) {
	seqIf := ToIterInt(l1)
	condFn := func(elem int) bool {
		return elem < lmean
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seqIf.Filter(condFn)
		v = v
	}
}

func Benchmark_IterInt_FilterNext(b *testing.B) {
	fn := func(elem int) bool {
		return elem < lmean
	}
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.FilterNext(fn), 0, false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterInt_Reduce(b *testing.B) {
	// make iterable
	seq := ToIterInt(l1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Reduce(func(x, y int) int { return x + y })
		v = v
	}
}

func Benchmark_IterInt_Map_Filter_Reduce(b *testing.B) {
	seq := ToIterInt(l1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Map(func(i int) int { return i * 2 }).Filter(func(i int) bool { return i > lmean }).Reduce(func(x, y int) int { return x + y })
		v = v
	}
	// restore l1
	copy(l1, l2)
}

func Benchmark_IterInt_MapInto_Filter_Reduce(b *testing.B) {
	seq := ToIterInt(l1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		v := seq.MapInto(func(i int) int { return i * 2 }).Filter(func(i int) bool { return i > lmean }).Reduce(func(x, y int) int { return x + y })
		v = v
		b.StopTimer()
		// restore l1
		copy(l1, l2)
	}
}

func Benchmark_IterInt_PairOp(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.PairOp(func(a, b int) int { return 2*a + b })
		t5 = t5
	}
}

func Benchmark_IterInt_DoubleOp(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.DoubleOp(func(p, a int) int { return p + a })
		t5 = t5
	}
}

func Benchmark_IterInt_DoubleComp(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.DoubleComp(func(p, a int) bool { return p < a })
		t5 = t5
	}
}

func Benchmark_IterInt_Tee5(b *testing.B) {
	seq := ToIterInt(l1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.Tee(5)
		t5 = t5
	}
}

// func Benchmark_IterInt_MMap(b *testing.B) {
// 	mapFn := func(iterVal []int) int {
// 		return iterVal[0] + iterVal[1] - iterVal[2]
// 	}
// 	mapp3 := MMapIterInt(mapFn,
// 		l1,
// 		l2,
// 		l1,
// 	)
// 	for r := 0; r < b.N; r++ {
// 		for i := 0; i < len(l2); i++ {
// 			mapp3()
// 		}
// 	}
// }

// IterableFloat64

func Benchmark_IterFloat64_ToIterFloat64(b *testing.B) {
	// b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq := ToIterFloat64(f1)
		seq = seq
	}
}

func Benchmark_IterFloat64_NextOverIterable(b *testing.B) {
	seq := ToIterFloat64(f1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq.Reset()
		for i, v := range f2 {
			if v != seq.Next() {
				fmt.Printf("element %3d unequal: is %v != should %v", i, seq.This(), v)
			}
		}
	}
}

func Benchmark_IterFloat64_Map(b *testing.B) {
	seq := ToIterFloat64(f1) //[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		s := seq.Map(func(elem float64) float64 { return elem / 10 })
		s = s
	}
}

func Benchmark_IterFloat64_MapInto(b *testing.B) {
	fn := func(elem float64) float64 {
		return 2 * elem
	}
	seq := ToIterFloat64(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		v := seq.MapInto(fn)
		v = v
		b.StopTimer()
		// restore f1
		copy(f1, f2)
	}

}

func Benchmark_IterFloat64_MapNext(b *testing.B) {
	fn := func(elem float64) float64 {
		return 10.0 * elem
	}
	seq := ToIterFloat64(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.MapNext(fn), 0.0, false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterFloat64_Filter(b *testing.B) {
	seq := ToIterFloat64(f1) // []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		s := seq.Filter(func(elem float64) bool { return elem < fmean })
		s = s
	}
}

func Benchmark_IterFloat64_FilterNext(b *testing.B) {
	fn := func(elem float64) bool {
		return elem < fmean
	}
	seq := ToIterFloat64(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.FilterNext(fn), 0.0, false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterFloat64_Map_ToList(b *testing.B) {
	seq := ToIterFloat64(f1) // []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		list := seq.Map(func(elem float64) float64 { return elem * 2 }).List()
		list = list
	}
}

func Benchmark_IterFloat64_Reduce(b *testing.B) {
	seq := ToIterFloat64(f1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Reduce(func(x, y float64) float64 { return x + y })
		v = v
	}
}

func Benchmark_IterFloat64_Map_Filter_Reduce(b *testing.B) {
	seq := ToIterFloat64(f1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq.Map(func(i float64) float64 { return i * 2 }).Filter(func(i float64) bool { return i > fmean }).Reduce(func(x, y float64) float64 { return x + y })
	}
}

func Benchmark_IterFloat64_MapInto_Filter_Reduce(b *testing.B) {
	seq := ToIterFloat64(f1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		s := seq.MapInto(func(i float64) float64 { return i * 2 }).Filter(func(i float64) bool { return i > fmean }).Reduce(func(x, y float64) float64 { return x + y })
		s = s
		b.StopTimer()
		// restore f1
		copy(f1, f2)
	}
}

func Benchmark_IterFloat64_Tee5(b *testing.B) {
	seq := ToIterFloat64(f1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.Tee(5)
		t5 = t5
	}
}

// func Benchmark_IterFloat64_MMap(b *testing.B) {
//     mapFn := func(iterVal []float64) float64 {
//         return iterVal[0] + iterVal[1] - iterVal[2]
//     }
//     mapp3 := MMapIterFloat64(mapFn,
//         f1,
//         f2,
//         f1,
//     )
//     for r := 0; r < b.N; r++ {

//         for i := 0; i < len(f2); i++ {
//             mapp3()
//         }
//     }
// }

// Iterable Interface

func Benchmark_IterIf_ToIterIf(b *testing.B) {
	// b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq := ToIterIf(f1)
		seq = seq
	}
}

func Benchmark_IterIf_NextOverIterable(b *testing.B) {
	seq := ToIterIf(f1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		seq.Reset()
		for i := 0; i < SAMPLELEN; i++ {
			seq.Next()
		}
	}
}

func Benchmark_IterIf_Map(b *testing.B) {
	fn := func(elem interface{}) interface{} {
		return 10.0 * elem.(float64)
	}
	seq := ToIterIf(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	for r := 0; r < b.N; r++ {
		v := seq.Map(fn)
		v = v
	}
}

func Benchmark_IterIf_MapInto(b *testing.B) {
	fn := func(elem interface{}) interface{} {
		return 2 * elem.(float64)
	}
	seq := ToIterIf(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		v := seq.MapInto(fn)
		v = v
		b.StopTimer()
		// restore f1
		copy(f1, f2)
	}
}

func Benchmark_IterIf_MapNext(b *testing.B) {
	fn := func(elem interface{}) interface{} {
		return 10.0 * elem.(float64)
	}
	seq := ToIterIf(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.MapNext(fn), interface{}(0.0), false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterIf_Filter(b *testing.B) {
	seqIf := ToIterIf(s1) // []string{"a", "b", "c", "A", "B", "C", "d", "D"})
	condFn := func(elem interface{}) bool {
		return string(elem.(string)[0]) == strings.ToLower(string(elem.(string)[0]))
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seqIf.Filter(condFn)
		v = v
	}
}

func Benchmark_IterIf_FilterNext(b *testing.B) {
	fn := func(elem interface{}) bool {
		return elem.(float64) < fmean
	}
	seq := ToIterIf(f1) //[]float64{.1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0},
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for step, v, ex := seq.FilterNext(fn), interface{}(0.0), false; ; {
			v, ex = step()
			if ex {
				break
			}
			v = v
		}
	}
}

func Benchmark_IterIf_Reduce(b *testing.B) {
	// make iterable
	seq := ToIterIf(l1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Reduce(func(x, y interface{}) interface{} { return x.(int) + y.(int) })
		v = v
	}
}

func Benchmark_IterIf_Map_Filter_Reduce(b *testing.B) {
	seq := ToIterIf(f1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		v := seq.Map(func(i interface{}) interface{} { return i.(float64) * 2 }).Filter(func(i interface{}) bool { return i.(float64) > fmean }).Reduce(func(x, y interface{}) interface{} { return x.(float64) + y.(float64) })
		v = v
	}
}

func Benchmark_IterIf_MapInto_Filter_Reduce(b *testing.B) {
	seq := ToIterIf(f1) // []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	for r := 0; r < b.N; r++ {
		b.StartTimer()
		v := seq.MapInto(func(i interface{}) interface{} { return i.(float64) * 2 }).Filter(func(i interface{}) bool { return i.(float64) > fmean }).Reduce(func(x, y interface{}) interface{} { return x.(float64) + y.(float64) })
		v = v
		b.StopTimer()
		// restore f1
		copy(f1, f2)
	}
}

func Benchmark_IterIf_Tee5(b *testing.B) {
	seq := ToIterIf(f1)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		t5 := seq.Tee(5)
		t5 = t5
	}
}

// func Benchmark_IterIf_MMap(b *testing.B) {
//     mapFn := func(iterVal []interface{}) interface{} {
//         return iterVal[0].(float64) + iterVal[1].(float64) - iterVal[2].(float64)
//     }
//     mapp3 := MMapIterIf(mapFn,
//         f1,
//         f2,
//         f1,
//     )
//     for r := 0; r < b.N; r++ {
//         for i := 0; i < len(f2); i++ {
//             mapp3()
//         }
//     }
// }

// func Benchmark_IterIf_Zip_and_MapNext(b *testing.B) {
// 	fn := func(elem interface{}) interface{} {
// 		return elem == strings.ToLower(elem.(string))
// 	}
// 	for r := 0; r < b.N; r++ {
// 		zipIf := ZipIterIf(
// 			s1, // []string{"a", "b", "c", "d", "e"},
// 			s2, // []string{"A", "B", "C", "D", "E"},
// 		)
// 		v := zipIf.Map(fn)
// 		v = v
// 	}
// }

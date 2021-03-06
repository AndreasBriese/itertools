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

// Code generated by github.com/AndreasBriese/itertools; DO NOT EDIT.
// File from command $go generate executed by go:generate in makeMoreItertools.go 
// Timestamp: 2018-09-21 17:31:33.635372325 +0200 CEST m=+0.002208485
// 

package itertools

// Type IterableInt16 an iterable over a slice of any type
// Use for any type of iterables that are CAN NOT use []int or []int16, i.e. string.
// IterableIf needs more memory
//      (because it likely will replicate the underlying slice
// 		and Next element in []iterface{} has one memoryword extra used for the elements type)
type IterableInt16 struct {
	// stepwise returns / does not destroy original underlying slice / no additional memory
	Reset, ToEnd                  func()
	Next, This, Back, First, Last func() int16
	Cycle                         func() func() int16
	MapNext                       func(func(int16) int16) func() (int16, bool)
	FilterNext                    func(func(int16) bool) func() (int16, bool)
	DoubleOpNext                  func(func(int16, int16) int16) func() (int16, bool)
	DoubleCompNext                func(func(int16, int16) bool) func() (int16, bool)
	PairOpNext                    func(func(int16, int16) int16, ...int) func() (int16, bool)

	// info / does not destroy original underlying slice
	Len      int
	Index    func() int
	SetIndex func(int) (int, bool)
	Any, All func(int16) bool
	Where    func(func(int16) bool) *IterableInt

	// conversions & abstractions / do not destroy the iterable nor the original underlying slice
	List   func() []int16
	ToList func() []int16 // needs additional memory
	Reduce func(func(int16, int16) int16) int16
	Tee    func(int) []*IterableInt16

	// return of new iterable(s) / does not destroy or change original / needs additional memory
	DoubleOp   func(func(int16, int16) int16) *IterableInt16
	DoubleComp func(func(int16, int16) bool) *IterableInt16
	PairOp     func(func(int16, int16) int16, ...int) *IterableInt16
	Filter     func(func(int16) bool) *IterableInt16
	Map        func(func(int16) int16) *IterableInt16

	// Return the iterabel BUT changes the underlying slice!
	MapInto func(func(int16) int16) *IterableInt16

	// Replace the iterable's underlying slice by a slice []int16 with length 0
	Destroy func()
}

// ZipToIterInt16(l1, l2) *IterableInt16
// takes to slices and returns a iterator over the zipping result
// zipps two slices and creates a 2*length []interface{} slice from it
//   -- Same size slices
//      or first smaller than the second (second will be cut-off an the first length)
func ZipToIterInt16(l1, l2 []int16) *IterableInt16 {
	if len(l1) != len(l2) {
		panic(ERR_SHORTER2)
	}
	l1l2 := make([]int16, 0, len(l1)<<1)
	for i := range l1 {
		l1l2 = append(l1l2, l1[i], l2[i])
	}
	return ToIterInt16(l1l2)
}

// ChainToIterInt16(...lists) takes at least 2 slices and returns an newly created iterable over these.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainToIterInt16(lists ...[]int16) *IterableInt16 {
	if len(lists) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range lists {
		totalLen += len(lists[i])
	}
	chain := make([]int16, totalLen)
	idx := 0
	for i := range lists {
		copy(chain[idx:], lists[i])
		idx += len(lists[i])
	}
	return ToIterInt16(chain)
}

// ChainIterInt16(...iters) takes at least 2 iterables of the same type and returns an newly created iterable over concat.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainIterInt16(iters ...*IterableInt16) *IterableInt16 {
	if len(iters) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range iters {
		totalLen += iters[i].Len
	}
	chain := make([]int16, totalLen)
	idx := 0
	for i := range iters {
		copy(chain[idx:], iters[i].List())
		idx += iters[i].Len
	}
	return ToIterInt16(chain)
}

// MMapToIterInt16(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
// maps a function fn to all slices given as comma separated parameter
// and returns an iterator over the slice with the results.
// func fn works on a slice with the i'th element from the iter conversed seq and needs to
//      the slice has a length of the number of slices (0..n) passed to MapIterIf(fn, seqs...)
//      code should adress the source seq by its index 0..n
//      i.e. a function to add the values of three byte slices and return that sum
//			fn := func(v []int16) int16{} { return v[0] + v[1] + v[2] }
// 			sum := MMapIterIf(fn, slice1, slice2, slice3)
// 			fmt.Println( sum.Next() ) // prints sum of the first element of the three given slices
//  		fmt.Println( sum.Next() ) // prints sum of the second element of the three given slices
// 			...
// Attention!  All slices to be multi-mapped need to be of the same length
// Needs additional memory (size of one of the given slices)
func MMapToIterInt16(fn func([]int16) int16, seqs ...[]int16) *IterableInt16 {
	// checks
	if len(seqs) < 2 {
		panic(ERR_SHORTER2)
	}
	slen := len(seqs[0])
	for _, seq := range seqs {
		if len(seq) != slen {
			panic(ERR_DIFFLEN)
		}
	}

	newIter := make([]int16, slen)
	vals := make([]int16, len(seqs))
	for i := range newIter {
		for ii := range vals {
			vals[ii] = seqs[ii][i]
		}
		newIter[i] = fn(vals)
	}
	return ToIterInt16(newIter)
}

// MMapIterInt16(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
// maps a function fn to all slices given as comma separated parameter
// and yields the results.
// func fn works on a slice with the i'th element from the iter conversed seq and needs to
//      the slice has a length of the number of slices (0..n) passed to MapIterIf(fn, seqs...)
//      code should adress the source seq by its index 0..n
//      i.e. a function to add the values of three byte slices and return that sum
//			fn := func(v []float32) float32{} { return v[0] + v[1] + v[2] }
// 			sum := MMapIterIf(fn, slice1, slice2, slice3)
// 			fmt.Println( sum() ) // prints sum of the first element of the three given slices
//  		fmt.Println( sum() ) // prints sum of the second element of the three given slices
// 			...
// Attention!  All slices to be multi-mapped need to be of the same length
func MMapIterInt16(fn func([]int16) int16, seqs ...[]int16) func() int16 {
	// checks
	if len(seqs) < 2 {
		panic(ERR_SHORTER2)
	}
	slen := len(seqs[0])
	for _, seq := range seqs {
		if len(seq) != slen {
			panic(ERR_DIFFLEN)
		}
	}

	iterInt16s := make([]*IterableInt16, len(seqs))
	for i, seq := range seqs {
		iterInt16s[i] = ToIterInt16(seq)
	}
	vals := make([]int16, len(seqs))
	return func() int16 {
		// get seq.Next values
		for i := range vals {
			vals[i] = iterInt16s[i].Next()
		}
		return fn(vals)
	}
}

// ToIterInt16(list interface{}) *IterableInt16
// ToIterInt16 takes a slice and returns an iterator over it
// Attention! If you change the slice while iterableInt16 exists,
// 	  further actions use the changed values
func ToIterInt16(s []int16) *IterableInt16 {

	// Declaration of contextual "global" state variables in scope
	const FIRSTIDX = 0
	var (
		IterLen = len(s)
		LastIdx = IterLen - 1
		ThisIdx = FIRSTIDX - 1
		Exhaust = false
	)

	// ? Assert len(s) ?
	// if IterLen < 2 {
	// 	panic(ERR_SHORTER2)
	// }

	// construct IterableInt16
	var iter = IterableInt16{}

	// checks before returning values
	constraint := func() int16 {
		// sanitize indexing and return elem
		switch {
		case ThisIdx < FIRSTIDX:
			ThisIdx = FIRSTIDX - 1
			Exhaust = true
			return s[FIRSTIDX]
		case ThisIdx > LastIdx:
			ThisIdx = IterLen
			Exhaust = true
			return s[LastIdx]
		}
		return s[ThisIdx]
	}

	// to get the underlying slices length
	iter.Len = IterLen
	// get the iterables internal index
	iter.Index = func() int { return ThisIdx }
	// set the index and return the resulting indexa and state of exhaustion
	iter.SetIndex = func(idx int) (int, bool) {
		ThisIdx = idx
		Exhaust = false
		constraint()
		return ThisIdx, Exhaust
	}

	// Reset these contextual "global" state variables in scope
	iter.Reset = func() {
		ThisIdx = FIRSTIDX - 1
		Exhaust = false
	}

	// Set Idx to the last elem i.e for reverse iteration with iter.Back
	iter.ToEnd = func() { ThisIdx = IterLen; Exhaust = false }

	// Return the first elem in iterable
	// No change in index nor reset (s. above)
	iter.First = func() int16 {
		return s[FIRSTIDX]
	}
	// Return the value at the actual index (again)
	iter.This = func() int16 {
		return constraint()
	}

	// Return the last elem in iterable
	// No change in index nor reset (s. above)
	iter.Last = func() int16 {
		return s[LastIdx]
	}

	// Next iteration and return the content; incr idx
	iter.Next = func() int16 {
		ThisIdx++
		return constraint()
	}

	// Go back one iteration and return value; decr idx
	iter.Back = func() int16 {
		ThisIdx--
		return constraint()
	}

	// cycle endlessly over the iterable (like a ring) returning its elements
	iter.Cycle = func() func() int16 {
		iter.Reset()
		return func() int16 {
			ThisIdx++
			if ThisIdx < IterLen {
				return s[ThisIdx]
			}
			ThisIdx = FIRSTIDX
			return s[ThisIdx]
		}
	}

	// Any(needle int16) reports if needle is an element of IterableFloat's underlying slice
	// The loop implementation allows to get the first index with .Index()
	// 		i.e.
	//      seq := itertools.ToIterInt16([]int16{1, 2, 3})
	//      if seq.Any(2) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.Any = func(needle int16) bool {
		for i, v := range s {
			if v == needle {
				ThisIdx = i
				return true
			}
		}
		return false
	}

	// All(needle int16) reports equality to needle of all elements in IterableFloat's underlying slice
	// The loop implementation allows to get the first index that is != needle with .Index()
	// 		i.e.
	//      seq := itertools.ToIterInt16([]int16{1, 2, 3})
	//      if seq.All(1) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.All = func(needle int16) bool {
		for i, v := range s {
			if v != needle {
				ThisIdx = i
				return false
			}
		}
		return true
	}

	// map

	// MapNext(mapFn) applies the map-function to every element before returning the element Next by Next
	// Does not change the underlying slice
	iter.MapNext = func(fn func(int16) int16) func() (int16, bool) {
		iter.Reset()
		return func() (int16, bool) {
			return fn(iter.Next()), Exhaust
		}
	}

	// MapInto(mapFn) applies the mapFunction to all elements changing the underlying slice to the result and returns itself
	// No additional memory needed
	iter.MapInto = func(fn func(int16) int16) *IterableInt16 {
		iter.Reset()
		for i, v := range s {
			s[i] = fn(v)
		}
		return &iter
	}

	// Map(mapFn) applies the mapFunction to all elements and returns a new iterator with the resulting values
	// Uses memory (new slice with the originals dimensions) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Map = func(fn func(int16) int16) *IterableInt16 {
		iter.Reset()
		newIter := make([]int16, IterLen)
		for i, v := range s {
			newIter[i] = fn(v)
		}
		return ToIterInt16(newIter)
	}

	// filter

	// FilterNext(condition) returns only elements that do meet the filter condition and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.FilterNext = func(cond func(int16) bool) func() (int16, bool) {
		iter.Reset()
		return func() (int16, bool) {
			for {
				cand := iter.Next()
				if Exhaust {
					break
				}
				if !cond(cand) {
					continue
				}
				return cand, Exhaust
			}
			return MININT16, true
		}
	}

	// Filter(filtercondition) returns a new iterable containing only that original elements that do meet the filtercondition
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Filter = func(cond func(int16) bool) *IterableInt16 {
		iter.Reset()
		newIter := make([]int16, 0, IterLen)
		for _, v := range s {
			if cond(v) {
				newIter = append(newIter, v)
			}
		}
		newIter = append(make([]int16, 0, len(newIter)), newIter...)
		return ToIterInt16(newIter)
	}

	// Where(filtercondition) returns an *IterableInt with the indices at which filtercondition is met
	// You might get a slice with the indices by using .List() with the result.
	// Uses memory (new []int slice up to the .Len)
	// Does not change the underlying original slice
	iter.Where = func(cond func(int16) bool) *IterableInt {
		iter.Reset()
		indices := make([]int, 0, IterLen)
		for i, v := range s {
			if cond(v) {
				indices = append(indices, i)
			}
		}
		indices = append(make([]int, 0, len(indices)), indices...)
		return ToIterInt(indices)
	}

	// reduce

	// Reduce(reducerFn) uses the reducerFn function to run over all elements and return one resulting value
	// Does not change the underlying original slice
	iter.Reduce = func(fn func(int16, int16) int16) int16 {
		iter.Reset()
		ThisIdx = FIRSTIDX
		state := constraint()
		for i := 1; i < IterLen; i++ {
			state = fn(state, s[i])
		}
		return state
	}

	// pairwise operation

	// PairOp(fn(x, x+1), [stepwidth=2]) returns a new iterable (len/2) that contains the result of the function applied successivly
	// to a pair of elements then jumping forward to the next pair by stepwidth
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.PairOp = func(fn func(int16, int16) int16, stp ...int) *IterableInt16 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		step, length := 2, IterLen>>1
		if len(stp) == 1 {
			step, length = stp[0], IterLen/stp[0]-1
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		newIter := make([]int16, length)
		ThisIdx = FIRSTIDX
		for i := range newIter {
			newIter[i] = fn(s[ThisIdx], s[ThisIdx+1])
			ThisIdx += step
		}
		return ToIterInt16(newIter)
	}
	// PairOpNext(fn(x, x+1), [stepwidth=2]) returns the result (len/2) of the function applied stepwise
	// to a pair of elements then jumping forward to the next pair by stepwidth and the state of
	// exhaustion of the iterable (index < length)
	// No additional memory use
	// Does not change the underlying original slice
	iter.PairOpNext = func(fn func(int16, int16) int16, stp ...int) func() (int16, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		if IterLen&1 == 1 {
			panic(ERR_ODDLEN)
		}
		step := 2
		if len(stp) == 1 {
			step = stp[0]
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		ThisIdx = FIRSTIDX
		return func() (int16, bool) {
			if Exhaust = ThisIdx > LastIdx; Exhaust {
				return iter.Last(), Exhaust
			}
			a, b := s[ThisIdx], s[ThisIdx+1]
			ThisIdx += step
			return fn(a, b), Exhaust
		}
	}

	// operation on the double previous-actual

	// DoubleOp(fn(prev, actual)) returns a new iterable (len -1) that contains the result of the function applied successivly
	// to the previos and the actual element
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.DoubleOp = func(fn func(int16, int16) int16) *IterableInt16 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		prev := constraint()
		newIter := make([]int16, LastIdx)
		for i := range newIter {
			val := iter.Next()
			newIter[i] = fn(prev, val)
			prev = val
		}
		return ToIterInt16(newIter)
	}
	// DoubleOpNext(fn(prev, actual)) returns the result of the function applied stepwise
	// to the previos and the actual element walking over Next element of iterable and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.DoubleOpNext = func(fn func(int16, int16) int16) func() (int16, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		prev := s[FIRSTIDX]
		val := s[FIRSTIDX]
		ThisIdx = FIRSTIDX
		return func() (int16, bool) {
			prev = val
			val = iter.Next()
			return fn(prev, val), Exhaust
		}
	}

	// comparison of the double previous-actual

	// DoubleComp(Filtercondition(prev, actual)) returns a new iterable that contains only that original elements that do meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.DoubleComp = func(cond func(int16, int16) bool) *IterableInt16 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		newIter := make([]int16, 0, IterLen)
		for ThisIdx < LastIdx {
			prev := val
			val = iter.Next()
			if cond(prev, val) {
				newIter = append(newIter, val)
			}

		}
		newIter = append(make([]int16, 0, len(newIter)), newIter...)
		return ToIterInt16(newIter)
	}
	// DoubleCompNext(Filtercondition(prev, actual)) returns stepwise Next original elements that does meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// No additional memory used as it iterates over the underlying slice
	// Does not change the underlying original slice
	iter.DoubleCompNext = func(cond func(int16, int16) bool) func() (int16, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		return func() (int16, bool) {
			prev := val
			val = iter.Next()
			for !Exhaust && !cond(prev, val) {
				prev = val
				val = iter.Next()
			}
			return val, Exhaust
		}
	}

	// List() returns underlying slice containing the elements of the iterable
	// leaves the iterable unchanged
	// Attention! If you chnge the returned list, that will change the iterable's underlying list too
	iter.List = func() []int16 {
		return s
	}

	// ToList() returns a New slice containing the elements of the iterable
	// allocates new memory
	iter.ToList = func() []int16 {
		list := make([]int16, IterLen)
		copy(list, s)
		return list
	}

	// Tee(number) breaks the iterable into a number of new iterables over the same underlying slice iterable uses
	// make shure not to change the underlying slaice of iter to prevent undesired consequences
	iter.Tee = func(n int) (iters []*IterableInt16) {
		if n < 1 {
			panic("Tee(n) with n smaller 2 ?")
		}
		iter.Reset()
		interval := IterLen / n
		allEqual := IterLen%n == 0
		if allEqual {
			iters = make([]*IterableInt16, 0, interval)
		} else {
			iters = make([]*IterableInt16, 0, interval)
			interval++
		}
		idx := FIRSTIDX
		for ; idx < IterLen-interval; idx += interval {
			iters = append(iters, ToIterInt16(s[idx:idx+interval]))
		}
		iters = append(iters, ToIterInt16(s[idx:IterLen]))
		return iters
	}

	// Destroy unreferences actual context
	iter.Destroy = func() {
		iter = *ToIterInt16(make([]int16, 0))
	}

	return &iter
}

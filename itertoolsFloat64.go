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

// Type IterableFloat64 an iterable over a slice of any type
// Use for any type of iterables that are CAN NOT use []int or []float64, i.e. string.
// IterableIf needs more memory
//      (because it likely will replicate the underlying slice
// 		and Next element in []iterface{} has one memoryword extra used for the elements type)
type IterableFloat64 struct {
	// stepwise returns / does not destroy original underlying slice / no additional memory
	Reset, ToEnd                  func()
	Next, This, Back, First, Last func() float64
	Cycle                         func() func() float64
	MapNext                       func(func(float64) float64) func() (float64, bool)
	FilterNext                    func(func(float64) bool) func() (float64, bool)
	DoubleOpNext                  func(func(float64, float64) float64) func() (float64, bool)
	DoubleCompNext                func(func(float64, float64) bool) func() (float64, bool)
	PairOpNext                    func(func(float64, float64) float64, ...int) func() (float64, bool)

	// info / does not destroy original underlying slice
	Len      int
	Index    func() int
	SetIndex func(int) (int, bool)
	Any, All func(float64) bool
	Where    func(func(float64) bool) *IterableInt

	// conversions & abstractions / do not destroy the iterable nor the original underlying slice
	List   func() []float64
	ToList func() []float64 // needs additional memory
	Reduce func(func(float64, float64) float64) float64
	Tee    func(int) []*IterableFloat64

	// return of new iterable(s) / does not destroy or change original / needs additional memory
	DoubleOp   func(func(float64, float64) float64) *IterableFloat64
	DoubleComp func(func(float64, float64) bool) *IterableFloat64
	PairOp     func(func(float64, float64) float64, ...int) *IterableFloat64
	Filter     func(func(float64) bool) *IterableFloat64
	Map        func(func(float64) float64) *IterableFloat64

	// Return the iterabel BUT changes the underlying slice!
	MapInto func(func(float64) float64) *IterableFloat64

	// Replace the iterable's underlying slice by a slice []float64 with length 0
	Destroy func()
}

// ZipToIterFloat64(l1, l2) *IterableFloat64
// takes to slices and returns a iterator over the zipping result
// zipps two slices and creates a 2*length []interface{} slice from it
//   -- Same size slices
//      or first smaller than the second (second will be cut-off an the first length)
func ZipToIterFloat64(l1, l2 []float64) *IterableFloat64 {
	if len(l1) != len(l2) {
		panic(ERR_SHORTER2)
	}
	l1l2 := make([]float64, 0, len(l1)<<1)
	for i := range l1 {
		l1l2 = append(l1l2, l1[i], l2[i])
	}
	return ToIterFloat64(l1l2)
}

// ChainToIterFloat64(...lists) takes at least 2 slices and returns an newly created iterable over these.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainToIterFloat64(lists ...[]float64) *IterableFloat64 {
	if len(lists) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range lists {
		totalLen += len(lists[i])
	}
	chain := make([]float64, totalLen)
	idx := 0
	for i := range lists {
		copy(chain[idx:], lists[i])
		idx += len(lists[i])
	}
	return ToIterFloat64(chain)
}

// ChainIterFloat64(...iters) takes at least 2 iterables of the same type and returns an newly created iterable over concat.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainIterFloat64(iters ...*IterableFloat64) *IterableFloat64 {
	if len(iters) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range iters {
		totalLen += iters[i].Len
	}
	chain := make([]float64, totalLen)
	idx := 0
	for i := range iters {
		copy(chain[idx:], iters[i].List())
		idx += iters[i].Len
	}
	return ToIterFloat64(chain)
}

// MMapToIterFloat64(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
// maps a function fn to all slices given as comma separated parameter
// and returns an iterator over the slice with the results.
// func fn works on a slice with the i'th element from the iter conversed seq and needs to
//      the slice has a length of the number of slices (0..n) passed to MapIterIf(fn, seqs...)
//      code should adress the source seq by its index 0..n
//      i.e. a function to add the values of three byte slices and return that sum
//			fn := func(v []float64) float64{} { return v[0] + v[1] + v[2] }
// 			sum := MMapIterIf(fn, slice1, slice2, slice3)
// 			fmt.Println( sum.Next() ) // prints sum of the first element of the three given slices
//  		fmt.Println( sum.Next() ) // prints sum of the second element of the three given slices
// 			...
// Attention!  All slices to be multi-mapped need to be of the same length
// Needs additional memory (size of one of the given slices)
func MMapToIterFloat64(fn func([]float64) float64, seqs ...[]float64) *IterableFloat64 {
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

	newIter := make([]float64, slen)
	vals := make([]float64, len(seqs))
	for i := range newIter {
		for ii := range vals {
			vals[ii] = seqs[ii][i]
		}
		newIter[i] = fn(vals)
	}
	return ToIterFloat64(newIter)
}

// MMapIterFloat64(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
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
func MMapIterFloat64(fn func([]float64) float64, seqs ...[]float64) func() float64 {
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

	iterFloat64s := make([]*IterableFloat64, len(seqs))
	for i, seq := range seqs {
		iterFloat64s[i] = ToIterFloat64(seq)
	}
	vals := make([]float64, len(seqs))
	return func() float64 {
		// get seq.Next values
		for i := range vals {
			vals[i] = iterFloat64s[i].Next()
		}
		return fn(vals)
	}
}

// ToIterFloat64(list interface{}) *IterableFloat64
// ToIterFloat64 takes a slice and returns an iterator over it
// Attention! If you change the slice while iterableFloat64 exists,
// 	  further actions use the changed values
func ToIterFloat64(s []float64) *IterableFloat64 {

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

	// construct IterableFloat64
	var iter = IterableFloat64{}

	// checks before returning values
	constraint := func() float64 {
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
	iter.First = func() float64 {
		return s[FIRSTIDX]
	}
	// Return the value at the actual index (again)
	iter.This = func() float64 {
		return constraint()
	}

	// Return the last elem in iterable
	// No change in index nor reset (s. above)
	iter.Last = func() float64 {
		return s[LastIdx]
	}

	// Next iteration and return the content; incr idx
	iter.Next = func() float64 {
		ThisIdx++
		return constraint()
	}

	// Go back one iteration and return value; decr idx
	iter.Back = func() float64 {
		ThisIdx--
		return constraint()
	}

	// cycle endlessly over the iterable (like a ring) returning its elements
	iter.Cycle = func() func() float64 {
		iter.Reset()
		return func() float64 {
			ThisIdx++
			if ThisIdx < IterLen {
				return s[ThisIdx]
			}
			ThisIdx = FIRSTIDX
			return s[ThisIdx]
		}
	}

	// Any(needle float64) reports if needle is an element of IterableFloat's underlying slice
	// The loop implementation allows to get the first index with .Index()
	// 		i.e.
	//      seq := itertools.ToIterFloat64([]float64{1, 2, 3})
	//      if seq.Any(2) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.Any = func(needle float64) bool {
		for i, v := range s {
			if v == needle {
				ThisIdx = i
				return true
			}
		}
		return false
	}

	// All(needle float64) reports equality to needle of all elements in IterableFloat's underlying slice
	// The loop implementation allows to get the first index that is != needle with .Index()
	// 		i.e.
	//      seq := itertools.ToIterFloat64([]float64{1, 2, 3})
	//      if seq.All(1) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.All = func(needle float64) bool {
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
	iter.MapNext = func(fn func(float64) float64) func() (float64, bool) {
		iter.Reset()
		return func() (float64, bool) {
			return fn(iter.Next()), Exhaust
		}
	}

	// MapInto(mapFn) applies the mapFunction to all elements changing the underlying slice to the result and returns itself
	// No additional memory needed
	iter.MapInto = func(fn func(float64) float64) *IterableFloat64 {
		iter.Reset()
		for i, v := range s {
			s[i] = fn(v)
		}
		return &iter
	}

	// Map(mapFn) applies the mapFunction to all elements and returns a new iterator with the resulting values
	// Uses memory (new slice with the originals dimensions) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Map = func(fn func(float64) float64) *IterableFloat64 {
		iter.Reset()
		newIter := make([]float64, IterLen)
		for i, v := range s {
			newIter[i] = fn(v)
		}
		return ToIterFloat64(newIter)
	}

	// filter

	// FilterNext(condition) returns only elements that do meet the filter condition and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.FilterNext = func(cond func(float64) bool) func() (float64, bool) {
		iter.Reset()
		return func() (float64, bool) {
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
			return MINFLOAT64, true
		}
	}

	// Filter(filtercondition) returns a new iterable containing only that original elements that do meet the filtercondition
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Filter = func(cond func(float64) bool) *IterableFloat64 {
		iter.Reset()
		newIter := make([]float64, 0, IterLen)
		for _, v := range s {
			if cond(v) {
				newIter = append(newIter, v)
			}
		}
		newIter = append(make([]float64, 0, len(newIter)), newIter...)
		return ToIterFloat64(newIter)
	}

	// Where(filtercondition) returns an *IterableInt with the indices at which filtercondition is met
	// You might get a slice with the indices by using .List() with the result.
	// Uses memory (new []int slice up to the .Len)
	// Does not change the underlying original slice
	iter.Where = func(cond func(float64) bool) *IterableInt {
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
	iter.Reduce = func(fn func(float64, float64) float64) float64 {
		iter.Reset()
		ThisIdx = FIRSTIDX
		state := constraint()
		for i := 1; i < IterLen; i++ {
			state = fn(state, s[i])
		}
		return state
	}

	// pairwise operation

	// PairOp(fn(prev, actual), [stepwidth=2]) returns a new iterable (len/2) that contains the result of the function applied successivly
	// to a pair of elements then jumping forward to the next pair by stepwidth
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.PairOp = func(fn func(float64, float64) float64, stp ...int) *IterableFloat64 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		step, length := 2, IterLen>>1
		if len(stp) == 1 && stp[0] != 2 {
			step, length = stp[0], IterLen/stp[0]-1
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		newIter := make([]float64, 0, length)
		for i := step - 1; i < len(s); i += step {
			newIter = append(newIter, fn(s[i-1], s[i]))
		}
		return ToIterFloat64(newIter)
	}
	// PairOpNext(fn(prev, actual), [stepwidth=2]) returns the result (len/2) of the function applied stepwise
	// to a pair of elements then jumping forward to the next pair by stepwidth and the state of
	// exhaustion of the iterable (index < length)
	// No additional memory use
	// Does not change the underlying original slice
	iter.PairOpNext = func(fn func(float64, float64) float64, stp ...int) func() (float64, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		if IterLen&1 == 1 {
			panic(ERR_ODDLEN)
		}
		step := 2
		if len(stp) == 1 && stp[0] != 2 {
			step = stp[0]
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		ThisIdx = step - 1
		return func() (float64, bool) {
			if Exhaust = ThisIdx > LastIdx; Exhaust {
				return iter.Last(), Exhaust
			}
			a, b := s[ThisIdx-1], s[ThisIdx]
			ThisIdx += step
			return fn(a, b), Exhaust
		}
	}

	// operation on the double previous-actual

	// DoubleOp(fn(prev, actual)) returns a new iterable (len -1) that contains the result of the function applied successivly
	// to the previos and the actual element
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.DoubleOp = func(fn func(float64, float64) float64) *IterableFloat64 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		prev := constraint()
		newIter := make([]float64, LastIdx)
		for i := range newIter {
			val := iter.Next()
			newIter[i] = fn(prev, val)
			prev = val
		}
		return ToIterFloat64(newIter)
	}
	// DoubleOpNext(fn(prev, actual)) returns the result of the function applied stepwise
	// to the previos and the actual element walking over Next element of iterable and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.DoubleOpNext = func(fn func(float64, float64) float64) func() (float64, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		prev := s[FIRSTIDX]
		val := s[FIRSTIDX]
		ThisIdx = FIRSTIDX
		return func() (float64, bool) {
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
	iter.DoubleComp = func(cond func(float64, float64) bool) *IterableFloat64 {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		newIter := make([]float64, 0, IterLen)
		for ThisIdx < LastIdx {
			prev := val
			val = iter.Next()
			if cond(prev, val) {
				newIter = append(newIter, val)
			}

		}
		newIter = append(make([]float64, 0, len(newIter)), newIter...)
		return ToIterFloat64(newIter)
	}
	// DoubleCompNext(Filtercondition(prev, actual)) returns stepwise Next original elements that does meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// No additional memory used as it iterates over the underlying slice
	// Does not change the underlying original slice
	iter.DoubleCompNext = func(cond func(float64, float64) bool) func() (float64, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		return func() (float64, bool) {
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
	iter.List = func() []float64 {
		return s
	}

	// ToList() returns a New slice containing the elements of the iterable
	// allocates new memory
	iter.ToList = func() []float64 {
		list := make([]float64, IterLen)
		copy(list, s)
		return list
	}

	// Tee(number) breaks the iterable into a number of new iterables over the same underlying slice iterable uses
	// make shure not to change the underlying slaice of iter to prevent undesired consequences
	iter.Tee = func(n int) (iters []*IterableFloat64) {
		if n < 1 {
			panic("Tee(n) with n smaller 2 ?")
		}
		iter.Reset()
		interval := IterLen / n
		allEqual := IterLen%n == 0
		if allEqual {
			iters = make([]*IterableFloat64, 0, interval)
		} else {
			iters = make([]*IterableFloat64, 0, interval)
			interval++
		}
		idx := FIRSTIDX
		for ; idx < IterLen-interval; idx += interval {
			iters = append(iters, ToIterFloat64(s[idx:idx+interval]))
		}
		iters = append(iters, ToIterFloat64(s[idx:IterLen]))
		return iters
	}

	// Destroy unreferences actual context
	iter.Destroy = func() {
		iter = *ToIterFloat64(make([]float64, 0))
	}

	return &iter
}

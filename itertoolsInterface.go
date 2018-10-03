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

// Type IterableIf an iterable over a slice of any type
// Use for any type of iterables that are CAN NOT use []int or []float64, i.e. string.
// IterableIf needs more memory
//      (because it likely will replicate the underlying slice
// 		and Next element in []iterface{} has one memoryword extra used for the elements type)
type IterableIf struct {
	// info
	Len      int
	Index    func() int
	SetIndex func(int) (int, bool)
	Any, All func(interface{}) bool
	Where    func(func(interface{}) bool) *IterableInt
	// stepwise returns / does not destroy original underlying slice / no additional memory
	Reset, ToEnd                  func()
	Next, This, Back, First, Last func() interface{}
	Cycle                         func() func() interface{}
	MapNext                       func(func(interface{}) interface{}) func() (interface{}, bool)
	FilterNext                    func(func(interface{}) bool) func() (interface{}, bool)
	DoubleOpNext                  func(func(interface{}, interface{}) interface{}) func() (interface{}, bool)
	DoubleCompNext                func(func(interface{}, interface{}) bool) func() (interface{}, bool)
	PairOpNext                    func(func(interface{}, interface{}) interface{}, ...int) func() (interface{}, bool)
	//
	Reduce func(func(interface{}, interface{}) interface{}) interface{}
	// return of a new iterable(s) / does not destroy or change original underlying slice / needs additional memory
	DoubleOp   func(func(interface{}, interface{}) interface{}) *IterableIf
	DoubleComp func(func(interface{}, interface{}) bool) *IterableIf
	PairOp     func(func(interface{}, interface{}) interface{}, ...int) *IterableIf
	Filter     func(func(interface{}) bool) *IterableIf
	Map        func(func(interface{}) interface{}) *IterableIf
	Tee        func(int) []*IterableIf
	// return of the very same iterable by changing original underlying slice
	MapInto func(func(interface{}) interface{}) *IterableIf
	// returns the underlying slice of the iterables / does not destroy original underlying slice nor iterable
	List func() []interface{}
	// returns a new slice with the iterables elements / does not destroy original underlying slice nor iterable
	ToList func() []interface{}
	// returns nothing / does not destroy original underlying slice / destroys the iterable by nil
	Destroy func()
}

// ZipToIterIf(s1, s2 interface{}) *IterableIf
// takes to slices and returns a iterator over the zipping result
// zipps two slices and creates a 2*length []interface{} slice from it
//   -- memory use: 4*InputSlice (because of 2 memwords per elem)
// Attention! No length check!
//   -- Same size slices
//      or first smaller than the second (second will be cut-off an the first length)
func ZipToIterIf(s1, s2 interface{}) *IterableIf {
	s1If, _ := convertToInterfaceSlice(s1)
	s2If, _ := convertToInterfaceSlice(s2)
	s1s2 := make([]interface{}, 0, len(s1If)<<1)
	for i := range s1If {
		s1s2 = append(s1s2, s1If[i], s2If[i])
	}
	return ToIterIf(s1s2)
}

// // ChainToIterIf(lists) takes at least 2 slices and returns an newly created iterable over these.
// // Uses additional memory to build an underlying slice for the new iterable
// // Does not change the original slices in lists
// func ChainToIterIf(lists ...interface{}) *IterableIf {
// 	if len(lists) < 2 {
// 		panic(ERR_SHORTER2)
// 	}
// 	totalLen := 0
// 	for i := range lists {
// 		totalLen += len(lists[i])
// 	}
// 	chain := make([]interface{}, totalLen)
// 	idx := 0
// 	for i := range lists {
// 		copy(chain[idx:], lists[i])
// 		idx += len(lists[i])
// 	}
// 	return ToIterIf(chain)
// }

// MMapIterIf(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
// maps a function fn to all slices given as comma separated parameter
// and yields the results.
// func fn works on a slice with the i'th element from the iter conversed seq and needs to
//      the slice has a length of the number of slices (0..n) passed to MapIterIf(fn, seqs...)
//      code should adress the source seq by its index 0..n
//      i.e. a function to add the values of three byte slices and return that sum
//			fn := func(v []interface{}) interface{} { return v[0].(uint8) + v[1].(uint8) + v[2].(uint8) }
// 			sum := MMapIterIf(fn, slice1, slice2, slice3)
// 			fmt.Println( sum() ) // prints sum of the first element of the three given slices
//  		fmt.Println( sum() ) // prints sum of the second element of the three given slices
// 			...
// Attention! There is no checking!
//   - all slices to be mapped need to be of the same type
//		   or mix but at least must all allow for the casting and the methods in the mapping
//   - best: all slices have same length - at least first slice must be the smallest
//   - mapfunction: make sure the internal type of the interfaced data allows the methods used (i.a '+')
//                  make shure the result fits into format
func MMapIterIf(fn func([]interface{}) interface{}, seqs ...interface{}) func() interface{} {
	iterIfs := make([]*IterableIf, len(seqs))
	for i, seq := range seqs {
		iterIfs[i] = ToIterIf(seq)
	}
	vals := make([]interface{}, len(seqs))
	return func() interface{} {
		// get seq.Next values
		for i := range vals {
			vals[i] = iterIfs[i].Next()
		}
		return fn(vals)
	}
}

// convertToInterfaceSlice(list interface{})
// creates a new []interface{} from a list given as a parameter
// returns []interface{} and errorVal with the minimum value distict to the type
// -- memory use 2*InputSlice Mem
func convertToInterfaceSlice(list interface{}) ([]interface{}, interface{}) {
	var (
		s      []interface{}
		errval interface{}
	)

	switch list.(type) {
	case []uint8:
		s = make([]interface{}, len(list.([]uint8)))
		for i, v := range list.([]uint8) {
			s[i] = v
		}
		errval = 0 // MinUint8
	case []int32:
		s = make([]interface{}, len(list.([]int32)))
		for i, v := range list.([]int32) {
			s[i] = v
		}
		errval = MININT32
	case []int:
		s = make([]interface{}, len(list.([]int)))
		for i, v := range list.([]int) {
			s[i] = v
		}
		errval = MININT
	case []float32:
		s = make([]interface{}, len(list.([]float32)))
		for i, v := range list.([]float32) {
			s[i] = v
		}
		errval = MINFLOAT32 // MinFloat32
	case []float64:
		s = make([]interface{}, len(list.([]float64)))
		for i, v := range list.([]float64) {
			s[i] = v
		}
		errval = MINFLOAT64
	case []string:
		s = make([]interface{}, len(list.([]string)))
		for i, v := range list.([]string) {
			s[i] = v
		}
		errval = "" // empty string
	case []interface{}:
		s = list.([]interface{})
		errval = interface{}(nil) // empty string
	}
	return s, errval
}

// ToIterIf(list interface{}) *IterableIf
// ToIterIf takes any kind of list and returns an iterator over it
// 		with replicating (Plus 2*memory of tthe slice) or changing the underlying slice
//		exception: slice is []interface{} which is taken by reference
// Attention! If you change the slice
func ToIterIf(list interface{}) *IterableIf {
	// declarations
	const FIRSTIDX = 0
	var (
		s, ErrorVal = convertToInterfaceSlice(list)
		IterLen     = len(s)
		LastIdx     = IterLen - 1
		ThisIdx     = FIRSTIDX - 1
		Exhaust     = false
	)

	// ? Assert len(s) ?
	// if IterLen < 2 {
	// 	panic(ERR_SHORTER2)
	// }

	// construct iterable
	var iter = IterableIf{}

	// checks before returning values
	constraint := func() interface{} {
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

	// Reset all contextual "global" state variables in scope
	iter.Reset = func() {
		ThisIdx = FIRSTIDX - 1
		Exhaust = false
	}

	// return the value at the actual index (again)
	iter.This = func() interface{} {
		return constraint()
	}

	// Next iteration and return the content
	iter.Next = func() interface{} {
		ThisIdx++
		return constraint()
	}

	// Set Idx to the Last Item i.e for reverse iteration with iter.Back
	iter.ToEnd = func() { ThisIdx = IterLen }

	// Go back one iteration and return value
	iter.Back = func() interface{} {
		ThisIdx--
		return constraint()
	}

	// return the first value in iterable
	// no change in index nor reset (s. above)
	iter.First = func() interface{} {
		return s[FIRSTIDX]
	}
	// return the last value in iterable
	// no change in index nor reset (s. above)
	iter.Last = func() interface{} {
		return s[LastIdx]
	}

	// cycle endlessly over the iterable (like a ring) returning its elements
	iter.Cycle = func() func() interface{} {
		iter.Reset()
		return func() interface{} {
			ThisIdx++
			if ThisIdx < IterLen {
				return s[ThisIdx]
			}
			ThisIdx = FIRSTIDX
			return s[ThisIdx]
		}
	}

	// Any(needle interface{}) reports if needle is an element of IterableIf's underlying slice
	// The loop implementation allows to get the first index with .Index()
	// 		i.e.
	//      seq := itertools.ToIterIf([]interface{}{1, 2, 3})
	//      if seq.Any(2) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.Any = func(needle interface{}) bool {
		for i, v := range s {
			if v == needle {
				ThisIdx = i
				return true
			}
		}
		return false
	}

	// All(needle interface{}) reports equality to needle of all elements in IterableIf's underlying slice
	// The loop implementation allows to get the first index that is != needle with .Index()
	// 		i.e.
	//      seq := itertools.ToIterIf([]interface{}{1, 2, 3})
	//      if seq.All(1) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.All = func(needle interface{}) bool {
		for i, v := range s {
			if v != needle {
				ThisIdx = i
				return false
			}
		}
		return true
	}

	// MapNext(mapFn) applies the map-function to every element before returning the element Next by Next
	// Does not change the underlying slice
	iter.MapNext = func(mapFn func(interface{}) interface{}) func() (interface{}, bool) {
		iter.Reset()
		return func() (interface{}, bool) {
			return mapFn(iter.Next()), Exhaust
		}
	}

	// MapInto(mapFn) applies the mapFunction to all elements changing the underlying slice to the result and returns itself
	// No additional memory needed
	iter.MapInto = func(fn func(interface{}) interface{}) *IterableIf {
		iter.Reset()
		for i, v := range s {
			s[i] = fn(v)
		}
		return &iter
	}

	// Map(mapFn) applies the mapFunction to all elements and returns a new iterator with the resulting values
	// Uses memory (new slice with the originals dimensions) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Map = func(fn func(interface{}) interface{}) *IterableIf {
		iter.Reset()
		newIter := make([]interface{}, IterLen)
		for i, v := range s {
			newIter[i] = fn(v)
		}
		return ToIterIf(newIter)
	}

	// FilterNext(condition) returns only elements that do meet the filter condition and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.FilterNext = func(cond func(interface{}) bool) func() (interface{}, bool) {
		iter.Reset()
		return func() (interface{}, bool) {
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
			return ErrorVal, true
		}
	}

	// Filter(condition) returns a new iterable that contains only that original elements that do meet the filter condition
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Filter = func(cond func(interface{}) bool) *IterableIf {
		iter.Reset()
		newIter := make([]interface{}, 0, IterLen)
		for _, v := range s {
			if cond(v) {
				newIter = append(newIter, v)
			}
		}
		newIter = append(make([]interface{}, 0, len(newIter)), newIter...)
		return ToIterIf(newIter)
	}

	// Where(filtercondition) returns an *IterableInt with the indices at which filtercondition is met
	// You might get a slice with the indices by using .List() with the result.
	// Uses memory (new []int slice up to the .Len)
	// Does not change the underlying original slice
	iter.Where = func(cond func(interface{}) bool) *IterableInt {
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

	// pairwise operation

	// PairOp(fn(x, x+1), [stepwidth=2]) returns a new iterable (len/2) that contains the result of the function applied successivly
	// to a pair of elements then jumping forward to the next pair by stepwidth
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.PairOp = func(fn func(interface{}, interface{}) interface{}, stp ...int) *IterableIf {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		step := 2
		if len(stp) == 1 {
			step = stp[0]
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		ThisIdx = FIRSTIDX
		newIter := make([]interface{}, IterLen/step)
		for i := range newIter {
			newIter[i] = fn(s[ThisIdx], s[ThisIdx+1])
			ThisIdx += step
		}
		return ToIterIf(newIter)
	}
	// PairOpNext(fn(x, x+1), [stepwidth=2]) returns the result (len/2) of the function applied stepwise
	// to a pair of elements then jumping forward to the next pair by stepwidth and the state of
	// exhaustion of the iterable (index < length)
	// No additional memory use
	// Does not change the underlying original slice
	iter.PairOpNext = func(fn func(interface{}, interface{}) interface{}, stp ...int) func() (interface{}, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		step := 2
		if len(stp) == 1 {
			step = stp[0]
		}
		if IterLen%step != 0 {
			panic(ERR_WRONGLEN)
		}
		ThisIdx = FIRSTIDX
		return func() (interface{}, bool) {
			a, b := s[ThisIdx], s[ThisIdx+1]
			Exhaust = ThisIdx >= LastIdx
			ThisIdx += step
			return fn(a, b), Exhaust
		}
	}

	// DoubleOp(fn(prev, actual)) returns a new iterable (len -1) that contains the result of the function applied successivly
	// to the previos and the actual element
	// Uses memory (new slice with the originals dimensions minus one) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.DoubleOp = func(fn func(interface{}, interface{}) interface{}) *IterableIf {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		prev := constraint()
		newIter := make([]interface{}, LastIdx)
		for i := range newIter {
			ThisIdx++
			val := s[ThisIdx]
			newIter[i] = fn(prev, val)
			prev = val
		}
		return ToIterIf(newIter)
	}

	// DoubleOpNext(fn(prev, actual)) returns the result of the function applied stepwise
	// to the previos and the actual element walking over Next element of iterable and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.DoubleOpNext = func(fn func(interface{}, interface{}) interface{}) func() (interface{}, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		prev := s[FIRSTIDX]
		val := s[FIRSTIDX]
		ThisIdx = FIRSTIDX
		return func() (interface{}, bool) {
			prev = val
			ThisIdx++
			val = constraint()
			return fn(prev, val), Exhaust
		}
	}

	// DoubleCom(Filtercondition(prev, actual)) returns a new iterable that contains only that original elements that do meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.DoubleComp = func(cond func(interface{}, interface{}) bool) *IterableIf {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		newIter := make([]interface{}, 0, IterLen)
		for ThisIdx < LastIdx {
			prev := val
			val = iter.Next()
			if cond(prev, val) {
				newIter = append(newIter, val)
			}
		}
		newIter = append(make([]interface{}, 0, len(newIter)), newIter...)
		return ToIterIf(newIter)
	}

	// DoubleComNext(Filtercondition(prev, actual)) returns stepwise Next original elements that does meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// No additional memory used as it iterates over the underlying slice
	// Does not change the underlying original slice
	iter.DoubleCompNext = func(cond func(interface{}, interface{}) bool) func() (interface{}, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		return func() (interface{}, bool) {
			prev := val
			val = iter.Next()
			for !Exhaust && !cond(prev, val) {
				prev = val
				val = iter.Next()
			}
			return val, Exhaust
		}
	}

	// Reduce(reducerFn) uses the reducerFn function to run over all elements and return one resulting value
	// Does not change the underlying original slice
	iter.Reduce = func(fn func(interface{}, interface{}) interface{}) interface{} {
		iter.Reset()
		ThisIdx = FIRSTIDX
		state := constraint()
		for i := 1; i < IterLen; i++ {
			state = fn(state, s[i])
		}
		return state
	}

	// List() returns underlying slice containing the elements of the iterable
	// leaves the iterable unchanged
	// Attention! If you chnge the returned list, that will change the iterable's underlying list too
	iter.List = func() []interface{} {
		return s
	}

	// ToList() returns a New slice containing the elements of the iterable
	// allocates new memory
	iter.ToList = func() []interface{} {
		list := make([]interface{}, IterLen)
		copy(list, s)
		return list
	}

	// Tee(number) breaks the iterable into a number of new iterables over the same underlying slice iterable uses
	// make shure not to change the underlying slaice of iter to prevent undesired consequences
	iter.Tee = func(n int) (iters []*IterableIf) {
		if n < 1 {
			panic("Tee(n) with n smaller 2 ?")
		}
		iter.Reset()
		interval := IterLen / n
		allEqual := IterLen%n == 0
		if allEqual {
			iters = make([]*IterableIf, 0, interval)
		} else {
			iters = make([]*IterableIf, 0, interval)
			interval++
		}
		idx := FIRSTIDX
		for ; idx < IterLen-interval; idx += interval {
			iters = append(iters, ToIterIf(s[idx:idx+interval]))
		}
		iters = append(iters, ToIterIf(s[idx:IterLen]))
		return iters
	}

	iter.Destroy = func() {
		iter = *ToIterIf(make([]interface{}, 0))
	}

	return &iter
}

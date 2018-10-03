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
// Timestamp: 2018-09-21 17:31:33.636594889 +0200 CEST m=+0.003431049
// 

package itertools

// Type IterableByte an iterable over a slice of any type
// Use for any type of iterables that are CAN NOT use []int or []byte, i.e. string.
// IterableIf needs more memory
//      (because it likely will replicate the underlying slice
// 		and Next element in []iterface{} has one memoryword extra used for the elements type)
type IterableByte struct {
	// stepwise returns / does not destroy original underlying slice / no additional memory
	Reset, ToEnd                  func()
	Next, This, Back, First, Last func() byte
	Cycle                         func() func() byte
	MapNext                       func(func(byte) byte) func() (byte, bool)
	FilterNext                    func(func(byte) bool) func() (byte, bool)
	DoubleOpNext                  func(func(byte, byte) byte) func() (byte, bool)
	DoubleCompNext                func(func(byte, byte) bool) func() (byte, bool)
	PairOpNext                    func(func(byte, byte) byte, ...int) func() (byte, bool)

	// info / does not destroy original underlying slice
	Len      int
	Index    func() int
	SetIndex func(int) (int, bool)
	Any, All func(byte) bool
	Where    func(func(byte) bool) *IterableInt

	// conversions & abstractions / do not destroy the iterable nor the original underlying slice
	List   func() []byte
	ToList func() []byte // needs additional memory
	Reduce func(func(byte, byte) byte) byte
	Tee    func(int) []*IterableByte

	// return of new iterable(s) / does not destroy or change original / needs additional memory
	DoubleOp   func(func(byte, byte) byte) *IterableByte
	DoubleComp func(func(byte, byte) bool) *IterableByte
	PairOp     func(func(byte, byte) byte, ...int) *IterableByte
	Filter     func(func(byte) bool) *IterableByte
	Map        func(func(byte) byte) *IterableByte

	// Return the iterabel BUT changes the underlying slice!
	MapInto func(func(byte) byte) *IterableByte

	// Replace the iterable's underlying slice by a slice []byte with length 0
	Destroy func()
}

// ZipToIterByte(l1, l2) *IterableByte
// takes to slices and returns a iterator over the zipping result
// zipps two slices and creates a 2*length []interface{} slice from it
//   -- Same size slices
//      or first smaller than the second (second will be cut-off an the first length)
func ZipToIterByte(l1, l2 []byte) *IterableByte {
	if len(l1) != len(l2) {
		panic(ERR_SHORTER2)
	}
	l1l2 := make([]byte, 0, len(l1)<<1)
	for i := range l1 {
		l1l2 = append(l1l2, l1[i], l2[i])
	}
	return ToIterByte(l1l2)
}

// ChainToIterByte(...lists) takes at least 2 slices and returns an newly created iterable over these.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainToIterByte(lists ...[]byte) *IterableByte {
	if len(lists) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range lists {
		totalLen += len(lists[i])
	}
	chain := make([]byte, totalLen)
	idx := 0
	for i := range lists {
		copy(chain[idx:], lists[i])
		idx += len(lists[i])
	}
	return ToIterByte(chain)
}

// ChainIterByte(...iters) takes at least 2 iterables of the same type and returns an newly created iterable over concat.
// Uses additional memory to build an underlying slice for the new iterable
// Does not change the original slices in lists
func ChainIterByte(iters ...*IterableByte) *IterableByte {
	if len(iters) < 2 {
		panic(ERR_SHORTER2)
	}
	totalLen := 0
	for i := range iters {
		totalLen += iters[i].Len
	}
	chain := make([]byte, totalLen)
	idx := 0
	for i := range iters {
		copy(chain[idx:], iters[i].List())
		idx += iters[i].Len
	}
	return ToIterByte(chain)
}

// MMapToIterByte(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
// maps a function fn to all slices given as comma separated parameter
// and returns an iterator over the slice with the results.
// func fn works on a slice with the i'th element from the iter conversed seq and needs to
//      the slice has a length of the number of slices (0..n) passed to MapIterIf(fn, seqs...)
//      code should adress the source seq by its index 0..n
//      i.e. a function to add the values of three byte slices and return that sum
//			fn := func(v []byte) byte{} { return v[0] + v[1] + v[2] }
// 			sum := MMapIterIf(fn, slice1, slice2, slice3)
// 			fmt.Println( sum.Next() ) // prints sum of the first element of the three given slices
//  		fmt.Println( sum.Next() ) // prints sum of the second element of the three given slices
// 			...
// Attention!  All slices to be multi-mapped need to be of the same length
// Needs additional memory (size of one of the given slices)
func MMapToIterByte(fn func([]byte) byte, seqs ...[]byte) *IterableByte {
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

	newIter := make([]byte, slen)
	vals := make([]byte, len(seqs))
	for i := range newIter {
		for ii := range vals {
			vals[ii] = seqs[ii][i]
		}
		newIter[i] = fn(vals)
	}
	return ToIterByte(newIter)
}

// MMapIterByte(fn func([]interface{}) interface{}, seqs ...[]interface{}) func() interface{}
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
func MMapIterByte(fn func([]byte) byte, seqs ...[]byte) func() byte {
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

	iterBytes := make([]*IterableByte, len(seqs))
	for i, seq := range seqs {
		iterBytes[i] = ToIterByte(seq)
	}
	vals := make([]byte, len(seqs))
	return func() byte {
		// get seq.Next values
		for i := range vals {
			vals[i] = iterBytes[i].Next()
		}
		return fn(vals)
	}
}

// ToIterByte(list interface{}) *IterableByte
// ToIterByte takes a slice and returns an iterator over it
// Attention! If you change the slice while iterableByte exists,
// 	  further actions use the changed values
func ToIterByte(s []byte) *IterableByte {

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

	// construct IterableByte
	var iter = IterableByte{}

	// checks before returning values
	constraint := func() byte {
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
	iter.First = func() byte {
		return s[FIRSTIDX]
	}
	// Return the value at the actual index (again)
	iter.This = func() byte {
		return constraint()
	}

	// Return the last elem in iterable
	// No change in index nor reset (s. above)
	iter.Last = func() byte {
		return s[LastIdx]
	}

	// Next iteration and return the content; incr idx
	iter.Next = func() byte {
		ThisIdx++
		return constraint()
	}

	// Go back one iteration and return value; decr idx
	iter.Back = func() byte {
		ThisIdx--
		return constraint()
	}

	// cycle endlessly over the iterable (like a ring) returning its elements
	iter.Cycle = func() func() byte {
		iter.Reset()
		return func() byte {
			ThisIdx++
			if ThisIdx < IterLen {
				return s[ThisIdx]
			}
			ThisIdx = FIRSTIDX
			return s[ThisIdx]
		}
	}

	// Any(needle byte) reports if needle is an element of IterableFloat's underlying slice
	// The loop implementation allows to get the first index with .Index()
	// 		i.e.
	//      seq := itertools.ToIterByte([]byte{1, 2, 3})
	//      if seq.Any(2) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.Any = func(needle byte) bool {
		for i, v := range s {
			if v == needle {
				ThisIdx = i
				return true
			}
		}
		return false
	}

	// All(needle byte) reports equality to needle of all elements in IterableFloat's underlying slice
	// The loop implementation allows to get the first index that is != needle with .Index()
	// 		i.e.
	//      seq := itertools.ToIterByte([]byte{1, 2, 3})
	//      if seq.All(1) {
	//			fmt.Printf("index: %v\n", seq.Index())
	//      }
	iter.All = func(needle byte) bool {
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
	iter.MapNext = func(fn func(byte) byte) func() (byte, bool) {
		iter.Reset()
		return func() (byte, bool) {
			return fn(iter.Next()), Exhaust
		}
	}

	// MapInto(mapFn) applies the mapFunction to all elements changing the underlying slice to the result and returns itself
	// No additional memory needed
	iter.MapInto = func(fn func(byte) byte) *IterableByte {
		iter.Reset()
		for i, v := range s {
			s[i] = fn(v)
		}
		return &iter
	}

	// Map(mapFn) applies the mapFunction to all elements and returns a new iterator with the resulting values
	// Uses memory (new slice with the originals dimensions) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Map = func(fn func(byte) byte) *IterableByte {
		iter.Reset()
		newIter := make([]byte, IterLen)
		for i, v := range s {
			newIter[i] = fn(v)
		}
		return ToIterByte(newIter)
	}

	// filter

	// FilterNext(condition) returns only elements that do meet the filter condition and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.FilterNext = func(cond func(byte) bool) func() (byte, bool) {
		iter.Reset()
		return func() (byte, bool) {
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
			return MINBYTE, true
		}
	}

	// Filter(filtercondition) returns a new iterable containing only that original elements that do meet the filtercondition
	// Uses memory (new slice with the originals dimensions minus the skipped elements) and the new iterable refers to this new slice
	// Does not change the underlying original slice
	iter.Filter = func(cond func(byte) bool) *IterableByte {
		iter.Reset()
		newIter := make([]byte, 0, IterLen)
		for _, v := range s {
			if cond(v) {
				newIter = append(newIter, v)
			}
		}
		newIter = append(make([]byte, 0, len(newIter)), newIter...)
		return ToIterByte(newIter)
	}

	// Where(filtercondition) returns an *IterableInt with the indices at which filtercondition is met
	// You might get a slice with the indices by using .List() with the result.
	// Uses memory (new []int slice up to the .Len)
	// Does not change the underlying original slice
	iter.Where = func(cond func(byte) bool) *IterableInt {
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
	iter.Reduce = func(fn func(byte, byte) byte) byte {
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
	iter.PairOp = func(fn func(byte, byte) byte, stp ...int) *IterableByte {
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
		newIter := make([]byte, length)
		ThisIdx = FIRSTIDX
		for i := range newIter {
			newIter[i] = fn(s[ThisIdx], s[ThisIdx+1])
			ThisIdx += step
		}
		return ToIterByte(newIter)
	}
	// PairOpNext(fn(x, x+1), [stepwidth=2]) returns the result (len/2) of the function applied stepwise
	// to a pair of elements then jumping forward to the next pair by stepwidth and the state of
	// exhaustion of the iterable (index < length)
	// No additional memory use
	// Does not change the underlying original slice
	iter.PairOpNext = func(fn func(byte, byte) byte, stp ...int) func() (byte, bool) {
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
		return func() (byte, bool) {
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
	iter.DoubleOp = func(fn func(byte, byte) byte) *IterableByte {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		prev := constraint()
		newIter := make([]byte, LastIdx)
		for i := range newIter {
			val := iter.Next()
			newIter[i] = fn(prev, val)
			prev = val
		}
		return ToIterByte(newIter)
	}
	// DoubleOpNext(fn(prev, actual)) returns the result of the function applied stepwise
	// to the previos and the actual element walking over Next element of iterable and a bool indicator
	// for the exhaustion of the iterable (index < length)
	// Does not change the underlying original slice
	iter.DoubleOpNext = func(fn func(byte, byte) byte) func() (byte, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		prev := s[FIRSTIDX]
		val := s[FIRSTIDX]
		ThisIdx = FIRSTIDX
		return func() (byte, bool) {
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
	iter.DoubleComp = func(cond func(byte, byte) bool) *IterableByte {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		newIter := make([]byte, 0, IterLen)
		for ThisIdx < LastIdx {
			prev := val
			val = iter.Next()
			if cond(prev, val) {
				newIter = append(newIter, val)
			}

		}
		newIter = append(make([]byte, 0, len(newIter)), newIter...)
		return ToIterByte(newIter)
	}
	// DoubleCompNext(Filtercondition(prev, actual)) returns stepwise Next original elements that does meet
	// the filtercondition between the previous and the actual element (i.e p == a)
	// Note! Allways starts with atleast the second element - first element cannot be compared as "actual" and thus is never included
	// No additional memory used as it iterates over the underlying slice
	// Does not change the underlying original slice
	iter.DoubleCompNext = func(cond func(byte, byte) bool) func() (byte, bool) {
		iter.Reset()
		if IterLen < 2 {
			panic(ERR_SHORTER2)
		}
		ThisIdx = FIRSTIDX
		val := constraint()
		return func() (byte, bool) {
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
	iter.List = func() []byte {
		return s
	}

	// ToList() returns a New slice containing the elements of the iterable
	// allocates new memory
	iter.ToList = func() []byte {
		list := make([]byte, IterLen)
		copy(list, s)
		return list
	}

	// Tee(number) breaks the iterable into a number of new iterables over the same underlying slice iterable uses
	// make shure not to change the underlying slaice of iter to prevent undesired consequences
	iter.Tee = func(n int) (iters []*IterableByte) {
		if n < 1 {
			panic("Tee(n) with n smaller 2 ?")
		}
		iter.Reset()
		interval := IterLen / n
		allEqual := IterLen%n == 0
		if allEqual {
			iters = make([]*IterableByte, 0, interval)
		} else {
			iters = make([]*IterableByte, 0, interval)
			interval++
		}
		idx := FIRSTIDX
		for ; idx < IterLen-interval; idx += interval {
			iters = append(iters, ToIterByte(s[idx:idx+interval]))
		}
		iters = append(iters, ToIterByte(s[idx:IterLen]))
		return iters
	}

	// Destroy unreferences actual context
	iter.Destroy = func() {
		iter = *ToIterByte(make([]byte, 0))
	}

	return &iter
}
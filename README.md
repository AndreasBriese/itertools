
__itertools__
=============

[![Build Status](https://travis-ci.org/AndreasBriese/itertools.png?branch=master)](http://travis-ci.org/AndreasBriese/itertools)

A fast iterator package for Go. 

Inspired by the python iterator functions this provides analogous abstraction and (kind'like) semantics.

For the sake of speed, the yield-like behaviour is implemented without channels and methods. 
Instead functional scope (context) provides state for the itertools iterators.

The package includes dedicated iterators for  "int", "int64", "int32", "int16", "int8", "float32", "float64" (beeing the template for the other distict type iterators using __go:generate__), "string", "byte" and an implementation using the empty interface{} (type iterableIf) for slices of any other type and for map operations on iterators of mixed type.

Whilst the first two are straightforward and stable to be used with functions to manipulate the emitted iterator elements, caution is needed to meet Go's type checking requirements with type iterableIf.   

NOTE: Any the abstraction costs speed in comparison to a generic slice operations. Unless you need a yield-like behaviour for your programming logic applying a map function to slice elements in a for loop is usually (about 10times) faster than using the iterable MapNext function(ality). Anyway have fun!

----
__Installation__

    go get -u github.com/AndreasBriese/itertools

Import to your code as usual or with a shorter alias.

    import iter "github.com/AndreasBriese/itertools"
    
----
__Functionality__

itertools implements iterators for slices of the following types &lt;t&gt;:

    []float64        IterableFloat64        (Template for go:generate in makeMoreItertools.go)
    
    []int            IterableInt
    []int64          IterableInt64
    []int32          IterableInt32
    []int16          IterableInt16
    []int8           IterableInt8
    []float32        IterableFloat32
    []string         IterableString
    []byte           IterableByte
    
    []interface{}    IterableIf
    
typed **Iterable&lt;T&gt;** (with front capital letter in &lt;T&gt; like IterableInt64) which are constructed by ToIter&lt;T&gt; (with front capital letter <T> i.e. ToIterInt64([]int{}).

Zipping two slices []&lt;t&gt; AAAA and BBBB to a new iterable of type Iterable&lt;T&gt; ABABABAB can be done with ZipToIter&lt;T&gt;(A, B []&lt;t&gt;). 

ChainToIterIf(...[]&lt;t&gt;) takes at least 2 slices and returns an newly created iterable over their concat.

MMapToIter&lt;T&gt; maps multiple slices []&lt;t&gt; with a mapping function to an Iterable&lt;T&gt; over a new slice with the results (Note! This does not work on IterableIf).
MMapIter&lt;T&gt; maps multiple slices []&lt;t&gt; with a mapping function and yields the result stepwise.

__Iterable&lt;T&gt;__ _Functions_

Functions that return stepwise elements without changing the underlying original slice (don't need additional memory):
    
    Next, This, Back, First, Last func() <T>
    Cycle                         func() func() <T>
    MapNext                       func(func(<T>) <T>) func() (<T>, bool)
    FilterNext                    func(func(<T>) bool) func() (<T>, bool)
    DoubleOpNext                  func(func(<T>, <T>) <T>) func() (<T>, bool)
    DoubleCompNext                func(func(<T>, <T>) bool) func() (<T>, bool)
    PairOpNext                    func(func(<T>, <T>) <T>, ...int) func() (<T>, bool)


Functions that return a new iterable without changing the underlying original slice - these need additional memory for the underlying slices:

    PairOp           func(func(<T>, <T>) <T>, ...int) *Iterable<T>
    DoubleOp         func(func(<T>, <T>) <T>) *Iterable<T>
    DoubleComp       func(func(<T>, <T>) bool) *Iterable<T>
    Filter           func(func(<T>) bool) *Iterable<T>
    Map              func(func(<T>) <T>) *Iterable<T>

Functions that return the very same iterable with changes to the underlying original slice - no additional memory needed:

    MapInto          func(func(<T>) <T>) *Iterable<T>

Function that returns iterables over the underlying origanal slice - ~~ little additional memory needed:

    Tee              func(int) []*Iterable<T>
    
Infos & setters about state, that do not change the underlying original slice:

    Len           int
    Index         func() int
    SetIndex      func(int) (int, bool)
    Reset, ToEnd  func()

converters & destroyer of the iterator:

    List   func() []<T> // returns the underlying slice but leaves the iterator intact
    ToList   func() []<T> // returns a copy of underlying slice but leaves the iterator intact, uses additional memory
    Reduce   func(func(<T>, <T>) <T>) <T> // leaves the iterator intact
    Destroy  func() // replaces the underlying slice by an empty slice
    
Godoc provides basic documentation and you might lookup the samples for usage and comparison.

__Some short notes about "specialties"__

1._Why having three variants of Map Iterable&lt;T&gt;.Map() | .MapInto() | .MapNext() ?_

Both, Iterable&lt;T&gt;.Map(fn) and .MapInto(fn), calculate the fn() over the iterable at once and return an Iterable&lt;T&gt; with the resulting elements.

Usage will depend on the use case:

If the original data in the underlying slice can be disposed, then use the faster .MapInto(). The bigger the data the larger the speed difference will be.
Iterable&lt;T&gt;.Map() needs additional memory but does not change the underlying slice.

    l := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    seq := iter.ToIterInt(l)
    fmt.Printf("Result: %v slice l: %v\n", seq.Map(func(x int) int { return 2 * x }).ToList(), l)
    fmt.Printf("Result: %v slice l: %v\n", seq.MapInto(func(x int) int { return 2 * x }).ToList(), l)
    
    // prints
    // Result: [0 2 4 6 8 10 12 14 16 18] slice l: [0 1 2 3 4 5 6 7 8 9]
    // Result: [0 2 4 6 8 10 12 14 16 18] slice l: [0 2 4 6 8 10 12 14 16 18]

Iterable&lt;T&gt;.MapNext() is slower than the others because the results can not produced by a loop over the iterable's slice __s__. Anyway this can be handy if you actually want to change the elements of __s__ dynamically in between the calculation steps or like apply fn repeatedly conditionally by storing its results within __s__ and using .Back(). __s__ is defined for the functions scope with ToIter&lt;T&gt;([]&lt;t&gt;) and might also be changed in the callers context. Anyway: __s__'s data can but it's length and capacity can not be changed (or more acurately: these can be changed i the callers context but the changes do not apply to the Iterables scope).   

2._Iterable&lt;T&gt;.PairOp(fn, step) & iter.PairOpNext(fn, step)_

Applies the function fn in the call to a pair of sequential elements (previous, actual) starting with elem[1] and jumps by _step_ to the next pair. Length of the returned iterable depends on step. That means:

1. iter.PairOp(fn, step=1) runs over the iterable with fn(elem[0], elem[1]), fn(elem[1], elem[2]) ... len(result) = len(iterable)-1 

2. iter.PairOp(fn, step=2) runs over the iterable with fn(elem[0], elem[1]), fn(elem[2], elem[3]) ... len(result) = len(iterable)/2

With step = 2 it can be used together with the ZipToIter&lt;T&gt;() function. 
With large data (and multiprocessor use), operating over the zip iterable has the benefit of looping over one slice instead 
of calling elements from two different slices at different memory locations. 
There is an example of PairOp in the sample directory.

3._Iterable&lt;T&gt;.DoubleComp(fn) & .DoubleCompNext(fn)_

    condFn := func(prev, actual interface{}) bool {
            return (prev.(string) != actual.(string))
    }
    s1 := []string{"A", "B", "T", "T", "A", "Su", "Su", "E", "Roy", "A", "A"}
    fmt.Println("no doubles:", iter.ToIterIf(s1).DoubleComp(condFn).ToList())
    // prints 'no doubles: [B T A Su E Roy A]'
    
    s1 = []string{"A", "A", "B", "T", "T", "A", "Su", "Su", "E", "Roy", "A"}
    fmt.Println("no doubles:", iter.ToIterIf(s1).DoubleComp(condFn).ToList())
    // prints 'no doubles: [B T A Su E Roy A]'
    

Against intuition but accurately this prints  *[B T A Su E Roy A]* because the DoubeComp* functions compares the actual and previous element and returns in case the actual element. Consequently the first element in the iterable can not be compared to previous nor returned. If you need it, call Iterable&lt;T&gt;.First().  

__Benchmarks__ (`go test -bench=. cpu=1`)

Some Benchmarks _IterableInt_ []int (slice length 1.000.000) 

    Benchmark_IterInt_ToIterInt                      1000000          1956 ns/op
    Benchmark_IterInt_ToList                             100       5871707 ns/op
    Benchmark_IterInt_List                         500000000             2.07 ns/op
    Benchmark_IterInt_MapInto                            500       2816308 ns/op
    Benchmark_IterInt_Map                                200       6779983 ns/op
    Benchmark_IterInt_MapNext                            200      10156138 ns/op
    Benchmark_IterInt_Filter                             100      14037458 ns/op
    Benchmark_IterInt_FilterNext                         100      15868025 ns/op
    Benchmark_IterInt_Reduce                             300       4018775 ns/op
    Benchmark_IterInt_Map_Filter_Reduce                  100      19697933 ns/op
    Benchmark_IterInt_MapInto_Filter_Reduce              100      16282103 ns/op
    Benchmark_IterInt_PairOp                             200       8981810 ns/op
    Benchmark_IterInt_DoubleOp                           100      11818775 ns/op
    Benchmark_IterInt_DoubleComp                         100      19677501 ns/op
    Benchmark_IterInt_Tee5                              2000        741668 ns/op

    Benchmark_IterIf_ToIterIf                             50      33606582 ns/op

Run `go test -bench=. cpu=1` to get further benchmarks. Other types will behave more-or-less similar, 
but IterableIf is slower to build the Iterable with its underlying []interface{} (i.e. ToIterIf([]uint8)) 
and in (other) cases where typecasting is involved (i.e. in mapping and filter functions).

Comparison slice []int vs. diverse iterables (slice length 1.000.000) 

_LoopOverIntSlice_ is gold standard (using native go slice for ops)

    Benchmark_Int_Sum_LoopOverIntSlice                  2000        702820 ns/op
    Benchmark_IterInt_Sum_NextOverIterable               200       8659626 ns/op
    Benchmark_Yanatan_Sum_NextOverIterable                 5     219389515 ns/op
    
    mapFunction f(x) = 2*x
    Benchmark_Int_Map_LoopOverIntSlice                  2000        662458 ns/op
    Benchmark_IterInt_Map_OverIterable                   500       3615911 ns/op
    Benchmark_IterInt_Map_IntoOverIterable               500       3073908 ns/op
    Benchmark_IterInt_Map_NextOverIterable               100      11960226 ns/op
    Benchmark_Yanatan_Map_NextOverIterable                 3     441119057 ns/op

(Yanatan refers to the iterables using channels: github.com/yanatan16/itertools)
See https://github.com/ewencp/golang-iterators-benchmark for more comparisons of different implementations.

Benchmarks on Apple MBP 2013 with i7 and 8GB Ram 
 
__License__   

(C)opyright 2018 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt, with MIT license - see the headers

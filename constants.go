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

//go:generate go run makeMoreItertools.go

const (
	MININT     = int(MININT64)
	MINFLOAT32 = float32(-3.4028235e+38)
	MINFLOAT64 = float64(-1.7976931348623157e+308)
	MININT8    = int8(-1 << 7)
	MININT16   = int16(-1 << 15)
	MININT32   = int32(-1 << 31)
	MININT64   = int64(-1 << 63)
	MINBYTE    = byte(0)
	MINSTRING  = ""

	ERR_SHORTER1 = "Iterable with lenght smaller 1"
	ERR_SHORTER2 = "Iterable with lenght smaller 2 - need at least 2 elements"
	ERR_DIFFTYPE = "Can not use different types in an iterable - need purity"
	ERR_DIFFLEN  = "Parameter error: underlying slices differ in length"
	ERR_ODDLEN   = "Pairwise operation: Underlying slices has odd length"
	ERR_WRONGLEN = "Length does not match op function steps"
)

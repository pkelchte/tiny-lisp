// H25.1/17 - H25.1/21 (鈴)

package arith

import (
	. "fmt"
	"math"
	. "math/big"
)

func ExampleIsNumber() {
	Println(IsNumber(int32(7)))
	Println(IsNumber(7.0))
	Println(IsNumber(NewRat(7, 2)))
	Println(IsNumber("7"))
	Println(IsNumber(int64(7)))
	// Output:
	// true
	// true
	// true
	// false
	// false
}

func ExampleAdd() {
	r := Add(7, 2)
	Printf("%T %v\n", r, r)
	r = Add(7, 2.0)
	Printf("%T %v\n", r, r)
	r = Add(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Add(7.0, 2)
	Printf("%T %v\n", r, r)
	r = Add(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = Add(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Add(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = Add(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = Add(NewRat(7, 1), NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Add(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	// Output:
	// int32 9
	// float64 9
	// int32 9
	// float64 9
	// float64 9
	// float64 9
	// int32 9
	// float64 9
	// int32 9
	// *big.Rat 25/6
}

func ExampleSubtract() {
	r := Subtract(7, 2)
	Printf("%T %v\n", r, r)
	r = Subtract(7, 2.0)
	Printf("%T %v\n", r, r)
	r = Subtract(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Subtract(7.0, 2)
	Printf("%T %v\n", r, r)
	r = Subtract(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = Subtract(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Subtract(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = Subtract(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = Subtract(NewRat(7, 1), NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Subtract(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	// Output:
	// int32 5
	// float64 5
	// int32 5
	// float64 5
	// float64 5
	// float64 5
	// int32 5
	// float64 5
	// int32 5
	// *big.Rat 17/6
}

func ExampleMultiply() {
	r := Multiply(7, 2)
	Printf("%T %v\n", r, r)
	r = Multiply(7, 2.0)
	Printf("%T %v\n", r, r)
	r = Multiply(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Multiply(7.0, 2)
	Printf("%T %v\n", r, r)
	r = Multiply(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = Multiply(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Multiply(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = Multiply(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = Multiply(NewRat(7, 1), NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Multiply(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	// Output:
	// int32 14
	// float64 14
	// int32 14
	// float64 14
	// float64 14
	// float64 14
	// int32 14
	// float64 14
	// int32 14
	// *big.Rat 7/3
}

func ExampleDivideReal() {
	r := DivideReal(7, 2)
	Printf("%T %v\n", r, r)
	r = DivideReal(7, 2.0)
	Printf("%T %v\n", r, r)
	r = DivideReal(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideReal(7.0, 2)
	Printf("%T %v\n", r, r)
	r = DivideReal(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = DivideReal(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideReal(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = DivideReal(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = DivideReal(NewRat(7, 1), NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideReal(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	// Output:
	// *big.Rat 7/2
	// float64 3.5
	// *big.Rat 7/2
	// float64 3.5
	// float64 3.5
	// float64 3.5
	// *big.Rat 7/2
	// float64 3.5
	// *big.Rat 7/2
	// *big.Rat 21/4
}

func ExampleDivideInt() {
	r := DivideInt(7, 2)
	Printf("%T %v\n", r, r)
	r = DivideInt(7, 2.0)
	Printf("%T %v\n", r, r)
	r = DivideInt(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideInt(7.0, 2)
	Printf("%T %v\n", r, r)
	r = DivideInt(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = DivideInt(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideInt(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = DivideInt(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = DivideInt(NewRat(7, 1), NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = DivideInt(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	// Output:
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 3
	// int32 5
}

func ExampleCompare() {
	r := Compare(7, 7)
	Printf("%T %v\n", r, r)
	r = Compare(7, 2)
	Printf("%T %v\n", r, r)
	r = Compare(2, 7)
	Printf("%T %v\n", r, r)
	r = Compare(7, 2.0)
	Printf("%T %v\n", r, r)
	r = Compare(7, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Compare(7.0, 2)
	Printf("%T %v\n", r, r)
	r = Compare(7.0, 2.0)
	Printf("%T %v\n", r, r)
	r = Compare(7.0, NewRat(2, 1))
	Printf("%T %v\n", r, r)
	r = Compare(NewRat(7, 1), 2)
	Printf("%T %v\n", r, r)
	r = Compare(NewRat(7, 1), 2.0)
	Printf("%T %v\n", r, r)
	r = Compare(NewRat(7, 1), 7)
	Printf("%T %v\n", r, r)
	r = Compare(NewRat(7, 2), NewRat(2, 3))
	Printf("%T %v\n", r, r)
	r = Compare(NewRat(2, 3), NewRat(7, 2))
	Printf("%T %v\n", r, r)
	// Output:
	// int 0
	// int 1
	// int -1
	// int 1
	// int 1
	// int 1
	// int 1
	// int 1
	// int 1
	// int 1
	// int 0
	// int 1
	// int -1
}

func ExampleFloat64() {
	r := Float64(7)
	Printf("%T %v\n", r, r)
	r = Float64(7.0)
	Printf("%T %v\n", r, r)
	r = Float64(NewRat(7, 1))
	Printf("%T %v\n", r, r)
	r = Float64(NewRat(-7, 2))
	Printf("%T %v\n", r, r)
	// Output:
	// float64 7
	// float64 7
	// float64 7
	// float64 -3.5
}

func ExampleTruncate() {
	r := Truncate(7)
	Printf("%T %v\n", r, r)
	r = Truncate(int32(7))
	Printf("%T %v\n", r, r)
	r = Truncate(int64(7))
	Printf("%T %v\n", r, r)
	r = Truncate(7.0)
	Printf("%T %v\n", r, r)
	r = Truncate(NewRat(7, 1))
	Printf("%T %v\n", r, r)
	r = Truncate(NewRat(7, 2))
	Printf("%T %v\n", r, r)
	r = Truncate(NewRat(7000000000000, 2))
	Printf("%T %v\n", r, r)
	r = Truncate(-3.5)
	Printf("%T %v\n", r, r)
	r = Truncate(NewRat(-7, 2))
	Printf("%T %v\n", r, r)
	r = Truncate(2147483647.999)
	Printf("%T %v\n", r, r)
	r = Truncate(2147483648.000)
	Printf("%T %v\n", r, r)
	r = Truncate(-2147483648.999)
	Printf("%T %v\n", r, r)
	r = Truncate(-2147483649.000)
	Printf("%T %v\n", r, r)
	r = Truncate(math.Inf(1))
	Printf("%T %v\n", r, r)
	// Output:
	// int32 7
	// int32 7
	// int32 7
	// int32 7
	// int32 7
	// int32 3
	// *big.Rat 3500000000000/1
	// int32 -3
	// int32 -3
	// int32 2147483647
	// *big.Rat 2147483648/1
	// int32 -2147483648
	// *big.Rat -2147483649/1
	// float64 +Inf
}

func ExampleString() {
	r := String(7)
	Printf("%T %v\n", r, r)
	r = String(int32(7))
	Printf("%T %v\n", r, r)
	r = String(int64(7))
	Printf("%T %v\n", r, r)
	r = String(7.0)
	Printf("%T %v\n", r, r)
	r = String(NewRat(7, 1))
	Printf("%T %v\n", r, r)
	r = String(NewRat(7, 2))
	Printf("%T %v\n", r, r)
	r = String(7.2)
	Printf("%T %v\n", r, r)
	r = String(7e20)
	Printf("%T %v\n", r, r)
	r = String(-7.0)
	Printf("%T %v\n", r, r)
	r = String(-7.2)
	Printf("%T %v\n", r, r)
	z, _ := new(Rat).SetString("123456789012345678901234567890")
	r = String(z)
	Printf("%T %v\n", r, r)
	// Output:
	// string 7
	// string 7
	// string 7
	// string 7.0
	// string 7
	// string 7/2
	// string 7.2
	// string 7e+20
	// string -7.0
	// string -7.2
	// string 123456789012345678901234567890
}

func ExampleNumber_factorial() {
	var n Number = 1
	for i := 1; i < 30; i++ {
		n = Multiply(n, i)
		Println(String(n))
	}
	// Output:
	// 1
	// 2
	// 6
	// 24
	// 120
	// 720
	// 5040
	// 40320
	// 362880
	// 3628800
	// 39916800
	// 479001600
	// 6227020800
	// 87178291200
	// 1307674368000
	// 20922789888000
	// 355687428096000
	// 6402373705728000
	// 121645100408832000
	// 2432902008176640000
	// 51090942171709440000
	// 1124000727777607680000
	// 25852016738884976640000
	// 620448401733239439360000
	// 15511210043330985984000000
	// 403291461126605635584000000
	// 10888869450418352160768000000
	// 304888344611713860501504000000
	// 8841761993739701954543616000000
}

/*
  Copyright (c) 2013 OKI Software Co., Ltd.

  Permission is hereby granted, free of charge, to any person obtaining a
  copy of this software and associated documentation files (the "Software"),
  to deal in the Software without restriction, including without limitation
  the rights to use, copy, modify, merge, publish, distribute, sublicense,
  and/or sell copies of the Software, and to permit persons to whom the
  Software is furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
  FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
  DEALINGS IN THE SOFTWARE.
*/

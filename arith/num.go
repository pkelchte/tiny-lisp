// H25.1/16 - H25.1/21 (鈴)

/*
 パッケージ arith は int32, float64, Rat ("math/big" の無限多倍長有理数)
 の自動的な相互変換を伴う算術演算を実装する.
 int32 の演算を内部的に int64 で行うことで桁あふれ時の Rat への自動的な
 変換を実現する.
*/
package arith

import (
	. "fmt"
	"math"
	. "math/big"
	"strconv"
	"strings"
)

// Number は int32, float64 または *Rat を表す型であるとする
// (実際の型検査は実行時に行う)。
// ただし，関数の引数としては便宜のため int, int64 も許す
// (それらは実行時の値によって int32 または *Rat へと変換される)。
type Number interface{}

// Number 引数を int32, float64, *Rat のどれかに当てはめる。
func regulateArg(a Number) Number {
	switch x := a.(type) {
	case int:
		var i int32 = int32(x)
		if int(i) == x {
			return i
		}
		return NewRat(int64(x), 1)
	case int64:
		return regulateInt64(x)
	case int32:
		return x
	case float64:
		return x
	case *Rat:
		return x
	}
	panic(Errorf("unsupported type: %T", a))
}

// int64 の値を可能ならば int32 の値に，できなければ有理数にする。
func regulateInt64(a int64) Number {
	var i int32 = int32(a)
	if int64(i) == a {
		return i
	}
	return NewRat(a, 1)
}

// 有理数を float64 の値にする。
// もしも範囲を越えるときは *strconv.NumError で panic する。
func ratToFloat64(a *Rat) float64 {
	s := a.FloatString(17)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// float64 の値をゼロの方向へ丸めた整数を返す。
// 無限大 (+Inf) などで丸めることができないときはそのままの値を返す。
func truncateFloat64(a float64) Number {
	if math.MinInt32-1 < a && a < math.MaxInt32+1 {
		return int32(a)
	}
	s := strconv.FormatFloat(a, 'g', -1, 64)
	r, ok := new(Rat).SetString(s)
	if !ok {
		return a
	}
	return truncateRat(r)
}

// float64 に対する符号関数
func signum(a float64) int {
	if a < 0 {
		return -1
	} else if a == 0 {
		return 0
	}
	return 1
}

// 有理数を可能ならば int32 の値にする。
func regulateRat(a *Rat) Number {
	if a.IsInt() {
		n := a.Num()
		if n.BitLen() < 32 { // 符号ビットを除き 31 ビット長以内か？
			return int32(n.Int64())
		}
	}
	return a
}

// 有理数をゼロの方向へ丸めた整数を返す。
func truncateRat(a *Rat) Number {
	n := a.Num()
	d := a.Denom()
	q := new(Int).Quo(n, d)
	if q.BitLen() < 32 { // 符号ビットを除き 31 ビット長以内か？
		return int32(q.Int64())
	}
	return new(Rat).SetInt(q)
}

// 以下，公開関数

// int32, float64 または *Rat ならば true を返す。
func IsNumber(a interface{}) bool {
	switch a.(type) {
	case int32:
		return true
	case float64:
		return true
	case *Rat:
		return true
	}
	return false
}

// 加算: Add(7, 2) => 9
func Add(a, b Number) Number {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			return regulateInt64(int64(x) + int64(y))
		case float64:
			return float64(x) + y
		case *Rat:
			z := NewRat(int64(x), 1)
			return regulateRat(z.Add(z, y))
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return x + float64(y)
		case float64:
			return x + y
		case *Rat:
			return x + ratToFloat64(y)
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return regulateRat(z.Add(x, z))
		case float64:
			return ratToFloat64(x) + y
		case *Rat:
			return regulateRat(new(Rat).Add(x, y))
		}
	}
	return Add(regulateArg(a), regulateArg(b))
}

// 減算: Subtract(7, 2) => 5
func Subtract(a, b Number) Number {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			return regulateInt64(int64(x) - int64(y))
		case float64:
			return float64(x) - y
		case *Rat:
			z := NewRat(int64(x), 1)
			return regulateRat(z.Sub(z, y))
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return x - float64(y)
		case float64:
			return x - y
		case *Rat:
			return x - ratToFloat64(y)
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return regulateRat(z.Sub(x, z))
		case float64:
			return ratToFloat64(x) - y
		case *Rat:
			return regulateRat(new(Rat).Sub(x, y))
		}
	}
	return Subtract(regulateArg(a), regulateArg(b))
}

// 乗算: Multiply(7, 2) => 14
func Multiply(a, b Number) Number {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			return regulateInt64(int64(x) * int64(y))
		case float64:
			return float64(x) * y
		case *Rat:
			z := NewRat(int64(x), 1)
			return regulateRat(z.Mul(z, y))
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return x * float64(y)
		case float64:
			return x * y
		case *Rat:
			return x * ratToFloat64(y)
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return regulateRat(z.Mul(x, z))
		case float64:
			return ratToFloat64(x) * y
		case *Rat:
			return regulateRat(new(Rat).Mul(x, y))
		}
	}
	return Multiply(regulateArg(a), regulateArg(b))
}

// 実数除算: DivideReal(7, 2) => 7/2
func DivideReal(a, b Number) Number {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			return regulateRat(NewRat(int64(x), int64(y)))
		case float64:
			return float64(x) / y
		case *Rat:
			z := NewRat(int64(x), 1)
			return regulateRat(z.Quo(z, y))
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return x / float64(y)
		case float64:
			return x / y
		case *Rat:
			return x / ratToFloat64(y)
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return regulateRat(z.Quo(x, z))
		case float64:
			return ratToFloat64(x) / y
		case *Rat:
			return regulateRat(new(Rat).Quo(x, y))
		}
	}
	return DivideReal(regulateArg(a), regulateArg(b))
}

// 整数除算: DivideInt(7, 2) => 3
func DivideInt(a, b Number) Number {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			return regulateInt64(int64(x) / int64(y))
		case float64:
			return truncateFloat64(float64(x) / y)
		case *Rat:
			z := NewRat(int64(x), 1)
			return truncateRat(z.Quo(z, y))
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return truncateFloat64(x / float64(y))
		case float64:
			return truncateFloat64(x / y)
		case *Rat:
			return truncateFloat64(x / ratToFloat64(y))
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return truncateRat(z.Quo(x, z))
		case float64:
			return truncateFloat64(ratToFloat64(x) / y)
		case *Rat:
			return truncateRat(new(Rat).Quo(x, y))
		}
	}
	return DivideInt(regulateArg(a), regulateArg(b))
}

// 比較: a, b について a が b に対して小さい/等しい/大きいとき，
// それぞれ -1/0/1 を返す。Compare(7, 2) => 1
func Compare(a, b Number) int {
	switch x := a.(type) {
	case int32:
		switch y := b.(type) {
		case int32:
			if x < y {
				return -1
			} else if x == y {
				return 0
			} else {
				return 1
			}
		case float64:
			return signum(float64(x) - y)
		case *Rat:
			z := NewRat(int64(x), 1)
			return z.Cmp(y)
		}
	case float64:
		switch y := b.(type) {
		case int32:
			return signum(x - float64(y))
		case float64:
			return signum(x - y)
		case *Rat:
			return signum(x - ratToFloat64(y))
		}
	case *Rat:
		switch y := b.(type) {
		case int32:
			z := NewRat(int64(y), 1)
			return x.Cmp(z)
		case float64:
			return signum(ratToFloat64(x) - y)
		case *Rat:
			return x.Cmp(y)
		}
	}
	return Compare(regulateArg(a), regulateArg(b))
}

// 数を浮動小数点数にした値を返す。Float64(7) => 7.0
func Float64(a Number) float64 {
	switch x := a.(type) {
	case int32:
		return float64(x)
	case float64:
		return x
	case *Rat:
		return ratToFloat64(x)
	}
	return Float64(regulateArg(a))
}

// 数をゼロの方向へ丸めた整数を返す。Truncate(7.2) => 7
func Truncate(a Number) Number {
	switch x := a.(type) {
	case int32:
		return x
	case float64:
		return truncateFloat64(x)
	case *Rat:
		return truncateRat(x)
	}
	return Truncate(regulateArg(a))
}

// 数の文字列表現を得る。String(7.0) => "7.0"
func String(a Number) string {
	switch x := a.(type) {
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case float64:
		s := strconv.FormatFloat(x, 'g', -1, 64)
		if !strings.ContainsAny(s, ".e") {
			s = s + ".0"
		}
		return s
	case *Rat:
		if x.IsInt() {
			return x.Num().String()
		}
		return x.String()
	}
	return String(regulateArg(a))
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

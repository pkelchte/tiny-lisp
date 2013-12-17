// H25.3/18 - 4/15 (鈴)

// このファイルは cons セルとシンボルその他データを実装する。

package lisp

import (
	"github.com/pkelchte/tiny-lisp/arith"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

// 任意の型.
// 注意: 空リストは *Cell 型の nil (つまり Any 型の非 nil) で表現される。
type Any interface{}

// Cons セル型
type Cell struct {
	Car Any
	Cdr *Cell // improper list にはしない
}

// シンボル型.  同じ文字列に対してアドレスは常に一意。
type Symbol struct {
	string
}

// 新しい cons セルを作る。
func Cons(car Any, cdr *Cell) *Cell {
	return &Cell{car, cdr}
}

var symbols = make(map[string]*Symbol)
var lock sync.Mutex

// 文字列からシンボルを作る。
func NewSymbol(name string) *Symbol {
	lock.Lock()
	sym, ok := symbols[name]
	if !ok {
		sym = &Symbol{name}
		symbols[name] = sym
	}
	lock.Unlock()
	return sym
}

// 定義済みのシンボルを用意する。

var QuoteSymbol = NewSymbol("quote")
var TSymbol = NewSymbol("t")
var NilSymbol = NewSymbol("nil")
var AmpRestSymbol = NewSymbol("&rest")

// デバッグの便宜のための cons セルの文字列表現
func (cell *Cell) String() string {
	return fmt.Sprintf("(%v . %v)", cell.Car, cell.Cdr)
}

// StringFor 関数で同じリストを再帰的に表示する深さ
var MaxPrintRecur = 4

// Lisp 式としての引数の文字列表現を返す。
func StringFor(a Any) string {
	return stringFor(a, MaxPrintRecur, make(map[*Cell]bool))
}

func stringFor(a Any, recurLevel int, printed map[*Cell]bool) string {
	switch x := a.(type) {
	case *Cell:
		if x != nil && x.Car == QuoteSymbol {
			if x.Cdr != nil && x.Cdr.Cdr == nil {
				return "'" + stringFor(x.Cdr.Car, recurLevel, printed)
			}
		}
		return "(" + stringForList(x, recurLevel, printed) + ")"
	case *Symbol:
		return x.string
	case string:
		return strconv.Quote(x)
	case *big.Rat:
		s1 := arith.String(x)
		s2 := fmt.Sprintf("%v", arith.Float64(x))
		if s1 == s2 {
			return s1
		} else if r, ok := NumberFor(s2); ok && arith.Compare(x, r) == 0 {
			return fmt.Sprintf("%s /*=%s*/", s1, s2)
		} else {
			return fmt.Sprintf("%s /*~%s*/", s1, s2)
		}
	}
	return fmt.Sprintf("%v", a)
}

func stringForList(x *Cell, recurLevel int, printed map[*Cell]bool) string {
	if x == nil {
		return ""
	}
	s := make([]string, 0, 10)
	var y *Cell
	for y = x; y != nil; y = y.Cdr {
		if _, ok := printed[y]; ok {
			recurLevel--
			if recurLevel < 0 {
				s = append(s, "...") // 循環リストを ... で表す
				break
			}
		} else {
			printed[y] = true
			recurLevel = MaxPrintRecur
		}
		e := stringFor(y.Car, recurLevel, printed)
		s = append(s, e)
	}
	if y == nil { // 最後まで到達できたならば非循環リストである
		for y := x; y != nil; y = y.Cdr {
			delete(printed, y)
		}
	}
	return strings.Join(s, " ")
}

// 文字列に対する整数または有理数を得る。変換できないとき論理値に偽を返す。
func NumberFor(text string) (arith.Number, bool) {
	i, err := strconv.ParseInt(text, 0, 32)
	if err == nil {
		return i, true
	}
	r := new(big.Rat)
	n := new(big.Int)
	_, ok := n.SetString(text, 0) // これは 0x… 0b… 0… を解釈する
	if ok {
		r.SetInt(n)
		return r, true
	}
	return r.SetString(text) //  これは 0x… 0b… 0… を解釈しない
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

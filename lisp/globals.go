// H25.3/18 - 4/1 (鈴)

// このファイルはトップレベルの環境を実装する。

package lisp

import (
	"github.com/pkelchte/tiny-lisp/arith"
	"fmt"
)

// トップレベルの環境
var Globals = &Env{map[*Symbol]Any{
	TSymbol:          TSymbol,
	NewSymbol("car"): carFunc, NewSymbol("cdr"): cdrFunc,
	NewSymbol("cons"):  consFunc,
	NewSymbol("listp"): listpFunc, NewSymbol("eq"): eqFunc,
	NewSymbol("rplaca"): rplacaFunc, NewSymbol("rplacd"): rplacdFunc,
	NewSymbol("list"): listFunc,
	NewSymbol("="):    eqOp, NewSymbol("/="): neOp,
	NewSymbol("<"): ltOp, NewSymbol("<="): leOp,
	NewSymbol(">"): gtOp, NewSymbol(">="): geOp,
	NewSymbol("+"): addOp, NewSymbol("-"): subtractOp,
	NewSymbol("*"): multiplyOp, NewSymbol("/"): divideOp,
	NewSymbol("gensym"): gensymFunc,
	NewSymbol("print"):  printFunc,
	QuoteSymbol:         quoteForm,
	NewSymbol("setq"):   setqForm,
	NewSymbol("progn"):  prognForm, NewSymbol("if"): ifForm,
	NewSymbol("lambda"): lambdaForm, NewSymbol("let"): letForm,
	NewSymbol("defun"): defunForm,
	NewSymbol("apply"): applyForm, NewSymbol("and"): andForm,
}, nil}

// 一般の関数

func carFunc(a []Any) Any {
	CheckArity(1, a)
	return a[0].(*Cell).Car
}

func cdrFunc(a []Any) Any {
	CheckArity(1, a)
	return a[0].(*Cell).Cdr
}

func consFunc(a []Any) Any {
	CheckArity(2, a)
	return Cons(a[0], a[1].(*Cell))
}

func listpFunc(a []Any) Any {
	CheckArity(1, a)
	_, ok := a[0].(*Cell)
	return LispBool(ok)
}

func eqFunc(a []Any) Any {
	CheckArity(2, a)
	return LispBool(a[0] == a[1])
}

func rplacaFunc(a []Any) Any {
	CheckArity(2, a)
	x := a[0].(*Cell)
	y := a[1]
	x.Car = y
	return y
}

func rplacdFunc(a []Any) Any {
	CheckArity(2, a)
	x := a[0].(*Cell)
	y := a[1].(*Cell)
	x.Cdr = y
	return y
}

func listFunc(a []Any) Any {
	var x *Cell = nil
	for i := len(a) - 1; i > -1; i-- {
		x = Cons(a[i], x)
	}
	return x
}

func eqOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) == 0)
}

func neOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) != 0)
}

func ltOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) < 0)
}

func leOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) <= 0)
}

func gtOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) > 0)
}

func geOp(a []Any) Any {
	CheckArity(2, a)
	return LispBool(arith.Compare(a[0], a[1]) >= 0)
}

func addOp(a []Any) Any {
	return inject(0, a, arith.Add)
}

func subtractOp(a []Any) Any {
	n := len(a)
	if n == 0 {
		return 0
	} else if n == 1 {
		return arith.Subtract(0, a[0])
	}
	return inject(a[0], a[1:], arith.Subtract)
}

func multiplyOp(a []Any) Any {
	return inject(1, a, arith.Multiply)
}

func divideOp(a []Any) Any {
	CheckArity(-2, a)
	return inject(a[0], a[1:], arith.DivideReal)
}

var gensymCount = 0

func gensymFunc(a []Any) Any {
	CheckArity(0, a)
	gensymCount++
	return NewSymbol(fmt.Sprintf("G%05d", gensymCount))
}

func printFunc(a []Any) Any {
	CheckArity(1, a)
	fmt.Println(StringFor(a[0]))
	return a[0]
}

// スペシャル・フォーム

// (quote expression)
func quoteForm(x *Cell, env *Env) (Any, *Env) {
	a := CheckForUnary(x)
	return a, nil
}

// (setq variable expression)
func setqForm(x *Cell, env *Env) (Any, *Env) {
	a, b := CheckForBinary(x)
	sym := a.(*Symbol)
	val := env.Eval(b)
	env.Set(sym, val)
	return val, nil
}

// (progn expression ...)
func prognForm(x *Cell, env *Env) (Any, *Env) {
	if x == nil {
		return (*Cell)(nil), nil
	}
	for ; x.Cdr != nil; x = x.Cdr {
		env.Eval(x.Car)
	}
	return x.Car, env
}

// (if condition then-expression [else-expression ...])
func ifForm(x *Cell, env *Env) (Any, *Env) {
	a, b, c := CheckForBinaryAndRest(x)
	if env.Eval(a) != (*Cell)(nil) {
		return b, env
	}
	return prognForm(c, env)
}

// (lambda ([variable...]) expression...)
func lambdaForm(x *Cell, lambdaEnv *Env) (Any, *Env) {
	a, b := CheckForUnaryAndRest(x)
	params, ok := a.(*Cell)
	if !ok {
		panic(fmt.Errorf("parameter list expected: %s", a))
	}
	return func(args *Cell, argsEnv *Env) (Any, *Env) {
		table := makeArgTable(params, args, argsEnv)
		return prognForm(b, &Env{table, lambdaEnv})
	}, nil
}

// (let ([var|(var expression)...]) expression...)
func letForm(x *Cell, env *Env) (Any, *Env) {
	a, b := CheckForUnaryAndRest(x)
	table := make(map[*Symbol]Any)
	for vars := a.(*Cell); vars != nil; vars = vars.Cdr {
		switch v := vars.Car.(type) {
		case *Symbol:
			table[v] = nil
		case *Cell:
			name, exp := CheckForBinary(v)
			table[name.(*Symbol)] = env.Eval(exp)
		default:
			panic(fmt.Errorf("symbol or (symbol expession) expected: %s",
				StringFor(v)))
		}
	}
	return prognForm(b, &Env{table, env})
}

// (defun name ([variable...]) expession...)
func defunForm(x *Cell, env *Env) (Any, *Env) {
	a, b := CheckForUnaryAndRest(x)
	sym := a.(*Symbol)
	lambda, _ := lambdaForm(b, env)
	env.Set(sym, lambda)
	return sym, nil
}

// (apply expession expession)
func applyForm(x *Cell, env *Env) (Any, *Env) {
	a, b := CheckForBinary(x)
	a = env.Eval(a)
	q := env.Eval(b).(*Cell)
	q = quoteList(q)
	return Cons(a, q), env
}

// (and expession...)
func andForm(x *Cell, env *Env) (Any, *Env) {
	if x == nil {
		return TSymbol, nil
	}
	for ; x.Cdr != nil; x = x.Cdr {
		if env.Eval(x.Car) == (*Cell)(nil) {
			return (*Cell)(nil), nil
		}
	}
	return x.Car, env
}

// 各種ユーティリティ

// 引数の個数検査. 負数の arity は |arity| 個以上の引数を意味する。
func CheckArity(arity int, a []Any) {
	n := len(a)
	if arity < 0 {
		if n < -arity {
			panic(fmt.Errorf("arity %d+; given %d", -arity, n))
		}
	} else if n != arity {
		panic(fmt.Errorf("arity %d; given %d", arity, n))
	}
}

// 真ならばシンボル t を，偽ならば空リストを返す。
func LispBool(t bool) Any {
	if t {
		return TSymbol
	}
	return (*Cell)(nil)
}

func inject(x arith.Number, a []Any,
	op func(arith.Number, arith.Number) arith.Number) arith.Number {
	for _, y := range a {
		x = op(x, y)
	}
	return x
}

// 長さ１のリストか確かめてその要素を返す。
func CheckForUnary(x *Cell) Any {
	if x != nil {
		a := x.Car
		if x.Cdr == nil {
			return a
		}
	}
	panic(fmt.Errorf("arity 1; given %s", StringFor(x)))
}

// 長さ１以上のリストか確かめてその要素を返す。
func CheckForUnaryAndRest(x *Cell) (Any, *Cell) {
	if x != nil {
		a := x.Car
		return a, x.Cdr
	}
	panic(fmt.Errorf("arity 1+; given %s", StringFor(x)))
}

// 長さ２のリストか確かめてその要素を返す。
func CheckForBinary(x *Cell) (Any, Any) {
	if x != nil {
		a := x.Car
		y := x.Cdr
		if y != nil {
			b := y.Car
			if y.Cdr == nil {
				return a, b
			}
		}
	}
	panic(fmt.Errorf("arity 2; given %s", StringFor(x)))
}

// 長さ２以上のリストか確かめてその要素を返す。
func CheckForBinaryAndRest(x *Cell) (Any, Any, *Cell) {
	if x != nil {
		a := x.Car
		y := x.Cdr
		if y != nil {
			b := y.Car
			return a, b, y.Cdr
		}
	}
	panic(fmt.Errorf("arity 2+; given %s", StringFor(x)))
}

// 仮引数名から実引数値への対応表を作る。
func makeArgTable(params *Cell, args *Cell, env *Env) map[*Symbol]Any {
	table := make(map[*Symbol]Any)
	for ; params != nil; params = params.Cdr {
		sym := params.Car.(*Symbol)
		if sym == AmpRestSymbol {
			sym = CheckForUnary(params.Cdr).(*Symbol)
			table[sym] = evlis(args, env)
			return table
		}
		if args == nil {
			panic(fmt.Errorf("missing argument for: %s", StringFor(sym)))
		}
		val := env.Eval(args.Car)
		table[sym] = val
		args = args.Cdr
	}
	if args != nil {
		panic(fmt.Errorf("unused rest: %s", StringFor(args)))
	}
	return table
}

// 剰余引数を評価したリストを作る。
func evlis(args *Cell, env *Env) *Cell {
	if args == nil {
		return nil
	}
	val := env.Eval(args.Car)
	return Cons(val, evlis(args.Cdr, env))
}

// 各要素をクォートしたリストを作る。(apply での二重評価を避けるため)
func quoteList(x *Cell) *Cell {
	if x == nil {
		return nil
	}
	val := Cons(QuoteSymbol, Cons(x.Car, nil))
	return Cons(val, quoteList(x.Cdr))
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

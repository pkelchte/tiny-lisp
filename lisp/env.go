// H25.3/18 - 4/15 (鈴)

// このファイルは Lisp の環境と評価器を実装する。

package lisp

import (
	"fmt"
	"sync"
)

// 環境. 現在の束縛変数に対する Table と自由変数に対する Next からなる。
// Table を参照するときは Lock で排他すること。
type Env struct {
	Table map[*Symbol]Any
	Next  *Env
	Lock  sync.Mutex
}

// シンボルに対する値を環境から得る。無ければパニックする。
func (env *Env) Get(sym *Symbol) Any {
	for ; env != nil; env = env.Next {
		env.Lock.Lock()
		val, ok := env.Table[sym]
		env.Lock.Unlock()
		if ok {
			return val
		}
	}
	panic(fmt.Errorf("unbound symbol: %s", sym.string))
}

// シンボルに対する値を環境にセットする。
// 未定義のシンボルをトップレベル以外でセットするとパニックする。
func (env *Env) Set(sym *Symbol, val Any) {
	for ev := env; ev != nil; ev = ev.Next {
		ev.Lock.Lock()
		_, ok := ev.Table[sym]
		if ok {
			ev.Table[sym] = val
			ev.Lock.Unlock()
			return
		}
		ev.Lock.Unlock()
	}
	if env.Next == nil { // トップレベルの環境ならば
		env.Lock.Lock()
		env.Table[sym] = val
		env.Lock.Unlock()
		return
	}
	panic(fmt.Errorf("global symbol created locally: %s", sym.string))
}

// 与えられた環境のもとで引数を評価する。
func (env *Env) Eval(a Any) Any {
	for {
		switch x := a.(type) {
		case *Cell:
			if x == nil {
				return (*Cell)(nil)
			}
			switch fn := env.Eval(x.Car).(type) {
			case (func(*Cell, *Env) (Any, *Env)): // スペシャル・フォーム
				a, env = fn(x.Cdr, env) // 末尾式ならば環境と共に式を返す。
				if env == nil {
					return a
				}
				// 次回のループで末尾式 a を環境 env (≠nil)で評価する。
			case (func([]Any) Any): // 一般の関数
				arg := make([]Any, 0, 4)
				for y := x.Cdr; y != nil; y = y.Cdr {
					e := env.Eval(y.Car)
					arg = append(arg, e)
				}
				return fn(arg)
			default:
				panic(fmt.Errorf("not function: %s", StringFor(x.Car)))
			}
		case *Symbol:
			return env.Get(x)
		default:
			return x
		}
	}
	panic(fmt.Errorf("never: %s", StringFor(a)))
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

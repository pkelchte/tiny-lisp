// H25.3/18 - 4/1 (鈴)

// このファイルは Lisp の字句解析と構文解析を実装する。

package lisp

import (
	"github.com/pkelchte/tiny-lisp/arith"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

// 字句解析器 (Lexical analyzer)
type Lex struct {
	scanner.Scanner
	Token rune // 現在のトークン
	Value Any  // 現在のトークンの値
}

// 入力ソースに対する字句解析器を返す。
func NewLex(src io.Reader) *Lex {
	var lex Lex
	lex.Init(src)
	lex.Mode &^= scanner.ScanChars | scanner.ScanRawStrings
	lex.NextToken()
	return &lex
}

// 文脈情報を伴った error でパニックを発生させる。
func (lex *Lex) Panic(msg string) {
	panic(fmt.Errorf("%s: %q", msg, lex.TokenText()))
}

// 次のトークンを lex.Token に得る。
// 識別子や数や文字列ならば値を lex.Value をセットする。
// 一般の記号類はすべて識別子扱いする。
func (lex *Lex) NextToken() {
	var token rune
	var text string
	for { // ; から行末までをコメントとして無視する
		token = lex.Scan()
		if token != ';' {
			break
		}
		for {
			r := lex.Next()
			if r == scanner.EOF {
				lex.Token = scanner.EOF
				return
			} else if r == '\n' {
				break
			}
		}
	}
	switch token {
	case '(', ')', '\'', scanner.EOF:
		lex.Token = token
		return
	case scanner.Int, scanner.Float:
		lex.Token = token
		lex.Value, text = parseNumber(lex)
		if text == "" {
			return
		}
	case '-': // -数字... は負数として扱う
		d := lex.Peek()
		if '0' <= d && d <= '9' {
			lex.Token = lex.Scan()
			lex.Value, text = parseNumber(lex)
			if text == "" {
				lex.Value = arith.Subtract(0, lex.Value)
				return
			}
			text = "-" + text
		} else {
			text = "-"
		}
	case scanner.String:
		text = lex.TokenText()
		s, err := strconv.Unquote(text)
		if err != nil {
			lex.Panic("invalid string")
		}
		lex.Value = s
		lex.Token = scanner.String
		return
	default:
		text = lex.TokenText()
	}
	for {
		r, ok := peekAndTest(lex)
		if ok {
			break
		}
		text = fmt.Sprintf("%s%c", text, r)
		lex.Next()
	}
	lex.Value = NewSymbol(text)
	lex.Token = scanner.Ident
	return
}

// 次の文字を見てそこが単語の切れ目かどうかテストする。
func peekAndTest(lex *Lex) (rune, bool) {
	r := lex.Peek()
	return r, (unicode.IsSpace(r) || strings.ContainsRune("()';.,", r) ||
		r == scanner.EOF)
}

func parseNumber(lex *Lex) (arith.Number, string) {
	text := lex.TokenText()
	num, ok := NumberFor(text)
	if ok {
		r, ok := peekAndTest(lex)
		if ok {
			return num, ""
		}
		if r == '/' && lex.Token == scanner.Int { // 有理数リテラル 整数/整数
			lex.Next()
			t := lex.Scan()
			text2 := lex.TokenText()
			if t == scanner.Int {
				num2, ok := NumberFor(text2)
				if ok {
					r, ok = peekAndTest(lex)
					if ok {
						return arith.DivideReal(num, num2), ""
					}
				}
			}
			return 0, text + "/" + text2
		}
	}
	return 0, text
}

// 字句解析器を使ってパースし Lisp 式を返す。
// 空リストは *Cell 型の nil (つまり Any 型の非 nil) として返す。
// EOF ならば Any 型の nil を返す。
func (lex *Lex) Read() Any {
	switch lex.Token {
	case ')':
		lex.Panic("')' unexpected")
	case '\'':
		lex.NextToken()
		return Cons(QuoteSymbol, Cons(lex.Read(), nil))
	case '(':
		lex.NextToken()
		x, ok := parseListBody(lex)
		if !ok {
			return nil
		}
		return x
	case scanner.EOF:
		return nil
	}
	value := lex.Value
	lex.NextToken()
	if value == NilSymbol { // 識別子 nil は空リストとして扱う
		return (*Cell)(nil)
	}
	return value
}

func parseListBody(lex *Lex) (*Cell, bool) {
	if lex.Token == ')' {
		lex.NextToken()
		return nil, true
	}
	var e1 Any = lex.Read()
	if e1 == nil {
		return nil, false
	}
	e2, ok := parseListBody(lex)
	return Cons(e1, e2), ok
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

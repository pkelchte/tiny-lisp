// H25.3/18 - 3/28 (鈴)
package main

import (
	"bufio"
	"fmt"
	"github.com/pkelchte/tiny-lisp/lisp"
	"os"
	"strings"
)

// 初期化スクリプト
var prelude = `
(defun null (x) (eq x nil))
(defun not (x) (eq x nil))

(defun length (x)
  (if (null x)
      0
    (+ 1 (length (cdr x)))))

(defun append (&rest x)
  (if (null x)
      nil
    (if (null (cdr x))
        (car x)
      (_append (car x) (apply append (cdr x))))))

(defun _append (x y)
  (if (null x)
      y
    (cons (car x) (_append (cdr x) y))))
`

// 文字列を読み込み式を評価して結果を表示する。
// 式が不完全ならば done 引数に false を与えて終わる。
func readEvalPrint(line string, done *bool) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				fmt.Printf("==> %s\n", e.Error())
			case string:
				fmt.Printf("===> %s\n", e)
			default:
				panic(e)
			}
		}
	}()
	*done = true
	if !lisp.ReadEvalPrint(strings.NewReader(line), os.Stdout) {
		*done = false
	}
}

// Lisp スクリプトまたは Lisp 対話セッションを実行する。
func main() {
	lisp.ReadAndEval(prelude)
	n := len(os.Args)
	if n >= 2 && os.Args[1] != "-" {
		lisp.ReadAndEvalFile(os.Args[1])
	}
	if n < 2 || os.Args[n-1] == "-" {
		// 対話セッションを始める。
		var rd *bufio.Reader = bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			line := ""
			for {
				s, err := rd.ReadString('\n')
				if err != nil {
					return
				}
				line += s
				var done bool
				readEvalPrint(line, &done)
				if done {
					break
				}
				fmt.Print("  ")
			}
		}
	}
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

// H25.3/18 - 3/28 (鈴)

// 簡単な Lisp インタープリタ.
// 数として有理数を扱う。
package lisp

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

// 文字列を読み込み式を評価する。
// ただし，読み込んだ式が不完全ならば false を返して終わる。
func ReadAndEval(line string) bool {
	var src io.Reader = strings.NewReader(line)
	lex := NewLex(src)
	for lex.Token != scanner.EOF {
		x := lex.Read()
		if x == nil {
			return false
		}
		Globals.Eval(x)
	}
	return true
}

// ファイルを読み込み式を評価する。
// ただし，読み込んだ式が不完全ならば false を返して終わる。
func ReadAndEvalFile(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	} else {
		defer func() {
			file.Close()
		}()
		lex := NewLex(file)
		for lex.Token != scanner.EOF {
			x := lex.Read()
			if x == nil {
				return false
			}
			Globals.Eval(x)
		}
	}
	return true
}

// 入力を読み込み式を評価し結果を 元の式 => 結果の値 という形式で出力する。
// ただし，読み込んだ式が不完全ならば false を返して終わる。
func ReadEvalPrint(input io.Reader, output io.Writer) bool {
	lex := NewLex(input)
	for lex.Token != scanner.EOF {
		x := lex.Read()
		if x == nil {
			return false
		}
		fmt.Fprintf(output, "%v => ", StringFor(x))
		y := Globals.Eval(x)
		fmt.Fprintf(output, "%v\n", StringFor(y))
	}
	return true
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

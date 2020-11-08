/*
Command assembly is an interpreter for the Assembly programming language.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dkmccandless/assembly/eval"
	"github.com/dkmccandless/assembly/lexer"
	"github.com/dkmccandless/assembly/object"
	"github.com/dkmccandless/assembly/parser"
)

func main() {
	helpmsg := `Command assembly is an interpreter for the Assembly programming language.

Usage:	assembly [resolution name]
`
	if len(os.Args) == 1 || strings.ToLower(os.Args[1]) == "help" {
		fmt.Println(helpmsg)
		return
	}
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ast, err := parser.New(lexer.New(string(b))).ParseResolution()
	if err != nil {
		fmt.Println(err)
		return
	}
	if obj := eval.Eval(ast, object.NewEnvironment()); err != nil {
		fmt.Println(obj.Inspect())
	}
}

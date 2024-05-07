package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/d2verb/mankey/evaluator"
	"github.com/d2verb/mankey/lexer"
	"github.com/d2verb/mankey/object"
	"github.com/d2verb/mankey/parser"
	"github.com/d2verb/mankey/repl"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Mankey %s\n", VERSION)
		repl.Start()
	} else {
		content, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		l := lexer.New(string(content))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				fmt.Println(msg)
			}
			os.Exit(1)
		}

		abspath, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		curdir := filepath.Dir(abspath)

		env := object.NewEnvironmentWithDir(curdir)
		evaluated := evaluator.Eval(program, env)

		if evaluated.Type() == object.INTEGER_OBJ {
			os.Exit(int(evaluated.(*object.Integer).Value))
		} else {
			os.Exit(0)
		}
	}
}

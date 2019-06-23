package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/d2verb/monkey/evaluator"
	"github.com/d2verb/monkey/lexer"
	"github.com/d2verb/monkey/object"
	"github.com/d2verb/monkey/parser"
	"github.com/d2verb/monkey/repl"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Monkey %s\n", VERSION)
		repl.Start(os.Stdin, os.Stdout)
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

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		if evaluated.Type() == object.INTEGER_OBJ {
			os.Exit(int(evaluated.(*object.Integer).Value))
		} else {
			os.Exit(0)
		}
	}
}

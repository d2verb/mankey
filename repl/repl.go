package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/d2verb/mankey/evaluator"
	"github.com/d2verb/mankey/lexer"
	"github.com/d2verb/mankey/object"
	"github.com/d2verb/mankey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	curdir, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	env := object.NewEnvironmentWithDir(curdir)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "Syntax error: "+msg+"\n")
	}
}

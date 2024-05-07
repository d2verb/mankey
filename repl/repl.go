package repl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c-bata/go-prompt"
	"github.com/d2verb/mankey/evaluator"
	"github.com/d2verb/mankey/lexer"
	"github.com/d2verb/mankey/object"
	"github.com/d2verb/mankey/parser"
)

const PROMPT = ">> "

func Start() {
	p := prompt.New(executor, completer, prompt.OptionPrefix(PROMPT))
	p.Run()
}

func executor(in string) {
	curdir, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	env := object.NewEnvironmentWithDir(curdir)

	l := lexer.New(in)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(os.Stderr, "Syntax error: "+msg+"\n")
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

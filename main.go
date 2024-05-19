package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/d2verb/mankey/evaluator"
	"github.com/d2verb/mankey/lexer"
	"github.com/d2verb/mankey/lsp"
	"github.com/d2verb/mankey/object"
	"github.com/d2verb/mankey/parser"
	"github.com/d2verb/mankey/repl"
	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mankey",
		Short: "mankey programming language",
		Run: func(cmd *cobra.Command, args []string) {
			runRepl()
		},
	}

	var runCodeCmd = &cobra.Command{
		Use:   "run <FILEPATH>",
		Short: "Read and execute mankey source file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires <FILEPATH>")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCode(args[0])
		},
	}

	var lspCmd = &cobra.Command{
		Use:   "lsp",
		Short: "Start Language Server",
		Run: func(cmd *cobra.Command, args []string) {
			runLsp()
		},
	}

	rootCmd.AddCommand(runCodeCmd)
	rootCmd.AddCommand(lspCmd)
	rootCmd.Execute()
}

func runRepl() {
	fmt.Printf("Mankey %s\n", VERSION)
	repl.Start()
}

func runCode(filename string) {
	content, err := os.ReadFile(filename)
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
	} else if evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func runLsp() {
	lsp.Serve()
}

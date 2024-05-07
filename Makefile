all:
	go build -o mankey main.go

install:
	go install github.com/d2verb/mankey

test:
	go test ./lexer
	go test ./ast
	go test ./parser
	go test ./evaluator

clean:
	rm -f mankey

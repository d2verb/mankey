all:
	go build -o monkey main.go

install:
	go install github.com/d2verb/monkey

test:
	go test ./lexer
	go test ./ast
	go test ./parser
	go test ./evaluator

clean:
	rm -f monkey

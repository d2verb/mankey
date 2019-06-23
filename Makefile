all:
	go build -o monkey main.go

install:
	go install github.com/d2verb/monkey

clean:
	rm -f monkey

package lsp

import (
	"context"
	"io"
	"log"
	"os"

	"go.lsp.dev/jsonrpc2"
)

type StdioReadWriteCloser struct {
	io.Reader
	io.Writer
}

func (rwc *StdioReadWriteCloser) Close() error {
	return nil
}

func Serve() {
	rwc := &StdioReadWriteCloser{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}

	stream := jsonrpc2.NewStream(rwc)
	conn := jsonrpc2.NewConn(stream)
	ctx := context.TODO()

	conn.Go(ctx, serverHandler)
	<-conn.Done()

	if err := conn.Err(); err != nil {
		log.Fatalf("LSP server exited with error: %v", err)
	}
}

func serverHandler(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	return nil
}

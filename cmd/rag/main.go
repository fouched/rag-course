package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rag-course/app"
	"rag-course/config"
	"syscall"
)

func main() {
	// We need to:
	// - set up the app
	// - set up config
	// - set up an LLM client
	// - set up the Read-Eval-Print loop (REPL)

	// shutdown cleanly
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Calling app.Run will actually call chat.RunREPL, which runs the Read-Eval-Print Loop for our llm chat.
	if err := app.Run(ctx, config.Load()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"go-rdbms/repl"
)

func main() {
	fmt.Println("Simple RDBMS - Type 'help' for commands, 'exit' to quit")

	repl, err := repl.NewRepl("./data")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	if err := repl.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

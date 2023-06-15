package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/replee/repl"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type S struct {
	Field int `json:"field"`
}

func (s *S) GetField() int {
	return s.Field
}

func main() {
	ctx := context.Background()
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := repl.BootstrapRepleeVM(ctx, vm); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	// Create a channel to receive OS signals
	sig := make(chan os.Signal, 1)
	// Notify the signal channel for SIGINT
	signal.Notify(sig, syscall.SIGINT)

	// Run a goroutine that waits for the SIGINT signal
	go func() {
		<-sig
		fmt.Println("\nGracefully shutting down...")
		os.Exit(0)
	}()

	for scanner.Scan() {
		input := scanner.Text()
		value, err := vm.RunString(input)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(value)
		}

		fmt.Print("> ")
	}

	if scanner.Err() != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
}

//vm := goja.New()
//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
//
//vm.Set("s", S{Field: 42})
//res, _ := vm.RunString(`s.field`) // without the mapper it would have been s.Field
//fmt.Println(res.Export())

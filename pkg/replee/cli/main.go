package main

import (
	"encoding/json"
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/handler"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: replee <action json>")
	}
	a := action.Action{}
	if err := json.Unmarshal([]byte(os.Args[1]), &a); err != nil {
		log.Fatal(err)
	}
	s, err := store.NewFSResolutionProblemStore(".")
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewRepleeHandler()
	if err := h.HandleAction(a, s); err != nil {
		log.Fatal(err)
	}
}

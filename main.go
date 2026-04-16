package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "kai: couldn't load config —", err)
		os.Exit(1)
	}

	mem := LoadMemory()
	args := os.Args[1:]

	// no args → interactive chat
	if len(args) == 0 {
		RunChat(cfg, mem)
		return
	}

	switch args[0] {
	case "--startup":
		RunStartup(cfg, mem)

	case "--context-check":
		RunContextCheck(cfg, mem)

	case "--error":
		code, cmd := parseErrorFlags(args[1:])
		RunError(cfg, code, cmd)

	default:
		// treat anything else as a quick question
		question := strings.Join(args, " ")
		RunQuickAsk(cfg, mem, question)
	}
}

func parseErrorFlags(args []string) (code string, cmd string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--code":
			if i+1 < len(args) {
				code = args[i+1]
				i++
			}
		case "--cmd":
			if i+1 < len(args) {
				cmd = args[i+1]
				i++
			}
		}
	}
	return
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunChat(cfg *Config, mem *Memory) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("kai: hey, what's up? (ctrl+c to exit)")

	for {
		fmt.Print("you: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		msgs := buildMessages(cfg, mem.ChatHistory, input)
		reply, err := Ask(cfg, msgs)
		if err != nil {
			fmt.Println("kai: something went wrong —", err)
			continue
		}

		fmt.Println("kai:", reply)

		mem.AddMessage("user", input, cfg.ChatHistoryLimit)
		mem.AddMessage("assistant", reply, cfg.ChatHistoryLimit)
		mem.Save()
	}
}

func RunQuickAsk(cfg *Config, mem *Memory, question string) {
	msgs := buildMessages(cfg, mem.ChatHistory, question)
	reply, err := Ask(cfg, msgs)
	if err != nil {
		fmt.Println("kai: something went wrong —", err)
		return
	}

	fmt.Println(reply)

	mem.AddMessage("user", question, cfg.ChatHistoryLimit)
	mem.AddMessage("assistant", reply, cfg.ChatHistoryLimit)
	mem.Save()
}

package main

import (
	"fmt"
	"time"
)

func RunStartup(cfg *Config, mem *Memory) {
	now := time.Now()
	if time.Since(mem.LastStartup) < time.Minute {
		return
	}

	hour := now.Hour()
	var timeContext string
	switch {
	case hour >= 5 && hour < 12:
		timeContext = "It's morning. Give a short energetic greeting."
	case hour >= 12 && hour < 17:
		timeContext = "It's afternoon. Give a chill, casual check-in."
	case hour >= 17 && hour < 23:
		timeContext = "It's evening. Give a relaxed wind-down greeting."
	default:
		timeContext = "It's late night. React like you're surprised they're still up."
	}
	prompt := fmt.Sprintf(
		"You are greeting the user as they open their terminal. %s Keep it under 2 sentences, no emojis.",
		timeContext,
	)
	msgs := buildMessages(cfg, nil, prompt)
	reply, err := Ask(cfg, msgs)
	if err != nil {
		fmt.Println("Error during startup greeting:", err)
		return
	}

	fmt.Println(reply)
	mem.LastStartup = now
	mem.Save()
}

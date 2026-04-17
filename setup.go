package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NeedsSetup(cfg *Config) bool {
	return cfg == nil || cfg.APIKey == "" || cfg.APIKey == "your_groq_key_here"
}

func RunSetup() (*Config, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│   ai-cli — first time setup          │")
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()

	// companion name
	fmt.Print("what do you want to call your companion? (default: Kai) → ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = "Kai"
	}

	// personality
	fmt.Println()
	fmt.Printf("describe %s's personality in one sentence.\n", name)
	fmt.Printf("(default: chill but sharp, keeps it short, actually helpful) → ")
	scanner.Scan()
	personality := strings.TrimSpace(scanner.Text())
	if personality == "" {
		personality = "chill but sharp, keeps it short, actually helpful"
	}

	// api key
	fmt.Println()
	fmt.Println("get your free Groq API key at → https://console.groq.com")
	fmt.Print("paste your Groq API key → ")
	scanner.Scan()
	apiKey := strings.TrimSpace(scanner.Text())
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	cfg := &Config{
		Provider:               "groq",
		APIKey:                 apiKey,
		Model:                  "llama-3.3-70b-versatile",
		Personality:            fmt.Sprintf("You are %s, a terminal companion. You are %s. Respond in plain text only, no markdown.", name, personality),
		ContextCooldownMinutes: 30,
		ChatHistoryLimit:       10,
	}

	if err := saveConfig(cfg); err != nil {
		return nil, fmt.Errorf("couldn't save config: %w", err)
	}

	fmt.Println()
	fmt.Printf("all set! %s is ready. try: ai \"hey\"\n", name)
	fmt.Println()

	return cfg, nil
}

func saveConfig(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".ai-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Provider               string `json:"provider"`
	APIKey                 string `json:"api_key"`
	Model                  string `json:"model"`
	Personality            string `json:"personality"`
	ContextCooldownMinutes int    `json:"context_cooldown_minutes"`
	ChatHistoryLimit       int    `json:"chat_history_limit"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".ai-cli", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

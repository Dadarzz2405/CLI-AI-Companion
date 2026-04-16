package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Memory struct {
	LastStartup      time.Time            `json:"last_startup"`
	ContextCooldowns map[string]time.Time `json:"context_cooldowns"`
	ChatHistory      []Message            `json:"chat_history"`
}

func memoryPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ai-cli", "memory.json")
}

func LoadMemory() *Memory {
	data, err := os.ReadFile(memoryPath())
	if err != nil {
		// no memory file yet, start fresh
		return &Memory{
			ContextCooldowns: make(map[string]time.Time),
			ChatHistory:      []Message{},
		}
	}

	var m Memory
	if err := json.Unmarshal(data, &m); err != nil {
		return &Memory{
			ContextCooldowns: make(map[string]time.Time),
			ChatHistory:      []Message{},
		}
	}

	return &m
}

func (m *Memory) Save() error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(memoryPath(), data, 0644)
}

func (m *Memory) AddMessage(role, content string, limit int) {
	m.ChatHistory = append(m.ChatHistory, Message{Role: role, Content: content})
	if len(m.ChatHistory) > limit {
		m.ChatHistory = m.ChatHistory[len(m.ChatHistory)-limit:]
	}
}

func (m *Memory) IsOnCooldown(dir string, minutes int) bool {
	last, ok := m.ContextCooldowns[dir]
	if !ok {
		return false
	}
	return time.Since(last) < time.Duration(minutes)*time.Minute
}

func (m *Memory) SetCooldown(dir string) {
	m.ContextCooldowns[dir] = time.Now()
}

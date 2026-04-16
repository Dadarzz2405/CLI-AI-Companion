package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"

type groqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type groqResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func Ask(cfg *Config, messages []Message) (string, error) {
	body, err := json.Marshal(groqRequest{
		Model:    cfg.Model,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", groqURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq error %d: %s", resp.StatusCode, string(data))
	}

	var result groqResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from groq")
	}

	return result.Choices[0].Message.Content, nil
}

func buildMessages(cfg *Config, history []Message, userPrompt string) []Message {
	msgs := []Message{
		{Role: "system", Content: cfg.Personality},
	}
	msgs = append(msgs, history...)
	msgs = append(msgs, Message{Role: "user", Content: userPrompt})
	return msgs
}

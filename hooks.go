package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var projectRootFiles = []string{
	"go.mod",
	"package.json",
	"requirements.txt",
	"pyproject.toml",
	"Cargo.toml",
	"pom.xml",
	"build.gradle",
}

func RunError(cfg *Config, exitCode string, cmd string) {
	prompt := fmt.Sprintf(
		"The user ran this command in their terminal: `%s`\n"+
			"It failed with exit code %s.\n"+
			"Briefly explain what likely went wrong and suggest a fix. Max 3 lines.",
		cmd, exitCode,
	)

	msgs := buildMessages(cfg, nil, prompt)
	reply, err := Ask(cfg, msgs)
	if err != nil {
		return // fail silently
	}

	fmt.Println(reply)
}

func RunContextCheck(cfg *Config, mem *Memory) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	// detect project root
	projectType := detectProject(cwd)
	if projectType == "" {
		return // not a project root, stay quiet
	}

	// check cooldown
	if mem.IsOnCooldown(cwd, cfg.ContextCooldownMinutes) {
		return
	}

	prompt := fmt.Sprintf(
		"The user just cd'd into a %s project at %s. "+
			"Give a short, casual acknowledgement — like you noticed. Max 1 sentence.",
		projectType, cwd,
	)

	msgs := buildMessages(cfg, nil, prompt)
	reply, err := Ask(cfg, msgs)
	if err != nil {
		return
	}

	fmt.Println(reply)

	mem.SetCooldown(cwd)
	mem.Save()
}

func detectProject(dir string) string {
	indicators := map[string]string{
		"go.mod":           "Go",
		"package.json":     "Node",
		"requirements.txt": "Python",
		"pyproject.toml":   "Python",
		"Cargo.toml":       "Rust",
		"pom.xml":          "Java",
		"build.gradle":     "Java",
		"Gemfile":          "Ruby",
	}

	for file, lang := range indicators {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return lang
		}
	}

	return ""
}

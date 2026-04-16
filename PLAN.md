# PLAN.md — Terminal AI Companion (V1)

## Goal
A lightweight AI companion that lives in your terminal. Reacts to errors, greets
you on startup, and notices when you enter a dev project — powered by Groq.

---

## Architecture

```
/ai-cli
  ├── main.go        — CLI entrypoint, parses flags/args, routes to handlers
  ├── chat.go        — interactive + quick-ask chat mode
  ├── startup.go     — --startup greeting logic (time-aware)
  ├── hooks.go       — --error and --context-check handlers
  ├── memory.go      — simple JSON-based session/short-term memory
  ├── groq.go        — Groq API client (raw net/http, no SDK)
  ├── config.go      — loads config.json, exposes settings
  └── config.json    — user config (API key, personality, cooldown, etc.)
```

---

## Command Interface

| command | args | description |
|---|---|---|
| `ai` | — | interactive chat session |
| `ai "question"` | string | single quick question |
| `ai --startup` | — | time-based greeting (called from .zshrc) |
| `ai --error` | `--code $?  --cmd "last cmd"` | reacts to a failed command |
| `ai --context-check` | — | scans cwd, reacts if project root found |

---

## Trigger System (.zshrc)

```zsh
# on terminal open
fastfetch
ai --startup

# capture last command for error hook
preexec() { _AI_LAST_CMD="$1" }

# on every prompt, check exit code
precmd() {
  local code=$?
  if [ $code -ne 0 ]; then
    ai --error --code $code --cmd "$_AI_LAST_CMD"
  fi
}

# on directory change
chpwd() {
  ai --context-check
}
```

---

## Reactive Triggers

### 1. Error Hook (`--error`)
- fires when exit code ≠ 0
- sends: exit code + the failed command string
- AI responds with: what likely went wrong + suggested fix
- always fires (no cooldown — errors are always relevant)

### 2. Context Check (`--context-check`)
- fires on every `cd`
- scans cwd for project root indicators:
  - `go.mod` → Go project
  - `requirements.txt` / `pyproject.toml` → Python project
  - `package.json` → Node project
  - `Cargo.toml` → Rust project
  - `pom.xml` / `build.gradle` → Java project
- only reacts if a project root file is found
- **cooldown**: tracks last-reacted directory + timestamp in a local state file
  - default cooldown: 30 minutes per directory
  - won't re-react to same dir within cooldown window

### 3. Startup Greeting (`--startup`)
- reads current hour
- morning (5–11): energetic greeting
- afternoon (12–17): chill check-in
- evening (18–22): wind-down vibe
- night (23–4): surprised you're up
- randomizes within each bucket so it doesn't feel repetitive

---

## Groq Integration (`groq.go`)

- hits `https://api.groq.com/openai/v1/chat/completions` directly via `net/http`
- model: `llama-3.3-70b-versatile` (fast, free tier)
- system prompt loaded from `config.json` (personality layer)
- each mode (error, context, chat, startup) injects its own context into the prompt
- streaming optional for V1 (non-streaming is simpler, still fast on Groq)

---

## Memory (`memory.go`)

Flat JSON file at `~/.ai-cli/memory.json`

```json
{
  "last_startup": "2025-04-16T08:00:00Z",
  "context_cooldowns": {
    "/home/user/projects/myapp": "2025-04-16T09:30:00Z"
  },
  "chat_history": [
    { "role": "user", "content": "..." },
    { "role": "assistant", "content": "..." }
  ]
}
```

- `chat_history`: last 10 exchanges only (keeps context window small)
- `context_cooldowns`: dir path → last reacted timestamp
- `last_startup`: prevents double-greeting in same session (edge case)

---

## Config (`config.json`)

```json
{
  "provider": "groq",
  "api_key": "YOUR_GROQ_KEY",
  "model": "llama-3.3-70b-versatile",
  "personality": "You are Kai, a chill but sharp terminal companion. You speak casually, keep responses short, and actually help — no fluff.",
  "context_cooldown_minutes": 30,
  "chat_history_limit": 10
}
```

> personality prompt is fully user-editable — this is where the vibe lives

---

## Future / V2 Ideas
- `--local` flag: swap Groq for a local LLM (Ollama-compatible endpoint)
  - local model reads directory tree directly for better context awareness
- smarter memory: summarize old chat history instead of truncating
- git-aware context: detect branch, recent commits, dirty status
- per-project personality override (`.ai-companion` file in project root)

---

## Build & Install

```bash
cd ai-cli
go build -o ai .
sudo mv ai /usr/local/bin/ai

# add to .zshrc
echo 'fastfetch' >> ~/.zshrc
echo 'ai --startup' >> ~/.zshrc
```
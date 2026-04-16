# ai — Terminal AI Companion

A lightweight AI companion that lives in your terminal.
Reacts to errors, greets you on startup, and notices when you enter a dev project.

Powered by [Groq](https://groq.com) — fast enough to feel instant.

```
 terminal opens
 └── fastfetch (your anime rice)
 └── ai --startup  →  "morning, you've got 3 PRs open lol"

 $ git psuh
 └── ai --error    →  "typo — did you mean `git push`?"

 $ cd ~/projects/myapp
 └── ai --context  →  "Go project detected. last time you were here you were fixing the auth bug"
```

---

## Install

**Prerequisites:** Go 1.21+, a [Groq API key](https://console.groq.com) (free)

```bash
git clone https://github.com/you/ai-cli
cd ai-cli
go build -o ai .
sudo mv ai /usr/local/bin/ai
```

### Setup config

```bash
mkdir -p ~/.ai-cli
cp config.json ~/.ai-cli/config.json
# then edit ~/.ai-cli/config.json and add your Groq key
```

### Hook into your shell (`.zshrc`)

```zsh
# visual layer
fastfetch

# ai companion
ai --startup

# capture commands for error hook
preexec() { _AI_LAST_CMD="$1" }

# error + context hooks
precmd() {
  local code=$?
  [[ $code -ne 0 ]] && ai --error --code $code --cmd "$_AI_LAST_CMD"
}

chpwd() { ai --context-check }
```

---

## Usage

| command | what it does |
|---|---|
| `ai` | open interactive chat |
| `ai "why is my goroutine leaking"` | quick one-shot question |
| `ai --startup` | time-aware greeting (used in .zshrc) |
| `ai --error --code 1 --cmd "git psuh"` | diagnose a failed command |
| `ai --context-check` | scan cwd and react if it's a project root |

---

## Configuration

Edit `~/.ai-cli/config.json`:

```json
{
  "provider": "groq",
  "api_key": "your_key_here",
  "model": "llama-3.3-70b-versatile",
  "personality": "You are Kai, a chill but sharp terminal companion. Keep it short and actually helpful.",
  "context_cooldown_minutes": 30,
  "chat_history_limit": 10
}
```

**`personality`** is the most important field — this is where you define the vibe.
Change the name, tone, style, whatever. It's just a system prompt.

---

## How Reactions Work

### Error hook
Fires when any command exits with a non-zero code.
Sends the exit code + the command you ran to the AI.
Gets back a short diagnosis + suggested fix.

### Context check
Fires on every `cd`.
Looks for project root files (`go.mod`, `package.json`, `requirements.txt`, etc.).
Only reacts if one is found, and only once per 30 minutes per directory (configurable).

### Startup greeting
Reads the current time and picks a tone — morning energy, afternoon chill, late night chaos.
Randomized within each time bucket so it doesn't get stale.

---

## Memory

Stored at `~/.ai-cli/memory.json` — plain JSON, you can read/edit it directly.

- last 10 chat messages kept for context in interactive mode
- cooldown timestamps per directory for context-check
- no cloud sync, no telemetry, everything stays local

---

## Roadmap

- [x] V1: Groq-powered, error + context + startup + chat
- [ ] V2: `--local` flag for Ollama-compatible local LLM
- [ ] V2: local model reads full directory tree for richer context
- [ ] V2: git-aware context (branch, dirty status, recent commits)
- [ ] V2: per-project `.ai-companion` config override

---

## License

MIT
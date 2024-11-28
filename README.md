# Go Shell Ask

A powerful command-line interface tool written in Go for interacting with various Large Language Models (LLMs) directly from your terminal. This is a Go implementation of the [shell-ask](https://github.com/egoist/shell-ask) project.

![Version](https://img.shields.io/badge/version-0.1.0-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 🤖 Support for multiple LLM providers:
  - OpenAI (GPT-3.5, GPT-4)
  - Anthropic (Claude)
  - Google (Gemini)
  - Groq
  - Ollama (local models)
- 🔄 Real-time streaming responses
- 📝 Pipe input support
- 🛠️ Custom commands
- ⚙️ Configuration via config files and environment variables
- 💬 Chat history management
- 🔍 Web search capability
- 🌐 URL content fetching
- 📊 Markdown rendering

## Installation

### Prerequisites

- Go 1.21 or higher
- Git (for installation from source)

### Using go install

```bash
go install github.com/acazau/shell-ask-go/cmd/ask@latest
```

### Building from source

```bash
git clone https://github.com/acazau/shell-ask-go.git
cd shell-ask-go
go build -o ask cmd/ask/main.go
```

## Configuration

Create a configuration file at `~/.config/shell-ask-go/config.json`:

```json
{
  "default_model": "gpt-4",
  "openai_api_key": "your-openai-key",
  "anthropic_api_key": "your-anthropic-key",
  "gemini_api_key": "your-gemini-key",
  "groq_api_key": "your-groq-key",
  "ollama_host": "http://localhost:11434"
}
```

Or use environment variables:
- `SHELL_ASK_OPENAI_API_KEY`
- `SHELL_ASK_ANTHROPIC_API_KEY`
- `SHELL_ASK_GEMINI_API_KEY`
- `SHELL_ASK_GROQ_API_KEY`
- `SHELL_ASK_OLLAMA_HOST`

## Usage

### Basic Usage

```bash
# Ask a question
ask "how to list all docker containers?"

# Get command-only output
ask -c "show me the git log for the last 5 commits"

# Use a specific model
ask -m claude-3 "explain quantum computing"

# Pipe input
cat main.go | ask "explain this code"

# Generate commit message
git diff | ask cm
```

### Command Line Flags

```
Flags:
  -m, --model string       Choose the LLM to use
  -c, --command           Ask LLM to return a command only
  -t, --type string       Define the shape of the response
  -u, --url string        Fetch URL content as context
  -s, --search            Enable web search
      --no-stream         Disable streaming output
  -r, --reply            Reply to previous conversation
  -h, --help             Help for ask
```

### Custom Commands

Define custom commands in your config file:

```json
{
  "commands": [
    {
      "command": "explain",
      "description": "Explain the code in the input",
      "prompt": "Explain the following code:\n{{input}}",
      "require_stdin": true
    }
  ]
}
```

Use custom commands:
```bash
cat main.go | ask explain
```

## Development

### Project Structure

```
shell-ask-go/
├── cmd/
│   └── ask/                 # Application entry point
├── internal/
│   ├── config/             # Configuration handling
│   ├── models/             # Model definitions
│   ├── providers/          # LLM provider implementations
│   ├── commands/           # Command handling
│   └── cli/               # CLI implementation
├── pkg/
│   ├── chat/              # Chat history management
│   ├── markdown/          # Markdown rendering
│   ├── stream/            # Response streaming
│   └── utils/             # Shared utilities
```

### Running Tests

```bash
go test ./...
```

### Building with Version Information

```bash
go build -ldflags "-X github.com/acazau/shell-ask-go/pkg/version.Version=1.0.0 -X github.com/acazau/shell-ask-go/pkg/version.GitCommit=$(git rev-parse HEAD)" ./cmd/ask
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Original [shell-ask](https://github.com/egoist/shell-ask) project by EGOIST
- All the amazing Go libraries used in this project
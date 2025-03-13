# co - AI-Powered Git Commit Message Generator

[![Build and Test](https://github.com/hamzabow/co/actions/workflows/build.yml/badge.svg)](https://github.com/hamzabow/co/actions/workflows/build.yml)
[![Release](https://github.com/hamzabow/co/actions/workflows/release.yml/badge.svg)](https://github.com/hamzabow/co/actions/workflows/release.yml)

`co` is a command-line tool that uses OpenAI's GPT-4o to generate meaningful, well-formatted git commit messages based on your staged changes.

## Features

- 🤖 Generates contextual commit messages based on your staged git changes
- 🔑 Securely manages your OpenAI API key
- ✏️ Interactive editor to review and modify suggested commit messages
- 📋 Supports multiple commit message formats:
  - Conventional Commits
  - Gitmoji (Unicode or shortcode format)
  - Simple format

## Prerequisites

- Git installed and configured
- An OpenAI API key (with access to GPT-4o)

## Installation

You'll need Go (Golang) installed so that you can run the following command to install `co`:

```bash
go install github.com/hamzabow/co@latest
```

Once installed, you can use `co` in any git repository. See the Usage section below for instructions.

## Setup

You can provide your OpenAI API key in one of two ways:

1. Set the `OPENAI_API_KEY` environment variable:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. Enter it when prompted on first use. The tool will ask for your API key if not found in environment variables.

## Usage

1. Stage your changes with git:
   ```bash
   git add .
   ```

2. Run the commit message generator:
   ```bash
   co
   ```

3. Review the generated message, edit if needed, and:
   - Press `Ctrl+Enter` to commit with the message
   - Press `Ctrl+C` to cancel

## How It Works

1. The tool retrieves the diff of your staged changes using `git diff --staged`
2. It sends this diff to OpenAI's API with a carefully crafted prompt
3. The AI generates a commit message following the specified format
4. You get to review and edit the message before committing
5. After confirmation, the tool executes `git commit -m "your message"`

## Configuration

The tool currently uses the Conventional Commits format by default. You can modify the format by editing the `internal/genmessage/genmessage.go` file to use one of the other prompt templates defined in `internal/prompts/prompts.go`.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests for new features, improvements, or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the terminal UI
- Uses the OpenAI API for generating commit messages

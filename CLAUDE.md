# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Build the binary
go build -o speedread .

# Run directly with go run
go run . filename.txt
go run . -wpm 300 filename.txt

# Run from stdin
echo "Hello world" | go run .
```

## Architecture

This is a single-file Go CLI application (`main.go`) that displays text one word at a time using ASCII art in the terminal.

### Key Components

- **ASCII Font Map** (lines 17-326): A `map[rune][]string` defining 5-row tall block characters for A-Z, 0-9, and common punctuation
- **renderWord()**: Scales and centers ASCII art words to fit terminal dimensions, using a reference word length of 8 characters for consistent sizing
- **Terminal handling**: Uses `golang.org/x/term` for raw mode input and terminal size detection; reads keyboard from `/dev/tty` to allow stdin piping while still capturing keypresses
- **Pause system**: Uses `atomic.Bool` for thread-safe pause state toggled by spacebar in a separate goroutine

### Input Flow

1. `readInput()` - reads from file argument or stdin
2. `tokenizeWords()` - splits text using `bufio.ScanWords`
3. Main loop displays each word with calculated delay based on WPM
4. Optional extra pause after punctuation (`.`, `,`, `!`, `?`, `;`, `:`, `"`, `'`)

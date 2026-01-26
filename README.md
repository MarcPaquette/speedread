# speedread

A terminal-based speed reading application that displays text one word at a time using large ASCII art with Spritz-style focal point highlighting.

## Installation

```bash
go build -o speedread .
```

## Usage

```bash
# Read from a file
./speedread filename.txt

# Read from a URL
./speedread https://example.com/article

# Read from stdin
cat filename.txt | ./speedread
echo "Hello, world!" | ./speedread
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-wpm` | Words per minute (10-1000) | 200 |
| `-punct-pause`, `-p` | Extra pause after punctuation in milliseconds | 0 |
| `-focal` | Enable focal point highlighting (Spritz-style) | true |
| `-focal-color`, `-c` | Focal point color (black, red, green, yellow, blue, magenta, cyan, white) | red |
| `-context` | Show surrounding words (previous/next) for context | false |

## Controls

| Key | Action |
|-----|--------|
| `Space` | Pause/unpause |
| `↑` | Increase WPM by 25 |
| `↓` | Decrease WPM by 25 |
| `←` | Rewind one word |
| `→` | Skip forward one word |
| `0-9` | Jump to percentage (0=0%, 1=10%, ..., 9=90%) |
| `Ctrl+C` | Exit (saves bookmark for files) |

## Features

- **Focal point highlighting**: Uses Spritz-style ORP (Optimal Recognition Point) to highlight the focal character in each word
- **Bookmarks**: Automatically saves your position when reading files; resume where you left off
- **Progress display**: Shows current WPM, time remaining, and progress bar
- **Session statistics**: Displays words read, total time, active time, and actual WPM at completion
- **URL support**: Fetch and read articles directly from URLs with automatic content extraction
- **Uniform text sizing**: Font size is based on the longest word for consistent display

## Examples

```bash
# Read at 300 words per minute
./speedread -wpm 300 book.txt

# Add a 500ms pause after punctuation for better comprehension
./speedread -wpm 250 -p 500 article.txt

# Disable focal point highlighting
./speedread -focal=false document.txt

# Use blue focal point color with context display
./speedread -c blue -context chapter.txt

# Read an online article
./speedread -wpm 350 https://example.com/blog/post

# Combine options
cat story.txt | ./speedread -wpm 350 -punct-pause 200
```

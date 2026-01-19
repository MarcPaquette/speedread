# speedread

A terminal-based speed reading application that displays text one word at a time using large ASCII art.

## Installation

```bash
go build -o speedread .
```

## Usage

```bash
# Read from a file
./speedread filename.txt

# Read from stdin
cat filename.txt | ./speedread
echo "Hello, world!" | ./speedread
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-wpm` | Words per minute (10-1000) | 200 |
| `-punct-pause`, `-p` | Extra pause after punctuation in milliseconds | 0 |

## Controls

| Key | Action |
|-----|--------|
| `Space` | Pause/unpause |
| `Ctrl+C` | Exit |

## Examples

```bash
# Read at 300 words per minute
./speedread -wpm 300 book.txt

# Add a 500ms pause after punctuation for better comprehension
./speedread -wpm 250 -p 500 article.txt

# Combine options
cat story.txt | ./speedread -wpm 350 -punct-pause 200
```

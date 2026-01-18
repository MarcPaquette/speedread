package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// ASCII art font - each letter is 5 rows tall
var font = map[rune][]string{
	'A': {
		"  █████  ",
		" ██   ██ ",
		" ███████ ",
		" ██   ██ ",
		" ██   ██ ",
	},
	'B': {
		" ██████  ",
		" ██   ██ ",
		" ██████  ",
		" ██   ██ ",
		" ██████  ",
	},
	'C': {
		"  █████  ",
		" ██      ",
		" ██      ",
		" ██      ",
		"  █████  ",
	},
	'D': {
		" ██████  ",
		" ██   ██ ",
		" ██   ██ ",
		" ██   ██ ",
		" ██████  ",
	},
	'E': {
		" ███████ ",
		" ██      ",
		" █████   ",
		" ██      ",
		" ███████ ",
	},
	'F': {
		" ███████ ",
		" ██      ",
		" █████   ",
		" ██      ",
		" ██      ",
	},
	'G': {
		"  █████  ",
		" ██      ",
		" ██  ███ ",
		" ██   ██ ",
		"  █████  ",
	},
	'H': {
		" ██   ██ ",
		" ██   ██ ",
		" ███████ ",
		" ██   ██ ",
		" ██   ██ ",
	},
	'I': {
		" ███████ ",
		"   ██    ",
		"   ██    ",
		"   ██    ",
		" ███████ ",
	},
	'J': {
		" ███████ ",
		"     ██  ",
		"     ██  ",
		" ██  ██  ",
		"  ████   ",
	},
	'K': {
		" ██   ██ ",
		" ██  ██  ",
		" █████   ",
		" ██  ██  ",
		" ██   ██ ",
	},
	'L': {
		" ██      ",
		" ██      ",
		" ██      ",
		" ██      ",
		" ███████ ",
	},
	'M': {
		" ██   ██ ",
		" ███ ███ ",
		" ██ █ ██ ",
		" ██   ██ ",
		" ██   ██ ",
	},
	'N': {
		" ██   ██ ",
		" ███  ██ ",
		" ██ █ ██ ",
		" ██  ███ ",
		" ██   ██ ",
	},
	'O': {
		"  █████  ",
		" ██   ██ ",
		" ██   ██ ",
		" ██   ██ ",
		"  █████  ",
	},
	'P': {
		" ██████  ",
		" ██   ██ ",
		" ██████  ",
		" ██      ",
		" ██      ",
	},
	'Q': {
		"  █████  ",
		" ██   ██ ",
		" ██   ██ ",
		" ██  ██  ",
		"  ████ █ ",
	},
	'R': {
		" ██████  ",
		" ██   ██ ",
		" ██████  ",
		" ██  ██  ",
		" ██   ██ ",
	},
	'S': {
		"  █████  ",
		" ██      ",
		"  █████  ",
		"      ██ ",
		"  █████  ",
	},
	'T': {
		" ███████ ",
		"   ██    ",
		"   ██    ",
		"   ██    ",
		"   ██    ",
	},
	'U': {
		" ██   ██ ",
		" ██   ██ ",
		" ██   ██ ",
		" ██   ██ ",
		"  █████  ",
	},
	'V': {
		" ██   ██ ",
		" ██   ██ ",
		" ██   ██ ",
		"  ██ ██  ",
		"   ███   ",
	},
	'W': {
		" ██   ██ ",
		" ██   ██ ",
		" ██ █ ██ ",
		" ███ ███ ",
		" ██   ██ ",
	},
	'X': {
		" ██   ██ ",
		"  ██ ██  ",
		"   ███   ",
		"  ██ ██  ",
		" ██   ██ ",
	},
	'Y': {
		" ██   ██ ",
		"  ██ ██  ",
		"   ███   ",
		"   ██    ",
		"   ██    ",
	},
	'Z': {
		" ███████ ",
		"     ██  ",
		"   ██    ",
		"  ██     ",
		" ███████ ",
	},
	'0': {
		"  █████  ",
		" ██  ███ ",
		" ██ █ ██ ",
		" ███  ██ ",
		"  █████  ",
	},
	'1': {
		"   ██    ",
		"  ███    ",
		"   ██    ",
		"   ██    ",
		" ███████ ",
	},
	'2': {
		"  █████  ",
		" ██   ██ ",
		"    ██   ",
		"  ██     ",
		" ███████ ",
	},
	'3': {
		"  █████  ",
		"      ██ ",
		"   ████  ",
		"      ██ ",
		"  █████  ",
	},
	'4': {
		" ██   ██ ",
		" ██   ██ ",
		" ███████ ",
		"      ██ ",
		"      ██ ",
	},
	'5': {
		" ███████ ",
		" ██      ",
		" ██████  ",
		"      ██ ",
		" ██████  ",
	},
	'6': {
		"  █████  ",
		" ██      ",
		" ██████  ",
		" ██   ██ ",
		"  █████  ",
	},
	'7': {
		" ███████ ",
		"     ██  ",
		"    ██   ",
		"   ██    ",
		"   ██    ",
	},
	'8': {
		"  █████  ",
		" ██   ██ ",
		"  █████  ",
		" ██   ██ ",
		"  █████  ",
	},
	'9': {
		"  █████  ",
		" ██   ██ ",
		"  ██████ ",
		"      ██ ",
		"  █████  ",
	},
	'.': {
		"         ",
		"         ",
		"         ",
		"         ",
		"   ██    ",
	},
	',': {
		"         ",
		"         ",
		"         ",
		"   ██    ",
		"  ██     ",
	},
	'!': {
		"   ██    ",
		"   ██    ",
		"   ██    ",
		"         ",
		"   ██    ",
	},
	'?': {
		"  █████  ",
		" ██   ██ ",
		"    ██   ",
		"         ",
		"    ██   ",
	},
	'\'': {
		"   ██    ",
		"  ██     ",
		"         ",
		"         ",
		"         ",
	},
	'"': {
		" ██  ██  ",
		" ██  ██  ",
		"         ",
		"         ",
		"         ",
	},
	'-': {
		"         ",
		"         ",
		" ███████ ",
		"         ",
		"         ",
	},
	' ': {
		"         ",
		"         ",
		"         ",
		"         ",
		"         ",
	},
}

const fontHeight = 5
const charWidth = 9            // All characters are exactly 9 columns wide
const targetHeightPercent = 0.5 // Use 50% of terminal height for text

func renderWord(word string, termWidth, termHeight int) []string {
	word = strings.ToUpper(word)

	// Calculate total width of the word (all chars are same width)
	totalWidth := len([]rune(word)) * charWidth

	// Calculate scale factor to fit terminal
	maxTextWidth := termWidth - 4 // Leave some margin
	maxTextHeight := termHeight - 4

	// Calculate base scale for uniform height across most words
	// Use the smaller of: height-based scale OR width-based scale for typical word
	heightScale := float64(maxTextHeight) * targetHeightPercent / float64(fontHeight)
	referenceWordLength := 8 // Typical word length for uniform scaling
	referenceWidth := referenceWordLength * charWidth
	widthScale := float64(maxTextWidth) / float64(referenceWidth)

	baseScale := heightScale
	if widthScale < heightScale {
		baseScale = widthScale
	}

	// Use baseScale for all words, but scale down further if word is too wide
	scale := baseScale
	scaledWidth := float64(totalWidth) * baseScale
	if scaledWidth > float64(maxTextWidth) {
		// Word too wide even at baseScale - scale down to fit
		scale = float64(maxTextWidth) / float64(totalWidth)
	}

	// Cap maximum scale at 3.0
	if scale > 3.0 {
		scale = 3.0
	}
	// Cap minimum scale at 1.0 (no scaling below base font)
	if scale < 1.0 {
		scale = 1.0
	}

	// Build the output lines
	lines := make([]string, fontHeight)
	for row := 0; row < fontHeight; row++ {
		var line strings.Builder
		for _, ch := range word {
			glyph, ok := font[ch]
			if !ok {
				glyph = font[' ']
			}
			if row < len(glyph) {
				// Ensure each glyph row is exactly charWidth characters
				glyphRow := glyph[row]
				runeCount := len([]rune(glyphRow))
				if runeCount < charWidth {
					glyphRow = glyphRow + strings.Repeat(" ", charWidth-runeCount)
				} else if runeCount > charWidth {
					glyphRow = string([]rune(glyphRow)[:charWidth])
				}
				line.WriteString(glyphRow)
			} else {
				line.WriteString(strings.Repeat(" ", charWidth))
			}
		}
		lines[row] = line.String()
	}

	// Scale up if scale > 1 (use rounding for better accuracy)
	if scale > 1.0 {
		scaleFactor := int(scale + 0.5)
		if scaleFactor < 1 {
			scaleFactor = 1
		}
		lines = scaleUp(lines, scaleFactor)
	}

	// Center horizontally
	for i, line := range lines {
		lineWidth := len([]rune(line))
		if lineWidth < termWidth {
			padding := (termWidth - lineWidth) / 2
			lines[i] = strings.Repeat(" ", padding) + line
		}
	}

	// Add vertical centering padding
	outputHeight := len(lines)
	if outputHeight < termHeight {
		topPadding := (termHeight - outputHeight) / 2
		paddedLines := make([]string, 0, termHeight)
		for i := 0; i < topPadding; i++ {
			paddedLines = append(paddedLines, "")
		}
		paddedLines = append(paddedLines, lines...)
		lines = paddedLines
	}

	return lines
}

func scaleUp(lines []string, factor int) []string {
	if factor <= 1 {
		return lines
	}

	var result []string
	for _, line := range lines {
		// Scale horizontally
		var scaledLine strings.Builder
		for _, ch := range line {
			for i := 0; i < factor; i++ {
				scaledLine.WriteRune(ch)
			}
		}
		scaled := scaledLine.String()

		// Scale vertically
		for i := 0; i < factor; i++ {
			result = append(result, scaled)
		}
	}
	return result
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func readInput(filename string) (string, error) {
	var reader io.Reader

	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		reader = file
	} else {
		// Check if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return "", fmt.Errorf("no input: provide a filename or pipe text to stdin")
		}
		reader = os.Stdin
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return string(content), nil
}

func tokenizeWords(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}

func getTerminalSize() (width, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default fallback
		return 80, 24
	}
	return width, height
}

func endsWithPunctuation(word string) bool {
	if len(word) == 0 {
		return false
	}
	lastChar := rune(word[len(word)-1])
	switch lastChar {
	case '.', ',', '!', '?', ';', ':', '"', '\'':
		return true
	}
	return false
}

func main() {
	wpm := flag.Int("wpm", 200, "Words per minute (10-1000)")
	punctPause := flag.Int("punct-pause", 0, "Extra pause after punctuation in milliseconds")
	flag.IntVar(punctPause, "p", 0, "Extra pause after punctuation in milliseconds (shorthand)")
	flag.Parse()

	// Validate WPM
	if *wpm < 10 {
		*wpm = 10
	}
	if *wpm > 1000 {
		*wpm = 1000
	}

	// Get filename from remaining args
	var filename string
	args := flag.Args()
	if len(args) > 0 {
		filename = args[0]
	}

	// Read input
	text, err := readInput(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Tokenize
	words := tokenizeWords(text)
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no words found in input")
		os.Exit(1)
	}

	// Calculate delay between words
	delay := time.Duration(float64(time.Minute) / float64(*wpm))

	// Display each word
	for i, word := range words {
		termWidth, termHeight := getTerminalSize()
		clearScreen()

		// Render and display the word
		lines := renderWord(word, termWidth, termHeight)
		for _, line := range lines {
			fmt.Println(line)
		}

		// Show progress at bottom
		progress := fmt.Sprintf("\n[%d/%d] %d WPM - Press Ctrl+C to exit", i+1, len(words), *wpm)
		fmt.Print(progress)

		time.Sleep(delay)

		// Add extra pause after punctuation
		if *punctPause > 0 && endsWithPunctuation(word) {
			time.Sleep(time.Duration(*punctPause) * time.Millisecond)
		}
	}

	// Final clear and message
	clearScreen()
	fmt.Println("Done! Read", len(words), "words.")
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	readability "github.com/go-shiori/go-readability"
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

func renderWord(word string, termWidth, termHeight int, focal bool) []string {
	word = strings.ToUpper(word)
	wordRunes := []rune(word)
	wordLen := len(wordRunes)

	// Calculate ORP index for focal point highlighting
	orpIndex := calculateORP(wordLen)

	// Calculate total width of the word (all chars are same width)
	totalWidth := wordLen * charWidth

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

	// Calculate integer scale factor
	scaleFactor := int(scale + 0.5)
	if scaleFactor < 1 {
		scaleFactor = 1
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
		lines = scaleUp(lines, scaleFactor)
	}

	// Calculate scaled character width
	scaledCharWidth := charWidth * scaleFactor

	// Calculate horizontal positioning
	lineWidth := len([]rune(lines[0]))
	var padding int
	if focal {
		// ORP-based centering: position ORP character center at screen center
		orpCenterCol := orpIndex*scaledCharWidth + (scaledCharWidth / 2)
		padding = (termWidth / 2) - orpCenterCol
	} else {
		// Traditional centering
		padding = (termWidth - lineWidth) / 2
	}

	// Ensure padding is non-negative
	if padding < 0 {
		padding = 0
	}

	// Apply horizontal padding and ORP coloring
	for i, line := range lines {
		if padding > 0 {
			lines[i] = strings.Repeat(" ", padding) + line
		}
		if focal && wordLen > 0 {
			// Calculate ORP column range (after padding)
			orpStartCol := padding + orpIndex*scaledCharWidth
			orpEndCol := orpStartCol + scaledCharWidth
			lines[i] = colorizeORPColumn(lines[i], orpStartCol, orpEndCol)
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

// calculateORP returns the Optimal Recognition Point index for a word.
// Based on Spritz algorithm: position depends on word length.
func calculateORP(wordLen int) int {
	switch {
	case wordLen <= 1:
		return 0
	case wordLen <= 5:
		return 1
	case wordLen <= 9:
		return 2
	case wordLen <= 13:
		return 3
	default:
		return 4
	}
}

// colorizeORPColumn applies ANSI red color to a specific column range in a line.
// startCol and endCol are 0-indexed rune positions.
func colorizeORPColumn(line string, startCol, endCol int) string {
	runes := []rune(line)
	if startCol < 0 || startCol >= len(runes) {
		return line
	}
	if endCol > len(runes) {
		endCol = len(runes)
	}

	var result strings.Builder
	result.WriteString(string(runes[:startCol]))
	result.WriteString("\033[31m") // Red
	result.WriteString(string(runes[startCol:endCol]))
	result.WriteString("\033[0m") // Reset
	if endCol < len(runes) {
		result.WriteString(string(runes[endCol:]))
	}
	return result.String()
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func isURL(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	article, err := readability.FromReader(resp.Body, nil)
	if err != nil {
		return "", fmt.Errorf("failed to extract content: %w", err)
	}

	return article.TextContent, nil
}

func readInput(input string) (string, error) {
	// Check if input is a URL
	if isURL(input) {
		return fetchURL(input)
	}

	var reader io.Reader

	if input != "" {
		file, err := os.Open(input)
		if err != nil {
			return "", fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		reader = file
	} else {
		// Check if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return "", fmt.Errorf("no input: provide a filename, URL, or pipe text to stdin")
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

func renderProgressBar(width, current, total int) string {
	// Reserve space for brackets and percentage: [████░░░░] 100%
	percentStr := fmt.Sprintf(" %d%%", current*100/total)
	barWidth := width - 2 - len(percentStr) // 2 for brackets
	if barWidth < 10 {
		barWidth = 10
	}

	filled := barWidth * current / total
	empty := barWidth - filled

	var bar strings.Builder
	bar.WriteString("[")
	bar.WriteString(strings.Repeat("█", filled))
	bar.WriteString(strings.Repeat("░", empty))
	bar.WriteString("]")
	bar.WriteString(percentStr)

	return bar.String()
}

func main() {
	wpm := flag.Int("wpm", 200, "Words per minute (10-1000)")
	punctPause := flag.Int("punct-pause", 0, "Extra pause after punctuation in milliseconds")
	flag.IntVar(punctPause, "p", 0, "Extra pause after punctuation in milliseconds (shorthand)")
	focal := flag.Bool("focal", true, "Enable focal point highlighting (Spritz-style)")
	flag.Parse()

	// Validate WPM
	if *wpm < 10 {
		*wpm = 10
	}
	if *wpm > 1000 {
		*wpm = 1000
	}

	// Atomic WPM for thread-safe adjustment during reading
	var currentWPM atomic.Int32
	currentWPM.Store(int32(*wpm))

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

	// Open /dev/tty for keyboard input (works even when stdin is piped)
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening tty: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	// Set up terminal raw mode for keyboard input
	oldState, err := term.MakeRaw(int(tty.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting raw mode: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(tty.Fd()), oldState)

	// Pause state
	var paused atomic.Bool

	// Current word index for navigation
	var currentIndex atomic.Int32
	totalWords := int32(len(words))

	// Goroutine to handle keyboard input
	go func() {
		buf := make([]byte, 3)
		for {
			n, _ := tty.Read(buf)
			if n == 0 {
				continue
			}

			// Check for escape sequence (arrow keys)
			if n >= 3 && buf[0] == 27 && buf[1] == '[' {
				switch buf[2] {
				case 'A': // Up arrow - increase WPM
					newWPM := currentWPM.Load() + 25
					if newWPM > 1000 {
						newWPM = 1000
					}
					currentWPM.Store(newWPM)
				case 'B': // Down arrow - decrease WPM
					newWPM := currentWPM.Load() - 25
					if newWPM < 10 {
						newWPM = 10
					}
					currentWPM.Store(newWPM)
				case 'D': // Left arrow - rewind
					newIdx := currentIndex.Load() - 1
					if newIdx < 0 {
						newIdx = 0
					}
					currentIndex.Store(newIdx)
				case 'C': // Right arrow - skip forward
					newIdx := currentIndex.Load() + 1
					if newIdx >= totalWords {
						newIdx = totalWords - 1
					}
					currentIndex.Store(newIdx)
				}
				continue
			}

			// Single character commands
			if buf[0] == ' ' {
				paused.Store(!paused.Load())
			} else if buf[0] == 3 { // Ctrl+C
				term.Restore(int(tty.Fd()), oldState)
				clearScreen()
				fmt.Print("Interrupted.\r\n")
				os.Exit(0)
			}
		}
	}()

	// Display each word
	for currentIndex.Load() < totalWords {
		i := int(currentIndex.Load())
		word := words[i]

		// Wait while paused
		for paused.Load() {
			// Re-read index in case user navigated while paused
			i = int(currentIndex.Load())
			word = words[i]

			termWidth, termHeight := getTerminalSize()
			clearScreen()
			lines := renderWord(word, termWidth, termHeight, *focal)
			for _, line := range lines {
				fmt.Print(line + "\r\n")
			}
			progressBar := renderProgressBar(termWidth, i+1, len(words))
			fmt.Print("\r\n" + progressBar)
			progress := fmt.Sprintf("\r\n%d WPM - PAUSED (space, ↑↓ speed, ←→ nav)", currentWPM.Load())
			fmt.Print(progress)
			time.Sleep(100 * time.Millisecond)
		}

		// Re-read index in case user navigated
		i = int(currentIndex.Load())
		word = words[i]

		termWidth, termHeight := getTerminalSize()
		clearScreen()

		// Render and display the word
		lines := renderWord(word, termWidth, termHeight, *focal)
		for _, line := range lines {
			fmt.Print(line + "\r\n")
		}

		// Show progress at bottom
		wpmNow := currentWPM.Load()
		progressBar := renderProgressBar(termWidth, i+1, len(words))
		fmt.Print("\r\n" + progressBar)
		progress := fmt.Sprintf("\r\n%d WPM - Space, ↑↓ speed, ←→ nav, Ctrl+C exit", wpmNow)
		fmt.Print(progress)

		// Calculate delay based on current WPM
		delay := time.Duration(float64(time.Minute) / float64(wpmNow))
		time.Sleep(delay)

		// Add extra pause after punctuation
		if *punctPause > 0 && endsWithPunctuation(word) {
			time.Sleep(time.Duration(*punctPause) * time.Millisecond)
		}

		// Advance to next word (if not navigated away)
		currentIndex.CompareAndSwap(int32(i), int32(i+1))
	}

	// Final clear and message
	clearScreen()
	fmt.Print("Done! Read ", len(words), " words.\r\n")
}

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func renderWord(word string, termWidth, termHeight int, focal bool, focalColorCode string, maxWordLen int) []string {
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

	// Calculate base scale for uniform height across all words
	// Use the smaller of: height-based scale OR width-based scale for longest word
	heightScale := float64(maxTextHeight) * targetHeightPercent / float64(fontHeight)
	referenceWidth := maxWordLen * charWidth
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

	// Apply horizontal padding, truncate to fit, and ORP coloring
	for i, line := range lines {
		if padding > 0 {
			lines[i] = strings.Repeat(" ", padding) + line
		}

		// Truncate to terminal width to prevent wrapping
		lineRunes := []rune(lines[i])
		if len(lineRunes) > termWidth {
			lines[i] = string(lineRunes[:termWidth])
		}

		if focal && wordLen > 0 {
			// Calculate ORP column range (after padding)
			orpStartCol := padding + orpIndex*scaledCharWidth
			orpEndCol := orpStartCol + scaledCharWidth
			lines[i] = colorizeORPColumn(lines[i], orpStartCol, orpEndCol, focalColorCode)
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

// colorToANSI converts a color name to its ANSI escape code
func colorToANSI(color string) string {
	colors := map[string]string{
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
	}
	if code, ok := colors[strings.ToLower(color)]; ok {
		return code
	}
	return "\033[31m" // Default to red
}

// colorizeORPColumn applies ANSI color to a specific column range in a line.
// startCol and endCol are 0-indexed rune positions.
func colorizeORPColumn(line string, startCol, endCol int, colorCode string) string {
	runes := []rune(line)
	if startCol < 0 || startCol >= len(runes) {
		return line
	}
	if endCol > len(runes) {
		endCol = len(runes)
	}

	var result strings.Builder
	result.WriteString(string(runes[:startCol]))
	result.WriteString(colorCode)
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

func findMaxWordLen(words []string) int {
	maxLen := 0
	for _, word := range words {
		wordLen := len([]rune(word))
		if wordLen > maxLen {
			maxLen = wordLen
		}
	}
	// Ensure minimum of 1 to avoid division by zero
	if maxLen < 1 {
		maxLen = 1
	}
	return maxLen
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

func endsWithSentence(word string) bool {
	if len(word) == 0 {
		return false
	}
	lastChar := rune(word[len(word)-1])
	switch lastChar {
	case '.', '!', '?':
		return true
	}
	return false
}

// Bookmark functions for saving/resuming reading position
func getBookmarkPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "speedread", "bookmarks.json")
}

func loadBookmarks() map[string]int {
	bookmarks := make(map[string]int)
	path := getBookmarkPath()
	if path == "" {
		return bookmarks
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return bookmarks
	}

	json.Unmarshal(data, &bookmarks)
	return bookmarks
}

func saveBookmark(filename string, position int) {
	path := getBookmarkPath()
	if path == "" || filename == "" {
		return
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	bookmarks := loadBookmarks()

	// Use absolute path as key
	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}

	if position <= 0 {
		delete(bookmarks, absPath) // Remove bookmark if at start
	} else {
		bookmarks[absPath] = position
	}

	data, err := json.Marshal(bookmarks)
	if err != nil {
		return
	}
	os.WriteFile(path, data, 0644)
}

func getBookmark(filename string) int {
	if filename == "" {
		return 0
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}

	bookmarks := loadBookmarks()
	return bookmarks[absPath]
}

func formatTimeRemaining(remainingWords, wpm int) string {
	if wpm <= 0 {
		return ""
	}
	// Calculate remaining seconds
	seconds := remainingWords * 60 / wpm
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	secs := seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
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
	focalColor := flag.String("focal-color", "red", "Focal point color (black, red, green, yellow, blue, magenta, cyan, white)")
	flag.StringVar(focalColor, "c", "red", "Focal point color (shorthand)")
	showContext := flag.Bool("context", false, "Show surrounding words (prev/next) for context")
	flag.Parse()

	// Validate WPM
	if *wpm < 10 {
		*wpm = 10
	}
	if *wpm > 1000 {
		*wpm = 1000
	}

	// Convert focal color to ANSI code
	focalColorCode := colorToANSI(*focalColor)

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

	// Find longest word for uniform font sizing
	maxWordLen := findMaxWordLen(words)

	// Check for saved bookmark (only for file input)
	startPosition := 0
	if filename != "" && !isURL(filename) {
		savedPos := getBookmark(filename)
		if savedPos > 0 && savedPos < len(words) {
			fmt.Printf("Found bookmark at word %d/%d (%.0f%%). Resume? [Y/n] ", savedPos+1, len(words), float64(savedPos)/float64(len(words))*100)
			var response string
			fmt.Scanln(&response)
			if response == "" || strings.ToLower(response) == "y" {
				startPosition = savedPos
			}
		}
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
	currentIndex.Store(int32(startPosition))
	totalWords := int32(len(words))

	// Session statistics tracking
	sessionStart := time.Now()
	var totalPauseTime time.Duration
	var pauseStart time.Time

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
			} else if buf[0] >= '0' && buf[0] <= '9' {
				// Number keys: jump to percentage (0=0%, 1=10%, ..., 9=90%)
				percent := int32(buf[0]-'0') * 10
				newIdx := totalWords * percent / 100
				if newIdx >= totalWords {
					newIdx = totalWords - 1
				}
				currentIndex.Store(newIdx)
			} else if buf[0] == 3 { // Ctrl+C
				// Save bookmark before exiting
				if filename != "" && !isURL(filename) {
					saveBookmark(filename, int(currentIndex.Load()))
				}
				term.Restore(int(tty.Fd()), oldState)
				clearScreen()
				fmt.Print("Interrupted. Position saved.\r\n")
				os.Exit(0)
			}
		}
	}()

	// Display each word
	for currentIndex.Load() < totalWords {
		i := int(currentIndex.Load())
		word := words[i]

		// Wait while paused
		if paused.Load() {
			pauseStart = time.Now()
		}
		for paused.Load() {
			// Re-read index in case user navigated while paused
			i = int(currentIndex.Load())
			word = words[i]

			termWidth, termHeight := getTerminalSize()
			clearScreen()

			// Show context: previous word (dimmed)
			if *showContext && i > 0 {
				prevWord := words[i-1]
				padding := (termWidth - len(prevWord)) / 2
				if padding < 0 {
					padding = 0
				}
				fmt.Printf("%s\033[2m%s\033[0m\r\n", strings.Repeat(" ", padding), prevWord)
			}

			lines := renderWord(word, termWidth, termHeight, *focal, focalColorCode, maxWordLen)
			for _, line := range lines {
				fmt.Print(line + "\r\n")
			}

			// Show context: next word (dimmed)
			if *showContext && i < len(words)-1 {
				nextWord := words[i+1]
				padding := (termWidth - len(nextWord)) / 2
				if padding < 0 {
					padding = 0
				}
				fmt.Printf("%s\033[2m%s\033[0m\r\n", strings.Repeat(" ", padding), nextWord)
			}

			wpmNow := int(currentWPM.Load())
			remaining := len(words) - i - 1
			timeLeft := formatTimeRemaining(remaining, wpmNow)
			progressBar := renderProgressBar(termWidth, i+1, len(words))
			fmt.Print("\r\n" + progressBar)
			progress := fmt.Sprintf("\r\n%d WPM | %s left - PAUSED (space, ↑↓, ←→, 0-9)", wpmNow, timeLeft)
			fmt.Print(progress)
			time.Sleep(100 * time.Millisecond)
		}
		if !pauseStart.IsZero() {
			totalPauseTime += time.Since(pauseStart)
			pauseStart = time.Time{}
		}

		// Re-read index in case user navigated
		i = int(currentIndex.Load())
		word = words[i]

		termWidth, termHeight := getTerminalSize()
		clearScreen()

		// Show context: previous word (dimmed)
		if *showContext && i > 0 {
			prevWord := words[i-1]
			padding := (termWidth - len(prevWord)) / 2
			if padding < 0 {
				padding = 0
			}
			fmt.Printf("%s\033[2m%s\033[0m\r\n", strings.Repeat(" ", padding), prevWord)
		}

		// Render and display the word
		lines := renderWord(word, termWidth, termHeight, *focal, focalColorCode, maxWordLen)
		for _, line := range lines {
			fmt.Print(line + "\r\n")
		}

		// Show context: next word (dimmed)
		if *showContext && i < len(words)-1 {
			nextWord := words[i+1]
			padding := (termWidth - len(nextWord)) / 2
			if padding < 0 {
				padding = 0
			}
			fmt.Printf("%s\033[2m%s\033[0m\r\n", strings.Repeat(" ", padding), nextWord)
		}

		// Show progress at bottom
		wpmNow := currentWPM.Load()
		remaining := len(words) - i - 1
		timeLeft := formatTimeRemaining(remaining, int(wpmNow))
		progressBar := renderProgressBar(termWidth, i+1, len(words))
		fmt.Print("\r\n" + progressBar)
		progress := fmt.Sprintf("\r\n%d WPM | %s left - Space, ↑↓, ←→, 0-9 jump", wpmNow, timeLeft)
		fmt.Print(progress)

		// Calculate delay based on current WPM with variable timing for word length
		baseDelay := float64(time.Minute) / float64(wpmNow)
		// Add 8% extra time per character above average length (5 chars)
		wordLen := len([]rune(word))
		if wordLen > 5 {
			extraChars := wordLen - 5
			baseDelay *= 1.0 + (float64(extraChars) * 0.08)
		}
		delay := time.Duration(baseDelay)
		time.Sleep(delay)

		// Add automatic pause at sentence boundaries (. ! ?)
		if endsWithSentence(word) {
			time.Sleep(150 * time.Millisecond)
		}

		// Add extra pause after other punctuation (user-configured)
		if *punctPause > 0 && endsWithPunctuation(word) && !endsWithSentence(word) {
			time.Sleep(time.Duration(*punctPause) * time.Millisecond)
		}

		// Advance to next word (if not navigated away)
		currentIndex.CompareAndSwap(int32(i), int32(i+1))
	}

	// Clear bookmark since reading is complete
	if filename != "" && !isURL(filename) {
		saveBookmark(filename, 0) // 0 removes the bookmark
	}

	// Final clear and session statistics
	clearScreen()
	sessionDuration := time.Since(sessionStart)
	activeTime := sessionDuration - totalPauseTime
	wordsRead := len(words)
	actualWPM := 0
	if activeTime.Minutes() > 0 {
		actualWPM = int(float64(wordsRead) / activeTime.Minutes())
	}

	fmt.Print("Session Complete!\r\n")
	fmt.Print("─────────────────\r\n")
	fmt.Printf("Words read:    %d\r\n", wordsRead)
	fmt.Printf("Total time:    %s\r\n", sessionDuration.Round(time.Second))
	fmt.Printf("Time paused:   %s\r\n", totalPauseTime.Round(time.Second))
	fmt.Printf("Active time:   %s\r\n", activeTime.Round(time.Second))
	fmt.Printf("Actual WPM:    %d\r\n", actualWPM)
}

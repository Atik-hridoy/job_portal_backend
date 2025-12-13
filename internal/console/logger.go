package console

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var useColor = detectColorSupport()

// Logger prints stylized console banners and status messages.
type Logger struct {
	prefix string
}

// New creates a console logger with a textual prefix to highlight the subsystem.
func New(prefix string) *Logger {
	return &Logger{prefix: strings.TrimSpace(prefix)}
}

// Banner prints a formatted multi-line banner with a headline and supporting copy.
func (l *Logger) Banner(title, message string) {
	headline := strings.TrimSpace(title)
	if l.prefix != "" {
		headline = fmt.Sprintf("%s · %s", strings.ToUpper(l.prefix), strings.ToUpper(headline))
	} else {
		headline = strings.ToUpper(headline)
	}
	border := strings.Repeat("═", len(headline)+6)

	log.Printf("\n%s\n%s %s %s\n%s\n%s%s\n\n",
		border,
		bannerEmoji(title),
		headline,
		bannerEmoji(title),
		border,
		indent(),
		message,
	)
}

// Status prints a single-line status with timestamp and label/detail pair.
func (l *Logger) Status(label, detail string) {
	timestamp := time.Now().Format("15:04:05")
	label = strings.ToUpper(label)
	log.Printf("[%s] %s %-12s%s %s", timestamp, cyan("▶"), label, reset(), detail)
}

// OTP logs the OTP dispatch in a consistent, highlighted format.
func (l *Logger) OTP(to, otp string) {
	log.Printf("%s OTP to %s → %s %s", cyan("🔐"), to, otp, reset())
}

func bannerEmoji(title string) string {
	slug := strings.ToLower(title)
	switch {
	case strings.Contains(slug, "error"):
		return "❌"
	case strings.Contains(slug, "warn"):
		return "⚠️"
	default:
		return "🚀"
	}
}

func indent() string {
	return strings.Repeat(" ", 4)
}

func cyan(text string) string {
	if !useColor {
		return text
	}
	return fmt.Sprintf("\033[36m%s\033[0m", text)
}

func reset() string {
	if !useColor {
		return ""
	}
	return "\033[0m"
}

func detectColorSupport() bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return false
	}
	if runtime.GOOS != "windows" {
		return true
	}
	// On Windows, enable color when running inside terminals known to support ANSI sequences.
	if _, ok := os.LookupEnv("WT_SESSION"); ok {
		return true
	}
	if _, ok := os.LookupEnv("ANSICON"); ok {
		return true
	}
	if term := os.Getenv("TERM"); strings.Contains(strings.ToLower(term), "xterm") {
		return true
	}
	return false
}

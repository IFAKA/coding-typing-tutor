package msgs

import (
	"time"

	"github.com/IFAKA/coding-typing-tutor/internal/snippets"
)

// Screen identifies which screen is currently active.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenTyping
	ScreenResults
	ScreenStats
)

// NavigateMsg tells the app router to switch to a different screen.
type NavigateMsg struct {
	To Screen
}

// StartTypingMsg carries the config and chosen snippet from menu → typing.
// For lesson mode, Code overrides Snippet.Code with generated text.
type StartTypingMsg struct {
	Snippet    snippets.Snippet
	Config     snippets.Config
	Code       string // non-empty in lesson mode; overrides Snippet.Code
	BestWPM    int    // personal best WPM for this language (for comparison)
	AvgWPM     int    // average WPM for this language
}

// RetryMsg tells the app to restart the current snippet.
type RetryMsg struct {
	Snippet snippets.Snippet
	Config  snippets.Config
	BestWPM int
	AvgWPM  int
}

// NextSnippetMsg tells the app router to pick and start the next snippet.
type NextSnippetMsg struct {
	Config  snippets.Config
	BestWPM int
	AvgWPM  int
}

// SnippetPlaceholder is a zero-value sentinel (unused externally).
var SnippetPlaceholder = snippets.Snippet{}

// KeyDelta tracks per-key attempt/error counts for one session.
type KeyDelta struct {
	Attempts int
	Errors   int
}

// TypingDoneMsg carries results from typing → results screen.
type TypingDoneMsg struct {
	Snippet    snippets.Snippet
	Config     snippets.Config
	WPM        int
	Accuracy   float64
	Errors     int
	Duration   time.Duration
	IsPersonalBest bool
	DiffFromAvg    int // WPM difference from personal average (can be negative)
	KeyDeltas      map[rune]KeyDelta
}

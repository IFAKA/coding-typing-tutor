package menu

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-typing-tutor/internal/history"
	"github.com/IFAKA/coding-typing-tutor/internal/keymap"
	"github.com/IFAKA/coding-typing-tutor/internal/lessons"
	"github.com/IFAKA/coding-typing-tutor/internal/snippets"
	"github.com/IFAKA/coding-typing-tutor/internal/sound"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/msgs"
)

// Model is the BubbleTea model for the menu screen.
type Model struct {
	langIdx    int
	diffIdx    int
	modeIdx    int
	lessonIdx  int // only active when mode == "lesson"
	entries    []history.Entry
	seenAt     map[string]time.Time
	progress   lessons.Progress
	width      int
	height     int
	activeRow  int // 0=lang, 1=diff, 2=mode, 3=lesson (lesson mode only)
}

// New creates a fresh menu model, loading history and preferences from disk.
func New(width, height int) Model {
	entries, _ := history.Load()
	prefs := history.LoadPrefs()
	return Model{
		langIdx:  prefs.LangIdx,
		diffIdx:  prefs.DiffIdx,
		modeIdx:  prefs.ModeIdx,
		entries:  entries,
		seenAt:   history.LastSeenMap(entries),
		progress: lessons.LoadProgress(),
		width:    width,
		height:   height,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenStats} }

		case "tab", "down", "j":
			m.activeRow = (m.activeRow + 1) % m.numRows()
			go sound.PlayNavRow()

		case "shift+tab", "up", "k":
			m.activeRow = (m.activeRow + m.numRows() - 1) % m.numRows()
			go sound.PlayNavRow()

		case "right", "l":
			m.cycleRight()
			go sound.PlayNavSelect()

		case "left", "h":
			m.cycleLeft()
			go sound.PlayNavSelect()

		case "enter", " ":
			return m, m.startTyping()
		}
	}
	return m, nil
}

func (m *Model) cycleRight() {
	switch m.activeRow {
	case 0:
		m.langIdx = (m.langIdx + 1) % len(snippets.Languages)
	case 1:
		m.diffIdx = (m.diffIdx + 1) % len(snippets.Difficulties)
	case 2:
		m.modeIdx = (m.modeIdx + 1) % len(snippets.Modes)
	case 3:
		m.lessonIdx = (m.lessonIdx + 1) % len(lessons.All)
	}
	m.savePrefs()
}

func (m *Model) cycleLeft() {
	switch m.activeRow {
	case 0:
		m.langIdx = (m.langIdx + len(snippets.Languages) - 1) % len(snippets.Languages)
	case 1:
		m.diffIdx = (m.diffIdx + len(snippets.Difficulties) - 1) % len(snippets.Difficulties)
	case 2:
		m.modeIdx = (m.modeIdx + len(snippets.Modes) - 1) % len(snippets.Modes)
	case 3:
		m.lessonIdx = (m.lessonIdx + len(lessons.All) - 1) % len(lessons.All)
	}
	m.savePrefs()
}

// numRows returns the number of navigable rows based on current mode.
func (m *Model) numRows() int {
	if m.isLessonMode() {
		return 4
	}
	return 3
}

// isLessonMode returns true when the current mode selection is "lesson".
func (m *Model) isLessonMode() bool {
	return snippets.Modes[m.modeIdx] == "lesson"
}

func (m *Model) savePrefs() {
	history.SavePrefs(history.Prefs{
		LangIdx: m.langIdx,
		DiffIdx: m.diffIdx,
		ModeIdx: m.modeIdx,
	})
}

func (m *Model) startTyping() tea.Cmd {
	lang := snippets.Languages[m.langIdx]
	diff := snippets.Difficulties[m.diffIdx]
	mode := snippets.Modes[m.modeIdx]

	if mode == "lesson" {
		lesson := lessons.All[m.lessonIdx]
		if !m.progress.Unlocked[lesson.Number] {
			// Lesson locked — do nothing (view shows lock indicator)
			return nil
		}
		code := lessons.Generate(lesson, 120)
		snippet := snippets.Snippet{
			ID:    fmt.Sprintf("lesson-%d", lesson.Number),
			Title: lesson.Name,
		}
		cfg := snippets.Config{
			Language:  "lesson",
			Mode:      "lesson",
			LessonNum: lesson.Number,
		}
		return func() tea.Msg {
			return msgs.StartTypingMsg{
				Snippet: snippet,
				Config:  cfg,
				Code:    code,
			}
		}
	}

	kstore, _ := keymap.Load()
	weak := keymap.WeakKeys(kstore, 0.15)
	snippet := snippets.Pick(lang, diff, m.seenAt, weak)
	if snippet == nil {
		return nil
	}

	stats := history.Compute(m.entries)
	langAvg := history.AvgWPMForLanguage(m.entries, lang)

	return func() tea.Msg {
		return msgs.StartTypingMsg{
			Snippet: *snippet,
			Config: snippets.Config{
				Language:   lang,
				Difficulty: diff,
				Mode:       mode,
			},
			BestWPM: stats.BestWPM,
			AvgWPM:  langAvg,
		}
	}
}

// ActiveLang returns the currently selected language.
func (m Model) ActiveLang() string { return snippets.Languages[m.langIdx] }

// ActiveDiff returns the currently selected difficulty.
func (m Model) ActiveDiff() string { return snippets.Difficulties[m.diffIdx] }

// ActiveMode returns the currently selected mode.
func (m Model) ActiveMode() string { return snippets.Modes[m.modeIdx] }

// ActiveRow returns the currently focused row (0=lang, 1=diff, 2=mode).
func (m Model) ActiveRow() int { return m.activeRow }

// Width returns current terminal width.
func (m Model) Width() int { return m.width }

// Height returns current terminal height.
func (m Model) Height() int { return m.height }

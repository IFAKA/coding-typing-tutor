package menu

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-type/internal/history"
	"github.com/IFAKA/coding-type/internal/snippets"
	"github.com/IFAKA/coding-type/internal/ui/msgs"
)

// Model is the BubbleTea model for the menu screen.
type Model struct {
	langIdx   int
	diffIdx   int
	modeIdx   int
	entries   []history.Entry
	seenAt    map[string]time.Time
	width     int
	height    int
	activeRow int // 0=lang, 1=diff, 2=mode
}

// New creates a fresh menu model, loading history and preferences from disk.
func New(width, height int) Model {
	entries, _ := history.Load()
	prefs := history.LoadPrefs()
	return Model{
		langIdx: prefs.LangIdx,
		diffIdx: prefs.DiffIdx,
		modeIdx: prefs.ModeIdx,
		entries: entries,
		seenAt:  history.LastSeenMap(entries),
		width:   width,
		height:  height,
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
			m.activeRow = (m.activeRow + 1) % 3

		case "shift+tab", "up", "k":
			m.activeRow = (m.activeRow + 2) % 3

		case "right", "l":
			m.cycleRight()

		case "left", "h":
			m.cycleLeft()

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
	}
	m.savePrefs()
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

	snippet := snippets.Pick(lang, diff, m.seenAt)
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

package stats

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-typing-tutor/internal/history"
	"github.com/IFAKA/coding-typing-tutor/internal/keymap"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/msgs"
)

// Model is the BubbleTea model for the stats/history screen.
type Model struct {
	stats    history.Stats
	keyStore keymap.Store
	activeTab int // 0 = overview, 1 = heatmap
	width    int
	height   int
}

// New creates a stats model by loading history from disk.
func New(width, height int) Model {
	entries, _ := history.Load()
	kstore, _ := keymap.Load()
	return Model{
		stats:    history.Compute(entries),
		keyStore: kstore,
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
		case "m", "esc":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenMenu} }
		case "tab":
			m.activeTab = 1 - m.activeTab
		}
	}
	return m, nil
}

func (m Model) Stats() history.Stats { return m.stats }
func (m Model) Width() int           { return m.width }
func (m Model) Height() int          { return m.height }

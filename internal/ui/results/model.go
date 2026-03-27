package results

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/msgs"
)

type resultTickMsg struct{}

// Model is the BubbleTea model for the results screen.
type Model struct {
	done       msgs.TypingDoneMsg
	width      int
	height     int
	displayWPM int
	displayAcc float64
	frame      int
	animDone   bool
}

// New creates a results model from a TypingDoneMsg.
func New(done msgs.TypingDoneMsg, width, height int) Model {
	return Model{done: done, width: width, height: height}
}

func (m Model) Init() tea.Cmd { return resultTickCmd() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case resultTickMsg:
		m.frame++
		if !m.animDone {
			// Ease-out count-up: large steps early, tiny steps near the target
			wpmStep := (m.done.WPM - m.displayWPM) / 4
			if wpmStep < 1 {
				wpmStep = 1
			}
			m.displayWPM = min(m.displayWPM+wpmStep, m.done.WPM)

			accStep := (m.done.Accuracy - m.displayAcc) / 4
			if accStep < 0.1 {
				accStep = 0.1
			}
			m.displayAcc = min(m.displayAcc+accStep, m.done.Accuracy)

			if m.displayWPM == m.done.WPM && m.displayAcc >= m.done.Accuracy-0.05 {
				m.displayAcc = m.done.Accuracy
				m.animDone = true
			}
		}
		// Keep ticking for sparkle animation after completion
		if !m.animDone || m.done.IsPersonalBest {
			return m, resultTickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			return m, func() tea.Msg {
				return msgs.RetryMsg{
					Snippet: m.done.Snippet,
					Config:  m.done.Config,
					BestWPM: m.currentBest(),
					AvgWPM:  m.done.WPM,
				}
			}

		case "n":
			return m, func() tea.Msg {
				return msgs.NextSnippetMsg{
					Config:  m.done.Config,
					BestWPM: m.currentBest(),
					AvgWPM:  m.done.WPM,
				}
			}

		case "m", "esc":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenMenu} }
		}
	}
	return m, nil
}

// currentBest returns the current personal best WPM after this session.
func (m Model) currentBest() int {
	if m.done.IsPersonalBest {
		return m.done.WPM
	}
	return 0
}

func resultTickCmd() tea.Cmd {
	return tea.Tick(40*time.Millisecond, func(_ time.Time) tea.Msg {
		return resultTickMsg{}
	})
}

// Done returns the underlying TypingDoneMsg.
func (m Model) Done() msgs.TypingDoneMsg { return m.done }
func (m Model) Width() int               { return m.width }
func (m Model) Height() int              { return m.height }

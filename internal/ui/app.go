package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-typing-tutor/internal/history"
	"github.com/IFAKA/coding-typing-tutor/internal/keymap"
	"github.com/IFAKA/coding-typing-tutor/internal/lessons"
	"github.com/IFAKA/coding-typing-tutor/internal/snippets"
	"github.com/IFAKA/coding-typing-tutor/internal/sound"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/menu"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/msgs"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/results"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/stats"
	"github.com/IFAKA/coding-typing-tutor/internal/ui/typing"
)

// App is the root BubbleTea model that routes between screens.
type App struct {
	screen  msgs.Screen
	menu    menu.Model
	typing  typing.Model
	results results.Model
	stats   stats.Model
	width   int
	height  int
}

// New creates the root app model.
func New() App {
	go sound.Init()
	return App{
		screen: msgs.ScreenMenu,
		menu:   menu.New(80, 24),
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height

	case msgs.NavigateMsg:
		return a.navigate(msg.To)

	case msgs.StartTypingMsg:
		a.screen = msgs.ScreenTyping
		a.typing = typing.New(msg, a.width, a.height)
		return a, a.typing.Init()

	case msgs.RetryMsg:
		a.screen = msgs.ScreenTyping
		a.typing = typing.New(msgs.StartTypingMsg{
			Snippet: msg.Snippet,
			Config:  msg.Config,
			BestWPM: msg.BestWPM,
			AvgWPM:  msg.AvgWPM,
		}, a.width, a.height)
		return a, a.typing.Init()

	case msgs.NextSnippetMsg:
		return a.pickAndStart(msg.Config, msg.BestWPM, msg.AvgWPM)

	case msgs.TypingDoneMsg:
		_ = history.Save(history.Entry{
			Timestamp:    time.Now(),
			Language:     msg.Config.Language,
			SnippetID:    msg.Snippet.ID,
			SnippetTitle: msg.Snippet.Title,
			WPM:          msg.WPM,
			Accuracy:     msg.Accuracy,
			DurationMs:   msg.Duration.Milliseconds(),
			Errors:       msg.Errors,
		})
		if len(msg.KeyDeltas) > 0 {
			km := make(map[rune]keymap.KeyDelta, len(msg.KeyDeltas))
			for r, d := range msg.KeyDeltas {
				km[r] = keymap.KeyDelta{Attempts: d.Attempts, Errors: d.Errors}
			}
			_ = keymap.Merge(km)
		}
		if msg.Config.Mode == "lesson" {
			p := lessons.LoadProgress()
			lessons.UpdateProgress(&p, msg.Config.LessonNum, msg.Accuracy)
			_ = lessons.SaveProgress(p)
		}
		a.screen = msgs.ScreenResults
		a.results = results.New(msg, a.width, a.height)
		return a, a.results.Init()
	}

	// Delegate to active screen
	return a.delegateToScreen(msg)
}

func (a App) navigate(to msgs.Screen) (tea.Model, tea.Cmd) {
	a.screen = to
	switch to {
	case msgs.ScreenMenu:
		a.menu = menu.New(a.width, a.height)
		return a, a.menu.Init()
	case msgs.ScreenStats:
		a.stats = stats.New(a.width, a.height)
		return a, a.stats.Init()
	}
	return a, nil
}

func (a App) pickAndStart(cfg snippets.Config, bestWPM, avgWPM int) (tea.Model, tea.Cmd) {
	entries, _ := history.Load()
	seenAt := history.LastSeenMap(entries)
	kstore, _ := keymap.Load()
	weak := keymap.WeakKeys(kstore, 0.15)
	snippet := snippets.Pick(cfg.Language, cfg.Difficulty, seenAt, weak)
	if snippet == nil {
		a.screen = msgs.ScreenMenu
		a.menu = menu.New(a.width, a.height)
		return a, nil
	}
	a.screen = msgs.ScreenTyping
	a.typing = typing.New(msgs.StartTypingMsg{
		Snippet: *snippet,
		Config:  cfg,
		BestWPM: bestWPM,
		AvgWPM:  avgWPM,
	}, a.width, a.height)
	return a, a.typing.Init()
}

func (a App) delegateToScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch a.screen {
	case msgs.ScreenMenu:
		m, cmd := a.menu.Update(msg)
		a.menu = m.(menu.Model)
		return a, cmd
	case msgs.ScreenTyping:
		m, cmd := a.typing.Update(msg)
		a.typing = m.(typing.Model)
		return a, cmd
	case msgs.ScreenResults:
		m, cmd := a.results.Update(msg)
		a.results = m.(results.Model)
		return a, cmd
	case msgs.ScreenStats:
		m, cmd := a.stats.Update(msg)
		a.stats = m.(stats.Model)
		return a, cmd
	}
	return a, nil
}

func (a App) View() string {
	switch a.screen {
	case msgs.ScreenMenu:
		return a.menu.View()
	case msgs.ScreenTyping:
		return a.typing.View()
	case msgs.ScreenResults:
		return a.results.View()
	case msgs.ScreenStats:
		return a.stats.View()
	}
	return ""
}

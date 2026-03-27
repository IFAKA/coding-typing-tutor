package typing

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-typing-tutor/internal/engine"
	kb "github.com/IFAKA/coding-typing-tutor/internal/keyboard"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

func (m Model) View() string {
	header := renderHeader(m)
	code := renderCode(m)
	statsBar := renderStats(m)
	help := "  " + theme.HelpKey.Render("esc") + " " + theme.HelpDesc.Render("menu") +
		"   " + theme.HelpKey.Render("ctrl+r") + " " + theme.HelpDesc.Render("restart")

	inner := strings.Join([]string{"", code, "", statsBar}, "\n")

	box := theme.RenderBox(inner, m.width, 0, 2)

	var curChar rune
	if m.state.Cursor < len(m.state.Target) {
		curChar = m.state.Target[m.state.Cursor]
	}

	var wrongFlash rune
	if m.wrongKeyFlash > 0 {
		wrongFlash = m.wrongExpected
	}

	keyboard := lipgloss.NewStyle().
		Width(lipgloss.Width(box)).
		Align(lipgloss.Center).
		Render(renderKeyboard(curChar, m.weakKeys, wrongFlash))

	// Finger hint shown for 2s after a wrong keypress
	var hint string
	if m.wrongKeyFlash > 0 {
		base, _ := kb.ResolveKey(m.wrongExpected)
		f := kb.ActiveFinger(base)
		if f >= 0 {
			keyLabel := string(m.wrongExpected)
			if m.wrongExpected == '\n' {
				keyLabel = "enter"
			} else if m.wrongExpected == ' ' {
				keyLabel = "space"
			}
			hint = "\n  " + lipgloss.NewStyle().Foreground(kb.FingerColor[f]).Render(
				fmt.Sprintf("use %s finger for '%s'", kb.FingerNames[f], keyLabel),
			)
		}
	}

	content := strings.Join([]string{header, box, help, "", keyboard}, "\n") + hint

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func renderHeader(m Model) string {
	dot := theme.Muted.Render(" · ")
	if m.config.Mode == "lesson" {
		badge := theme.HeaderBadge.Render(" lesson ")
		num := theme.Muted.Render(fmt.Sprintf("level %d", m.config.LessonNum))
		title := theme.Muted.Render(m.snippet.Title)
		return "  " + badge + dot + num + dot + title + "\n"
	}
	lang := theme.HeaderBadge.Render(" " + m.config.Language + " ")
	diff := diffStyle(m.config.Difficulty).Render(m.config.Difficulty)
	mode := theme.Muted.Render(m.config.Mode)
	title := theme.Muted.Render(m.snippet.Title)
	return "  " + lang + dot + diff + dot + mode + dot + title + "\n"
}

func renderCode(m Model) string {
	s := m.state
	var sb strings.Builder
	for i, r := range s.Target {
		style := charStyle(m, i)
		if r == '\n' {
			if i == s.Cursor {
				sb.WriteString(style.Render("↵"))
			}
			sb.WriteRune('\n')
		} else {
			sb.WriteString(style.Render(string(r)))
		}
	}
	return sb.String()
}

func charStyle(m Model, i int) lipgloss.Style {
	s := m.state
	if i == s.Cursor {
		if !m.cursorVisible {
			// Blink off: render as dim untyped so cursor position is still faintly visible
			if i < len(s.SyntaxColors) {
				return lipgloss.NewStyle().Foreground(s.SyntaxColors[i])
			}
			return theme.UntypedChar
		}
		if m.errorFlash > 0 {
			return theme.CursorError
		}
		return theme.CursorChar
	}

	switch s.States[i] {
	case engine.Correct:
		return theme.CorrectChar
	case engine.Incorrect:
		return theme.IncorrectChar
	default:
		if i < len(s.SyntaxColors) {
			return lipgloss.NewStyle().Foreground(s.SyntaxColors[i]).Faint(true)
		}
		return theme.UntypedChar
	}
}

func renderStats(m Model) string {
	s := m.state

	wpm := fmt.Sprintf("WPM %s", theme.StatValue.Render(fmt.Sprintf("%d", s.WPM())))
	acc := fmt.Sprintf("ACC %s", theme.StatValue.Render(fmt.Sprintf("%.0f%%", s.Accuracy())))

	var timerStr string
	if m.config.Mode == "timed" && s.Started {
		remaining := timedDuration - time.Since(s.StartedAt)
		if remaining < 0 {
			remaining = 0
		}
		secs := int(remaining.Seconds())
		timerStyle := theme.StatValue
		if secs <= 10 {
			timerStyle = lipgloss.NewStyle().Foreground(theme.Red).Bold(true)
		}
		timerStr = timerStyle.Render(fmt.Sprintf("%02d:%02d", secs/60, secs%60))
	} else {
		elapsed := s.ElapsedSeconds()
		timerStr = theme.StatValue.Render(fmt.Sprintf("%02d:%02d", elapsed/60, elapsed%60))
	}

	progress := fmt.Sprintf("%d/%d", s.Cursor, len(s.Target))

	dot := theme.Muted.Render("  ·  ")
	return "  " + strings.Join([]string{wpm, acc, timerStr, theme.Muted.Render(progress)}, dot)
}

func diffStyle(diff string) lipgloss.Style {
	switch diff {
	case "easy":
		return theme.DiffEasy
	case "medium":
		return theme.DiffMedium
	case "hard":
		return theme.DiffHard
	default:
		return theme.Muted
	}
}

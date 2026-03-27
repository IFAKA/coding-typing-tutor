package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-typing-tutor/internal/history"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

func (m Model) View() string {
	s := m.stats

	if s.TotalSessions == 0 {
		return m.emptyView()
	}

	if m.activeTab == 1 {
		return m.heatmapView()
	}

	aggregates := renderAggregates(s)
	recentTable := renderRecent(s.Recent)

	sep := theme.Separator.Render(strings.Repeat("─", 48))

	help := "  " + theme.HelpKey.Render("tab") + " " + theme.HelpDesc.Render("heatmap") +
		"   " + theme.HelpKey.Render("m") + " " + theme.HelpDesc.Render("menu") +
		"   " + theme.HelpKey.Render("q") + " " + theme.HelpDesc.Render("quit")

	inner := strings.Join([]string{
		"",
		aggregates,
		"",
		"  " + sep,
		"",
		recentTable,
		"",
		help,
		"",
	}, "\n")

	box := theme.RenderBox(inner, m.width, 0, 0)

	header := theme.Title.Render("  your stats")
	content := strings.Join([]string{header, box}, "\n")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func renderAggregates(s history.Stats) string {
	left := renderStat("best wpm", fmt.Sprintf("%d", s.BestWPM)) + "\n" +
		renderStat("avg wpm", fmt.Sprintf("%d", s.AvgWPM)) + "\n" +
		renderStat("avg acc", fmt.Sprintf("%.1f%%", s.AvgAccuracy))

	streakVal := fmt.Sprintf("%d", s.Streak)
	if s.Streak == 1 {
		streakVal += " day"
	} else {
		streakVal += " days"
	}
	if s.Streak >= 3 {
		streakVal += " 🔥"
	}

	right := renderStat("sessions", fmt.Sprintf("%d", s.TotalSessions)) + "\n" +
		renderStat("streak", streakVal) + "\n" +
		"  "

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(27).Render(left),
		lipgloss.NewStyle().Width(27).Render(right),
	)
}

func renderStat(label, value string) string {
	return "  " + theme.StatLabel.Render(fmt.Sprintf("%-12s", label)) +
		"  " + theme.StatValue.Render(value)
}

func renderRecent(entries []history.Entry) string {
	if len(entries) == 0 {
		return "  " + theme.Muted.Render("no sessions yet")
	}

	header := "  " + theme.Muted.Render(
		fmt.Sprintf("%-6s  %-6s  %-22s  %4s  %5s",
			"date", "lang", "snippet", "wpm", "acc"))

	var rows []string
	rows = append(rows, header)
	for _, e := range entries {
		date := e.Timestamp.Format("Jan 02")
		lang := e.Language
		if len(lang) > 6 {
			lang = lang[:6]
		}
		title := e.SnippetTitle
		if len(title) > 22 {
			title = title[:19] + "..."
		}
		row := fmt.Sprintf("  %-6s  %-6s  %-22s  %4d  %4.0f%%",
			date, lang, title, e.WPM, e.Accuracy)
		rows = append(rows, theme.HelpDesc.Render(row))
	}
	return strings.Join(rows, "\n")
}

func (m Model) emptyView() string {
	msg := theme.Muted.Render("no history yet — complete a snippet to see stats here")
	help := "  " + theme.HelpKey.Render("m") + " " + theme.HelpDesc.Render("menu")

	inner := "\n\n  " + msg + "\n\n" + help + "\n"
	box := theme.RenderBox(inner, m.width, 0, 0)
	header := theme.Title.Render("  your stats")
	content := strings.Join([]string{header, box}, "\n")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// keep time import used
var _ = time.Now

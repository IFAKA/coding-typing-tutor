package results

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

var sparkles = []string{"✦", "✧", "·", "✦", "·", "✧", "·", "✦"}

func (m Model) View() string {
	d := m.done

	wpmLabel := theme.StatLabel.Render("wpm")
	wpmValue := lipgloss.NewStyle().
		Foreground(theme.Yellow).
		Bold(true).
		Render(fmt.Sprintf("%d", m.displayWPM))

	// Sparkle decoration once animation completes on a personal best
	var wpmLine string
	if m.animDone && d.IsPersonalBest {
		sp := lipgloss.NewStyle().Foreground(theme.Yellow).Bold(true).
			Render(sparkles[m.frame%len(sparkles)])
		wpmLine = sp + " " + wpmValue + " " + sp + "  " + wpmLabel
	} else {
		wpmLine = wpmValue + "  " + wpmLabel
	}

	accLabel := theme.StatLabel.Render("accuracy")
	accValue := lipgloss.NewStyle().
		Foreground(theme.Green).
		Bold(true).
		Render(fmt.Sprintf("%.1f%%", m.displayAcc))

	wpmBlock := lipgloss.JoinVertical(lipgloss.Left,
		wpmLine,
		wpmSub(d.DiffFromAvg, d.IsPersonalBest),
	)

	accBlock := lipgloss.JoinVertical(lipgloss.Left,
		accValue+"  "+accLabel,
		accSub(d.Accuracy),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(22).Render(wpmBlock),
		lipgloss.NewStyle().Width(22).Render(accBlock),
	)

	elapsed := int(d.Duration.Seconds())
	details := theme.Muted.Render(
		fmt.Sprintf("time %02d:%02d   errors %d   chars %d",
			elapsed/60, elapsed%60, d.Errors, len([]rune(d.Snippet.Code))),
	)

	sep := theme.Separator.Render(strings.Repeat("─", 44))

	help := "  " + strings.Join([]string{
		theme.HelpKey.Render("r") + " " + theme.HelpDesc.Render("retry"),
		theme.HelpKey.Render("n") + " " + theme.HelpDesc.Render("next"),
		theme.HelpKey.Render("m") + " " + theme.HelpDesc.Render("menu"),
		theme.HelpKey.Render("q") + " " + theme.HelpDesc.Render("quit"),
	}, theme.Muted.Render("   "))

	snippetInfo := theme.Muted.Render(d.Snippet.Title) +
		theme.Muted.Render("  ·  ") +
		theme.Muted.Render(d.Config.Language)

	inner := strings.Join([]string{
		"",
		"  " + snippetInfo,
		"",
		"  " + topRow,
		"",
		"  " + details,
		"",
		"  " + sep,
		"",
		help,
		"",
	}, "\n")

	box := theme.RenderBox(inner, m.width, 0, 0)

	header := theme.Title.Render("  results")

	content := strings.Join([]string{header, box}, "\n")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func wpmSub(diff int, isPersonalBest bool) string {
	if isPersonalBest {
		return theme.PersonalBest.Render("  personal best ✓")
	}
	if diff > 0 {
		return theme.Success.Render(fmt.Sprintf("  +%d from avg", diff))
	}
	if diff < 0 {
		return theme.HelpDesc.Render(fmt.Sprintf("  %d from avg", diff))
	}
	return theme.Muted.Render("  avg")
}

func accSub(acc float64) string {
	switch {
	case acc >= 99:
		return theme.Success.Render("  perfect")
	case acc >= 95:
		return theme.Success.Render("  great")
	case acc >= 85:
		return theme.HelpDesc.Render("  good")
	default:
		return lipgloss.NewStyle().Foreground(theme.Red).Render("  needs work")
	}
}

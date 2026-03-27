package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	kb "github.com/IFAKA/coding-typing-tutor/internal/keyboard"
	"github.com/IFAKA/coding-typing-tutor/internal/keymap"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

func (m Model) heatmapView() string {
	help := "  " + theme.HelpKey.Render("tab") + " " + theme.HelpDesc.Render("overview") +
		"   " + theme.HelpKey.Render("m") + " " + theme.HelpDesc.Render("menu") +
		"   " + theme.HelpKey.Render("q") + " " + theme.HelpDesc.Render("quit")

	legend := renderLegend()
	keyboard := renderHeatmapKeyboard(m.keyStore)
	weakList := renderWeakList(m.keyStore)

	inner := strings.Join([]string{
		"",
		legend,
		"",
		keyboard,
		"",
		weakList,
		"",
		help,
		"",
	}, "\n")

	box := theme.RenderBox(inner, m.width, 0, 0)
	header := theme.Title.Render("  key heatmap")
	content := strings.Join([]string{header, box}, "\n")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func renderLegend() string {
	never  := lipgloss.NewStyle().Foreground(theme.Surface1).Render("■ never typed")
	good   := lipgloss.NewStyle().Foreground(theme.Green).Render("■ ≤5% errors")
	medium := lipgloss.NewStyle().Foreground(theme.Yellow).Render("■ 5–15%")
	weak   := lipgloss.NewStyle().Foreground(theme.Red).Render("■ >15%")
	return "  " + strings.Join([]string{never, good, medium, weak}, "  ")
}

func heatColor(ks keymap.KeyStats) lipgloss.Color {
	if ks.Attempts == 0 {
		return theme.Surface1
	}
	rate := keymap.ErrorRate(ks)
	switch {
	case rate <= 0.05:
		return theme.Green
	case rate <= 0.15:
		return theme.Yellow
	default:
		return theme.Red
	}
}

func renderHeatmapKeyboard(store keymap.Store) string {
	renderKey := func(k kb.KeyDef) string {
		label := k.Label()
		color := heatColor(store[k.Ch])
		return lipgloss.NewStyle().Foreground(color).Render(label)
	}

	shiftStyle := lipgloss.NewStyle().Foreground(theme.Surface1)

	var lines []string

	mutedStyle := lipgloss.NewStyle().Foreground(theme.Surface1)
	lines = append(lines, "  "+mutedStyle.Render("LEFT")+"               "+mutedStyle.Render("RIGHT"))

	for rowIdx, row := range kb.KbRows {
		var sb strings.Builder

		if rowIdx == 3 {
			sb.WriteString("   ")
			sb.WriteString(shiftStyle.Render("⇧"))
			sb.WriteString(" ")
		} else {
			sb.WriteString(kb.RowIndent[rowIdx])
		}

		for i, k := range row {
			if i == kb.GapAfter[rowIdx]+1 {
				sb.WriteString(" ")
			}
			sb.WriteString(renderKey(k))
			sb.WriteString(" ")
		}

		if rowIdx == 3 {
			sb.WriteString(shiftStyle.Render("⇧"))
		}

		lines = append(lines, sb.String())
	}

	// Space bar
	spaceColor := heatColor(store[' '])
	var spaceSB strings.Builder
	spaceSB.WriteString("        ")
	spaceSB.WriteString(lipgloss.NewStyle().Foreground(spaceColor).Render("___________"))
	lines = append(lines, spaceSB.String())

	return strings.Join(lines, "\n")
}

func renderWeakList(store keymap.Store) string {
	type keyEntry struct {
		r    rune
		ks   keymap.KeyStats
		rate float64
	}

	var entries []keyEntry
	for r, ks := range store {
		if ks.Attempts >= 5 {
			entries = append(entries, keyEntry{r, ks, keymap.ErrorRate(ks)})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].rate > entries[j].rate
	})

	if len(entries) == 0 {
		return "  " + theme.Muted.Render("type more to see weak key analysis")
	}

	header := "  " + theme.Muted.Render("weakest keys:")
	var parts []string
	limit := 8
	if len(entries) < limit {
		limit = len(entries)
	}
	for _, e := range entries[:limit] {
		label := string(e.r)
		if e.r == '\n' {
			label = "↵"
		} else if e.r == ' ' {
			label = "space"
		}
		color := heatColor(e.ks)
		part := lipgloss.NewStyle().Foreground(color).Render(
			fmt.Sprintf("%s %.0f%%", label, e.rate*100),
		)
		parts = append(parts, part)
	}
	return header + "  " + strings.Join(parts, theme.Muted.Render("  ·  "))
}

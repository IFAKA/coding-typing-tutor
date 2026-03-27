package typing

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	kb "github.com/IFAKA/coding-typing-tutor/internal/keyboard"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

// renderKeyboard renders the keyboard visualizer.
//   - currentChar: the next character the user needs to type
//   - weakKeys:    keys the user historically struggles with (highlighted in red)
//   - wrongFlash:  the expected key during a wrong-key flash (0 = no flash)
func renderKeyboard(currentChar rune, weakKeys map[rune]bool, wrongFlash rune) string {
	base, needsShift := kb.ResolveKey(currentChar)
	af := kb.ActiveFinger(base)

	activeStyle   := lipgloss.NewStyle().Background(theme.Mauve).Foreground(theme.Base).Bold(true)
	wrongStyle    := lipgloss.NewStyle().Background(theme.Red).Foreground(theme.Base).Bold(true)
	shiftOnStyle  := lipgloss.NewStyle().Background(theme.Teal).Foreground(theme.Base).Bold(true)
	shiftOffStyle := lipgloss.NewStyle().Foreground(theme.Surface1)
	spaceOnStyle  := lipgloss.NewStyle().Background(theme.Teal).Foreground(theme.Base).Bold(true)
	spaceOffStyle := lipgloss.NewStyle().Foreground(theme.Surface1)

	renderKey := func(k kb.KeyDef) string {
		label := k.Label()
		isHit        := base != 0 && k.Ch == base
		isWrongFlash := wrongFlash != 0 && k.Ch == wrongFlash
		isSameFinger := af >= 0 && k.F == af && !isHit
		isWeak       := weakKeys[k.Ch] && !isHit && !isSameFinger

		switch {
		case isWrongFlash:
			// Flash the correct key in red after a wrong keypress
			return wrongStyle.Render(label)
		case isHit:
			return activeStyle.Render(label)
		case isSameFinger:
			return lipgloss.NewStyle().Foreground(kb.FingerColor[k.F]).Render(label)
		case isWeak:
			return lipgloss.NewStyle().Foreground(theme.Red).Faint(true).Render(label)
		default:
			return lipgloss.NewStyle().Foreground(kb.FingerColor[k.F]).Faint(true).Render(label)
		}
	}

	var lines []string

	// Hand labels
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Surface1)
	lines = append(lines, "  "+mutedStyle.Render("LEFT")+"               "+mutedStyle.Render("RIGHT"))

	// Finger label strip
	fingerLabels := []struct {
		label string
		f     kb.Finger
	}{
		{"P", kb.LP}, {"R", kb.LR}, {"M", kb.LM}, {"I", kb.LI}, {"I", kb.LI},
		{"I", kb.RI}, {"I", kb.RI}, {"M", kb.RM}, {"R", kb.RR}, {"P", kb.RP},
	}
	var labelSB strings.Builder
	labelSB.WriteString("  ")
	for i, fl := range fingerLabels {
		if i == 5 {
			labelSB.WriteString(" ")
		}
		isActive := af >= 0 && fl.f == af
		if isActive {
			labelSB.WriteString(lipgloss.NewStyle().Foreground(kb.FingerColor[fl.f]).Bold(true).Render(fl.label))
		} else {
			labelSB.WriteString(lipgloss.NewStyle().Foreground(kb.FingerColor[fl.f]).Faint(true).Render(fl.label))
		}
		if i < len(fingerLabels)-1 {
			labelSB.WriteString(" ")
		}
	}
	lines = append(lines, labelSB.String())

	for rowIdx, row := range kb.KbRows {
		var sb strings.Builder

		if rowIdx == 3 {
			sb.WriteString("   ")
			if needsShift {
				sb.WriteString(shiftOnStyle.Render("⇧"))
			} else {
				sb.WriteString(shiftOffStyle.Render("⇧"))
			}
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
			if needsShift {
				sb.WriteString(shiftOnStyle.Render("⇧"))
			} else {
				sb.WriteString(shiftOffStyle.Render("⇧"))
			}
		}

		lines = append(lines, sb.String())
	}

	// Space bar
	const spaceBarLabel = "___________"
	var spaceSB strings.Builder
	spaceSB.WriteString("        ")
	if base == ' ' {
		spaceSB.WriteString(spaceOnStyle.Render(spaceBarLabel))
	} else {
		spaceSB.WriteString(spaceOffStyle.Render(spaceBarLabel))
	}
	lines = append(lines, spaceSB.String())

	return strings.Join(lines, "\n")
}

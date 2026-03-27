package engine

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

// SyntaxColors returns a lipgloss.Color per rune for the given code string.
// Falls back to theme.Text if the lexer cannot process the code.
func SyntaxColors(code, language string) []lipgloss.Color {
	runes := []rune(code)
	colors := make([]lipgloss.Color, len(runes))
	for i := range colors {
		colors[i] = theme.Text
	}

	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iter, err := lexer.Tokenise(nil, code)
	if err != nil {
		return colors
	}

	runeIdx := 0
	for {
		tok := iter()
		if tok == chroma.EOF {
			break
		}
		color := tokenColor(tok.Type)
		for range []rune(tok.Value) {
			if runeIdx < len(colors) {
				colors[runeIdx] = color
			}
			runeIdx++
		}
	}

	return colors
}

// tokenColor maps a chroma token type to a Catppuccin Mocha color.
func tokenColor(t chroma.TokenType) lipgloss.Color {
	tInt := int(t)

	switch {
	// Keywords (1000–1999)
	case tInt >= 1000 && tInt < 2000:
		return theme.Mauve

	// Specific name types
	case t == chroma.NameFunction, t == chroma.NameFunctionMagic:
		return theme.Blue
	case t == chroma.NameBuiltin, t == chroma.NameBuiltinPseudo:
		return theme.Teal
	case t == chroma.NameDecorator:
		return theme.Peach

	// Strings (3100–3199)
	case tInt >= 3100 && tInt < 3200:
		return theme.Green

	// Numbers (3200–3299)
	case tInt >= 3200 && tInt < 3300:
		return theme.Peach

	// Comments (6000–6999)
	case tInt >= 6000 && tInt < 7000:
		return theme.Gray

	// Operators (7000–7999)
	case tInt >= 7000 && tInt < 8000:
		return theme.Sky

	// Punctuation (5000–5999)
	case tInt >= 5000 && tInt < 6000:
		return theme.Subtext

	default:
		return theme.Text
	}
}

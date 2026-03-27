// Package keyboard provides the shared keyboard layout data used by both
// the typing visualizer and the stats heatmap.
package keyboard

import (
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-typing-tutor/internal/theme"
)

// Finger identifies one of 8 finger positions (left pinky → right pinky).
type Finger int

const (
	LP Finger = iota // left pinky
	LR               // left ring
	LM               // left middle
	LI               // left index
	RI               // right index
	RM               // right middle
	RR               // right ring
	RP               // right pinky
)

// FingerColor maps each finger to its display color.
var FingerColor = [8]lipgloss.Color{
	theme.Mauve,  // LP
	theme.Blue,   // LR
	theme.Sky,    // LM
	theme.Teal,   // LI
	theme.Green,  // RI
	theme.Yellow, // RM
	theme.Peach,  // RR
	theme.Red,    // RP
}

// FingerNames maps each finger to a human-readable label.
var FingerNames = [8]string{
	"left pinky", "left ring", "left middle", "left index",
	"right index", "right middle", "right ring", "right pinky",
}

// KeyDef describes a single key on the keyboard.
type KeyDef struct {
	Ch      rune   // base character (lowercase)
	Display rune   // display override (0 = use Ch)
	F       Finger // finger assignment
}

// Label returns the string to display for this key.
func (k KeyDef) Label() string {
	if k.Display != 0 {
		return string(k.Display)
	}
	return string(k.Ch)
}

var (
	Row0 = []KeyDef{
		{'`', 0, LP}, {'1', 0, LP}, {'2', 0, LR}, {'3', 0, LM}, {'4', 0, LI}, {'5', 0, LI},
		{'6', 0, RI}, {'7', 0, RI}, {'8', 0, RM}, {'9', 0, RR}, {'0', 0, RP}, {'-', 0, RP}, {'=', 0, RP},
	}
	Row1 = []KeyDef{
		{'q', 0, LP}, {'w', 0, LR}, {'e', 0, LM}, {'r', 0, LI}, {'t', 0, LI},
		{'y', 0, RI}, {'u', 0, RI}, {'i', 0, RM}, {'o', 0, RR}, {'p', 0, RP},
		{'[', 0, RP}, {']', 0, RP}, {'\\', 0, RP},
	}
	Row2 = []KeyDef{
		{'a', 0, LP}, {'s', 0, LR}, {'d', 0, LM}, {'f', 0, LI}, {'g', 0, LI},
		{'h', 0, RI}, {'j', 0, RI}, {'k', 0, RM}, {'l', 0, RR}, {';', 0, RP}, {'\'', 0, RP},
		{'\n', '↵', RP},
	}
	Row3 = []KeyDef{
		{'z', 0, LP}, {'x', 0, LR}, {'c', 0, LM}, {'v', 0, LI}, {'b', 0, LI},
		{'n', 0, RI}, {'m', 0, RI}, {',', 0, RM}, {'.', 0, RR}, {'/', 0, RP},
	}

	// KbRows is all rows in order.
	KbRows = [][]KeyDef{Row0, Row1, Row2, Row3}

	// GapAfter gives the index after which the left/right hand gap is inserted.
	GapAfter = []int{5, 4, 4, 4}

	// RowIndent is the leading whitespace for each row.
	RowIndent = []string{"", " ", "  ", ""}
)

// ShiftMap maps shifted characters to their base key.
var ShiftMap = map[rune]rune{
	'!': '1', '@': '2', '#': '3', '$': '4', '%': '5',
	'^': '6', '&': '7', '*': '8', '(': '9', ')': '0',
	'_': '-', '+': '=',
	'{': '[', '}': ']', '|': '\\',
	':': ';', '"': '\'',
	'<': ',', '>': '.', '?': '/',
	'~': '`',
}

// ResolveKey returns the base key and whether Shift is needed.
func ResolveKey(ch rune) (base rune, needsShift bool) {
	switch ch {
	case 0, '\t':
		return 0, false
	case '\n':
		return '\n', false
	case ' ':
		return ' ', false
	}
	if unicode.IsUpper(ch) {
		return unicode.ToLower(ch), true
	}
	if b, ok := ShiftMap[ch]; ok {
		return b, true
	}
	return ch, false
}

// ActiveFinger returns the finger responsible for the given base key, or -1.
func ActiveFinger(base rune) Finger {
	for _, row := range KbRows {
		for _, k := range row {
			if k.Ch == base {
				return k.F
			}
		}
	}
	return -1
}

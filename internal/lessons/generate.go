package lessons

import (
	_ "embed"
	"math/rand"
	"strings"
	"unicode"
)

//go:embed data/words.txt
var wordData string

// Generate produces a typing exercise string for the given lesson.
// The result is approximately targetLen characters long.
func Generate(lesson Lesson, targetLen int) string {
	if lesson.AllowedKeys == nil {
		// Free-code / all-symbols: use full word list
		return buildText(allWords(), targetLen)
	}

	allowed := allowedSet(lesson.AllowedKeys)
	words := filterWords(allWords(), allowed)

	if len(words) < 5 {
		// Fallback for very restricted lessons (e.g., home-row only)
		words = syntheticWords(lesson.AllowedKeys, 40)
	}

	return buildText(words, targetLen)
}

// allWords parses the embedded word list.
func allWords() []string {
	lines := strings.Split(strings.TrimSpace(wordData), "\n")
	words := make([]string, 0, len(lines))
	for _, l := range lines {
		w := strings.TrimSpace(l)
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

// allowedSet builds a set of allowed runes (case-insensitive base keys).
func allowedSet(keys []rune) map[rune]bool {
	set := make(map[rune]bool, len(keys)+1)
	set[' '] = true // space always allowed
	for _, r := range keys {
		set[unicode.ToLower(r)] = true
		set[unicode.ToUpper(r)] = true
	}
	return set
}

// filterWords keeps only words whose characters are all in allowed.
func filterWords(words []string, allowed map[rune]bool) []string {
	var out []string
	for _, w := range words {
		ok := true
		for _, r := range w {
			if !allowed[r] {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, w)
		}
	}
	return out
}

// syntheticWords generates pseudo-words from allowed keys for sparse lessons.
func syntheticWords(keys []rune, count int) []string {
	if len(keys) == 0 {
		return []string{"a"}
	}
	words := make([]string, count)
	for i := range words {
		length := 2 + rand.Intn(5)
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteRune(keys[rand.Intn(len(keys))])
		}
		words[i] = sb.String()
	}
	return words
}

// buildText shuffles words and joins them until targetLen chars is reached.
func buildText(words []string, targetLen int) string {
	if len(words) == 0 {
		return ""
	}
	// Shuffle a copy
	pool := make([]string, len(words))
	copy(pool, words)
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	var sb strings.Builder
	for sb.Len() < targetLen {
		for _, w := range pool {
			if sb.Len() > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(w)
			if sb.Len() >= targetLen {
				break
			}
		}
	}
	result := sb.String()
	// Trim to a word boundary near targetLen
	if len(result) > targetLen {
		cut := result[:targetLen]
		if idx := strings.LastIndex(cut, " "); idx > targetLen/2 {
			result = cut[:idx]
		} else {
			result = cut
		}
	}
	return result
}

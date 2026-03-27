package snippets

import (
	"embed"
	"encoding/json"
	"math/rand"
	"sort"
	"time"
)

//go:embed data/*.json
var dataFS embed.FS

var allSnippets []Snippet

func init() {
	files := []string{
		"data/python.json",
		"data/javascript.json",
		"data/typescript.json",
		"data/go.json",
		"data/cpp.json",
		"data/nextjs.json",
	}
	for _, f := range files {
		b, err := dataFS.ReadFile(f)
		if err != nil {
			continue
		}
		var s []Snippet
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		allSnippets = append(allSnippets, s...)
	}
}

// All returns all loaded snippets.
func All() []Snippet {
	return allSnippets
}

// Filter returns snippets matching the given language and difficulty.
// Empty string matches all values for that field.
func Filter(language, difficulty string) []Snippet {
	var result []Snippet
	for _, s := range allSnippets {
		if language != "" && s.Language != language {
			continue
		}
		if difficulty != "" && s.Difficulty != difficulty {
			continue
		}
		result = append(result, s)
	}
	return result
}

// Pick selects a snippet weighted by recency (snippets seen longer ago score higher).
// seenAt maps snippet ID → last time it was played (zero value = never played).
// weakKeys optionally boosts snippets containing keys the user struggles with.
func Pick(language, difficulty string, seenAt map[string]time.Time, weakKeys map[rune]bool) *Snippet {
	pool := Filter(language, difficulty)
	if len(pool) == 0 {
		// Fallback: ignore difficulty filter
		pool = Filter(language, "")
	}
	if len(pool) == 0 {
		return nil
	}

	now := time.Now()
	type scored struct {
		s     Snippet
		score float64
	}

	scored_ := make([]scored, len(pool))
	for i, s := range pool {
		var base float64
		last, ok := seenAt[s.ID]
		if !ok {
			base = 1e9
		} else {
			base = now.Sub(last).Hours()
		}
		// Adaptive bonus: snippets containing more weak keys score higher.
		var bonus float64
		if len(weakKeys) > 0 {
			count := 0
			for _, r := range []rune(s.Code) {
				if weakKeys[r] {
					count++
				}
			}
			total := len([]rune(s.Code))
			if total > 0 {
				bonus = float64(count) / float64(total) * 10.0
			}
		}
		scored_[i] = scored{s, base + bonus}
	}

	// Sort descending by score
	sort.Slice(scored_, func(i, j int) bool {
		return scored_[i].score > scored_[j].score
	})

	// Weighted random pick from top half (or at least 3)
	topN := len(scored_) / 2
	if topN < 3 {
		topN = len(scored_)
	}

	// Assign weights: higher score = higher weight
	totalWeight := 0.0
	for i := 0; i < topN; i++ {
		totalWeight += scored_[i].score
	}

	r := rand.Float64() * totalWeight
	cum := 0.0
	for i := 0; i < topN; i++ {
		cum += scored_[i].score
		if r <= cum {
			return &scored_[i].s
		}
	}
	return &scored_[0].s
}

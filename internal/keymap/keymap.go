package keymap

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// KeyStats holds cumulative attempt/error counts for a single key.
type KeyStats struct {
	Attempts int `json:"attempts"`
	Errors   int `json:"errors"`
}

// KeyDelta is the per-session delta passed from the typing model.
type KeyDelta struct {
	Attempts int
	Errors   int
}

// Store maps rune → cumulative KeyStats.
type Store map[rune]KeyStats

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		home, err2 := os.UserHomeDir()
		if err2 != nil {
			return "", err2
		}
		base = filepath.Join(home, ".config")
	}
	dir := filepath.Join(base, "coding-type")
	return dir, os.MkdirAll(dir, 0755)
}

// Load reads the keymap from disk. Returns empty Store if no file exists.
func Load() (Store, error) {
	dir, err := configDir()
	if err != nil {
		return Store{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "keymap.json"))
	if os.IsNotExist(err) {
		return Store{}, nil
	}
	if err != nil {
		return Store{}, err
	}
	// Keys are stored as single-character strings for human readability.
	var raw map[string]KeyStats
	if err := json.Unmarshal(data, &raw); err != nil {
		return Store{}, nil
	}
	s := make(Store, len(raw))
	for k, v := range raw {
		runes := []rune(k)
		if len(runes) == 1 {
			s[runes[0]] = v
		}
	}
	return s, nil
}

// Save writes the store to disk.
func Save(s Store) error {
	raw := make(map[string]KeyStats, len(s))
	for r, ks := range s {
		raw[string(r)] = ks
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "keymap.json"), data, 0644)
}

// Merge loads the store, applies the delta, and saves.
func Merge(delta map[rune]KeyDelta) error {
	s, _ := Load()
	if s == nil {
		s = make(Store)
	}
	for r, d := range delta {
		ks := s[r]
		ks.Attempts += d.Attempts
		ks.Errors += d.Errors
		s[r] = ks
	}
	return Save(s)
}

// ErrorRate returns the error rate for a key (0.0–1.0).
func ErrorRate(ks KeyStats) float64 {
	if ks.Attempts == 0 {
		return 0
	}
	return float64(ks.Errors) / float64(ks.Attempts)
}

// WeakKeys returns keys whose error rate exceeds the threshold (min 5 attempts).
func WeakKeys(s Store, threshold float64) map[rune]bool {
	weak := make(map[rune]bool)
	for r, ks := range s {
		if ks.Attempts >= 5 && ErrorRate(ks) > threshold {
			weak[r] = true
		}
	}
	return weak
}

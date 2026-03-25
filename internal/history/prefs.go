package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Prefs holds the user's last menu selections.
type Prefs struct {
	LangIdx int `json:"lang_idx"`
	DiffIdx int `json:"diff_idx"`
	ModeIdx int `json:"mode_idx"`
}

func prefsPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "prefs.json"), nil
}

// LoadPrefs reads saved menu preferences. Returns zero-value Prefs if no file exists.
func LoadPrefs() Prefs {
	path, err := prefsPath()
	if err != nil {
		return Prefs{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Prefs{}
	}
	var p Prefs
	if err := json.Unmarshal(data, &p); err != nil {
		return Prefs{}
	}
	return p
}

// SavePrefs writes the current menu preferences to disk.
func SavePrefs(p Prefs) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	path, err := prefsPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

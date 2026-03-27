package lessons

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Progress tracks lesson completion state.
type Progress struct {
	// ConsecutivePass[lessonNum] = number of consecutive sessions with ≥90% accuracy.
	ConsecutivePass map[int]int  `json:"consecutive_pass"`
	// Unlocked[lessonNum] = true if the lesson is available to play.
	Unlocked        map[int]bool `json:"unlocked"`
}

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

// LoadProgress reads lesson progress from disk.
func LoadProgress() Progress {
	p := Progress{
		ConsecutivePass: make(map[int]int),
		Unlocked:        map[int]bool{1: true}, // lesson 1 always unlocked
	}
	dir, err := configDir()
	if err != nil {
		return p
	}
	data, err := os.ReadFile(filepath.Join(dir, "progress.json"))
	if err != nil {
		return p
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p
	}
	if p.ConsecutivePass == nil {
		p.ConsecutivePass = make(map[int]int)
	}
	if p.Unlocked == nil {
		p.Unlocked = map[int]bool{1: true}
	}
	// Lesson 1 is always unlocked
	p.Unlocked[1] = true
	return p
}

// SaveProgress writes lesson progress to disk.
func SaveProgress(p Progress) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "progress.json"), data, 0644)
}

// UpdateProgress records the result of a lesson session.
// Achieving ≥90% accuracy 3 times in a row unlocks the next lesson.
func UpdateProgress(p *Progress, lessonNum int, accuracy float64) {
	if accuracy >= 90.0 {
		p.ConsecutivePass[lessonNum]++
		if p.ConsecutivePass[lessonNum] >= 3 {
			// Unlock next lesson
			next := lessonNum + 1
			if next <= len(All) {
				p.Unlocked[next] = true
			}
		}
	} else {
		p.ConsecutivePass[lessonNum] = 0
	}
}

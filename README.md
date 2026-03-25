# coding-type

A terminal typing tutor for coding interview prep. Practice real algorithm patterns — binary search, BFS, dynamic programming — in Python, JavaScript, TypeScript, Go, and C++.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/IFAKA/coding-type/main/install.sh | sh
```

## Uninstall

```sh
curl -fsSL https://raw.githubusercontent.com/IFAKA/coding-type/main/uninstall.sh | sh
```

## Usage

```sh
coding-type
```

### Controls

| Key | Action |
|-----|--------|
| `←` / `→` | Change selected option |
| `↑` / `↓` or `Tab` | Move between rows |
| `Enter` | Start typing |
| `s` | View stats |
| `Backspace` | Delete last character |
| `Ctrl+R` | Restart current snippet |
| `Esc` | Back to menu |
| `q` | Quit |

## Features

- **5 languages** — Python, JavaScript, TypeScript, Go, C++
- **3 difficulty levels** — easy, medium, hard
- **2 modes** — practice (untimed) or timed (60s)
- **Syntax highlighting** on untyped code (Catppuccin Mocha theme)
- **Live WPM + accuracy** updating as you type
- **Session history** saved to `~/.config/coding-type/history.json`
- **Persistent menu selection** — language, difficulty, and mode are remembered across sessions
- **Smart snippet selection** — avoids recently seen snippets

## Reset

Delete history and preferences:

```sh
rm -rf ~/.config/coding-type
```

Delete only preferences (keep history):

```sh
rm ~/.config/coding-type/prefs.json
```

## Build from source

```sh
go install github.com/IFAKA/coding-type@latest
```

Or:

```sh
git clone https://github.com/IFAKA/coding-type
cd coding-type
go build -o coding-type .
```

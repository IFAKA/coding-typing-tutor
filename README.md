# coding-type: Terminal-based Coding Typing Tutor ⌨️🚀

[![Go Report Card](https://goreportcard.com/badge/github.com/IFAKA/coding-typing-tutor)](https://goreportcard.com/report/github.com/IFAKA/coding-typing-tutor)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/IFAKA/coding-typing-tutor.svg)](https://pkg.go.dev/github.com/IFAKA/coding-typing-tutor)

**`coding-type`** is a high-performance, terminal-based typing tutor designed specifically for software engineers and students preparing for **coding interviews**. Practice real-world **algorithm patterns** and **LeetCode-style problems** directly in your CLI to improve your coding speed, accuracy, and muscle memory.

<p align="center">
  <img width="1512" height="918" alt="coding-type menu" src="https://github.com/user-attachments/assets/32377351-b371-4c44-b2e0-438205527b84" />
</p>

## Why coding-type?

Modern technical interviews aren't just about solving the problem; they're about communicating and implementing your solution efficiently. `coding-type` helps you:
- **Master Algorithm Patterns**: Practice common templates for BFS, DFS, Binary Search, and Dynamic Programming.
- **Improve Muscle Memory**: Type syntax-heavy code in Python, JavaScript, TypeScript, Go, and C++ faster.
- **Stay in the Flow**: Practice directly in your terminal, mirroring a real development environment.
- **Track Progress**: Monitor your WPM (Words Per Minute) and accuracy with a persistent history log.

<p align="center">
  <img width="1512" height="918" alt="coding-type practice session" src="https://github.com/user-attachments/assets/1e5e792e-bcba-46f2-93b1-e63765f6e516" />
</p>

## Features

- **5 Languages Supported**: Python, JavaScript, TypeScript, Go, and C++.
- **Algorithm-Focused Snippets**: Curated snippets covering real-world algorithm patterns.
- **3 Difficulty Levels**: Easy, medium, and hard snippets to challenge your skills.
- **Timed & Practice Modes**: Race against the clock or practice at your own pace.
- **Syntax Highlighting**: Beautiful Catppuccin Mocha theme highlighting to simulate your IDE.
- **Persistent Progress**: Your history and preferences are saved locally at `~/.config/coding-type/`.
- **Smart Selection**: The engine avoids showing you the same snippet twice in a row.

## Installation

### Quick Install (Shell)

```sh
curl -fsSL https://raw.githubusercontent.com/IFAKA/coding-type/main/install.sh | sh
```

### Go Install

```sh
go install github.com/IFAKA/coding-typing-tutor@latest
```

## Usage

Simply run:
```sh
coding-type
```

### Controls

| Key | Action |
|-----|--------|
| `←` / `→` | Change selected option |
| `↑` / `↓` or `Tab` | Move between rows |
| `Enter` | Start typing |
| `s` | View session stats |
| `Backspace` | Delete last character |
| `Ctrl+R` | Restart current snippet |
| `Esc` | Back to main menu |
| `q` | Quit application |

## Development & Build

If you want to build from source:

```sh
git clone https://github.com/IFAKA/coding-typing-tutor
cd coding-type
go build -o coding-type .
```

Or run directly without building:

```sh
go run .
```

## Contributing

Contributions are welcome! If you have more algorithm snippets or feature ideas, feel free to open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Made for developers who want to type as fast as they think.*

# venv-killer

**venv-killer** is a blazing fast CLI tool written in Go that recursively finds and kills (deletes) Python virtual environments. It acts like `npkill` but for Python `venv` folders.

It identifies virtual environments by looking for `pyvenv.cfg` files or `bin/activate` scripts.

## Features

- 🚀 **Fast**: Recursively scans directories to find venvs.
- 🖥️ **Interactive TUI**: Select which venvs to delete using a simple terminal interface.
- 📊 **Progress Bar**: Visual feedback during deletion.
- 🛡️ **Safe**: Asks for confirmation (via selection) before deleting.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap your-github-username/tap
brew install venv-killer
```

### Go Install

```bash
go install github.com/matiassotose/venv-killer@latest
```

### Manual

Download the latest binary from the [Releases](https://github.com/your-github-username/venv-killer/releases) page.

## Usage

Run the tool in the current directory:

```bash
venv-killer
```

Or specify a directory to search:

```bash
venv-killer /path/to/search
```

### Controls

- **Up / Down / k / j**: Navigate the list.
- **Space**: Select/Deselect a venv.
- **Enter**: Delete selected venvs.
- **q / Ctrl+c**: Quit.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
Built with ❤️ using [Bubbletea](https://github.com/charmbracelet/bubbletea).

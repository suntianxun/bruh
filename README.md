# Bruh 🍺

A beautiful, animated, and highly functional Terminal User Interface (TUI) for managing your Homebrew packages. Built with Go and the [Charm](https://charm.sh) ecosystem.

![Bruh TUI](https://github.com/suntianxun/bruh/blob/main/internal/ui/components/bruh_logo.png)

## Features

- **Beautiful UI:** Catppuccin Mocha themed with gradients, ASCII art header, and a physics-based spring animation.
- **Unified Data Table:** View all your installed Formulae and Casks in a clean, aligned, and automatically truncated data table.
- **Live Filtering:** Press `/` to seamlessly filter through your installed packages as you type.
- **Remote Search:** Press `s` to pull up a modal to search the remote Homebrew catalog. Highlights already-installed packages in green.
- **Quick Actions:** Easily `[u]pgrade`, `[d]ninstall`, or `[r]einstall` packages with a single keystroke.
- **Live Logs:** Watch real-time streaming `brew` output right in the TUI when installing or upgrading packages, so you never have to wonder if it's stuck.

## Installation

You can install `bruh` via Homebrew:

```bash
brew install suntianxun/tap/bruh
```

Or build from source:

```bash
git clone https://github.com/suntianxun/bruh.git
cd bruh
go build .
./bruh
```

## Shortcuts

**Main View:**
- `/` - Filter installed packages
- `s` - Open remote search popup
- `u` - Upgrade the selected package
- `d` - Uninstall the selected package
- `r` - Reinstall the selected package
- `q` / `Ctrl+C` - Quit

**Search View:**
- `Enter` - Search for the typed query
- `i` - Install the selected package
- `d` - Uninstall the selected package
- `q` / `Esc` - Close search and return to main view

## License

MIT License. See [LICENSE](LICENSE) for more details.

// internal/ui/theme.go
package ui

import "github.com/charmbracelet/lipgloss"

var (
	MochaMauve    = lipgloss.Color("#cba6f7")
	MochaPink     = lipgloss.Color("#f5c2e7")
	MochaFlamingo = lipgloss.Color("#f2cdcd")
	MochaBase     = lipgloss.Color("#1e1e2e")
	MochaText     = lipgloss.Color("#cdd6f4")
	MochaSurface0 = lipgloss.Color("#313244")
)

var (
	HeaderStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true).
		Foreground(MochaBase).
		Background(MochaMauve)
)
// internal/ui/components/tabs.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#cba6f7")). // MochaMauve
		Foreground(lipgloss.Color("#cba6f7")).
		Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#313244")). // MochaSurface0
		Foreground(lipgloss.Color("#a6adc8")).       // MochaSubtext0
		Padding(0, 1)
)

func RenderTabs(tabs []string, activeIdx int, width int) string {
	var renderedTabs []string

	for i, t := range tabs {
		if i == activeIdx {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(t))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	fillerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#313244")) // MochaSurface0
	
	fillerWidth := width - lipgloss.Width(row)
	if fillerWidth < 0 {
		fillerWidth = 0
	}
	filler := fillerStyle.Render(strings.Repeat(" ", fillerWidth))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, filler)
}
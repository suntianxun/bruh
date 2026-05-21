// internal/ui/components/details.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
)

func RenderDetails(pkg brew.PackageInfo, width int) string {
	md := fmt.Sprintf("# %s\n\n", pkg.Name)
	if pkg.LatestVersion != "" {
		md += fmt.Sprintf("**Version:** %s\n\n", pkg.LatestVersion)
	}
	md += fmt.Sprintf("%s\n\n", pkg.Desc)

	if len(pkg.Dependencies) > 0 {
		md += "## Dependencies\n"
		md += strings.Join(pkg.Dependencies, ", ") + "\n"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width - 4),
	)
	
	out, err := r.Render(md)
	if err != nil {
		return "Error rendering details"
	}
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		Padding(1, 2).
		Width(width).
		Render(out)
}
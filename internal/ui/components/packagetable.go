// internal/ui/components/packagetable.go
package components

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
)

func NewPackageTable(pkgs []brew.PackageInfo, width, height int) table.Model {
	columns := []table.Column{
		{Title: "Package", Width: 25},
		{Title: "Description", Width: width - 25 - 15 - 15 - 10 - 10}, // Calculate remaining space
		{Title: "Installed", Width: 15},
		{Title: "Latest", Width: 15},
		{Title: "Type", Width: 10},
	}

	// Ensure description has a minimum width
	if columns[1].Width < 20 {
		columns[1].Width = 20
	}

	var rows []table.Row
	for _, p := range pkgs {
		pkgType := "Formula"
		if p.IsCask {
			pkgType = "Cask"
		}
		
		desc := p.Desc
		if len(desc) > columns[1].Width {
			desc = desc[:columns[1].Width-3] + "..."
		}

		installed := p.CurrentVersion
		if installed == "" {
			installed = "-"
		}

		latest := p.LatestVersion
		if latest == "" {
			latest = "-"
		}

		rows = append(rows, table.Row{
			p.Name,
			desc,
			installed,
			latest,
			pkgType,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height-4),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#313244")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#1e1e2e")).
		Background(lipgloss.Color("#cba6f7")).
		Bold(false)
		
	t.SetStyles(s)

	return t
}
// internal/ui/components/actions.go
package components

import (
	"github.com/charmbracelet/huh"
	"github.com/user/brew-tui/internal/brew"
)

func NewConfirmForm(action string, pkg brew.PackageInfo) *huh.Form {
	var confirm bool
	
	// Default theme matches lipgloss relatively well, or we can use huh.ThemeCatppuccin()
	theme := huh.ThemeCatppuccin()
	
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(action + " " + pkg.Name + "?").
				Value(&confirm),
		),
	).WithTheme(theme)
}
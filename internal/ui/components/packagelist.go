// internal/ui/components/packagelist.go
package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/user/brew-tui/internal/brew"
)

type PackageItem brew.PackageInfo

func (i PackageItem) Title() string       { return i.Name }
func (i PackageItem) Description() string { return i.Desc }
func (i PackageItem) FilterValue() string { return i.Name }

func NewList(pkgs []brew.PackageInfo, width, height int) list.Model {
	var items []list.Item
	for _, p := range pkgs {
		items = append(items, PackageItem(p))
	}

	m := list.New(items, list.NewDefaultDelegate(), width, height)
	m.Title = "Installed Packages"
	return m
}
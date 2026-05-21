// internal/ui/components/packagetable.go
package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/user/brew-tui/internal/brew"
)

var (
	colorOutdated = lipgloss.Color("#f38ba8") // Mocha Red
	colorUpToDate = lipgloss.Color("#bac2de") // Mocha Subtext1 (softer)
	colorSelected = lipgloss.Color("#cba6f7") // Mocha Mauve text
	colorSelectedBg = lipgloss.Color("#313244") // Mocha Surface0 subtle background
	colorText     = lipgloss.Color("#cdd6f4")
	colorSurface  = lipgloss.Color("#313244")
)

type PackageItem brew.PackageInfo

func (i PackageItem) Title() string       { return i.Name }
func (i PackageItem) Description() string { return i.Desc }
func (i PackageItem) FilterValue() string { return i.Name }

type packageDelegate struct {
	width int
}

func (d packageDelegate) Height() int                             { return 1 }
func (d packageDelegate) Spacing() int                            { return 0 }
func (d packageDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d packageDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PackageItem)
	if !ok {
		return
	}

	col1W := 25
	col3W := 15
	col4W := 15
	col5W := 10
	col2W := d.width - col1W - col3W - col4W - col5W - 8 // padding
	if col2W < 20 {
		col2W = 20
	}

	pkgType := "Formula"
	if i.IsCask {
		pkgType = "Cask"
	}

	installed := i.CurrentVersion
	if installed == "" {
		installed = "-"
	}

	latest := i.LatestVersion
	if latest == "" {
		latest = "-"
	}

	// Truncate plain strings BEFORE applying colors
	name := runewidth.Truncate(i.Name, col1W, "…")
	desc := runewidth.Truncate(i.Desc, col2W, "…")
	inst := runewidth.Truncate(installed, col3W, "…")
	lat := runewidth.Truncate(latest, col4W, "…")
	ptype := runewidth.Truncate(pkgType, col5W, "…")

	var fg lipgloss.Color = colorText
	if i.CurrentVersion != "" {
		fg = colorUpToDate
		if i.IsOutdated {
			fg = colorOutdated
		}
	}

	bg := lipgloss.Color("")
	bold := false
	if index == m.Index() {
		fg = colorSelected
		bg = colorSelectedBg
		bold = true
	}

	style := lipgloss.NewStyle().Foreground(fg).Background(bg).Bold(bold)
	
	// Format strings to exact width using left-alignment and spaces
	fmt1 := fmt.Sprintf("%%-%ds", col1W)
	fmt2 := fmt.Sprintf("%%-%ds", col2W)
	fmt3 := fmt.Sprintf("%%-%ds", col3W)
	fmt4 := fmt.Sprintf("%%-%ds", col4W)
	fmt5 := fmt.Sprintf("%%-%ds", col5W)

	col1 := fmt.Sprintf(fmt1, name)
	col2 := fmt.Sprintf(fmt2, desc)
	col3 := fmt.Sprintf(fmt3, inst)
	col4 := fmt.Sprintf(fmt4, lat)
	col5 := fmt.Sprintf(fmt5, ptype)

	row := fmt.Sprintf("%s %s %s %s %s", col1, col2, col3, col4, col5)
	
	// Render row width with background spanning full table width
	fullRow := style.Render(runewidth.FillRight(row, d.width))
	fmt.Fprint(w, fullRow)
}

func RenderTableHeader(width int) string {
	col1W := 25
	col3W := 15
	col4W := 15
	col5W := 10
	col2W := width - col1W - col3W - col4W - col5W - 8
	if col2W < 20 {
		col2W = 20
	}

	s := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#b4befe")). // Mocha Sapphire
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#585b70")). // Mocha Surface2
		PaddingBottom(0)
	
	fmt1 := fmt.Sprintf("%%-%ds", col1W)
	fmt2 := fmt.Sprintf("%%-%ds", col2W)
	fmt3 := fmt.Sprintf("%%-%ds", col3W)
	fmt4 := fmt.Sprintf("%%-%ds", col4W)
	fmt5 := fmt.Sprintf("%%-%ds", col5W)

	col1 := fmt.Sprintf(fmt1, "PACKAGE")
	col2 := fmt.Sprintf(fmt2, "DESCRIPTION")
	col3 := fmt.Sprintf(fmt3, "VERSION")
	col4 := fmt.Sprintf(fmt4, "LATEST")
	col5 := fmt.Sprintf(fmt5, "TYPE")

	row := fmt.Sprintf("%s %s %s %s %s", col1, col2, col3, col4, col5)
	return s.Render(runewidth.FillRight(row, width))
}

func NewPackageTable(pkgs []brew.PackageInfo, width, height int) list.Model {
	var items []list.Item
	for _, p := range pkgs {
		items = append(items, PackageItem(p))
	}

	d := packageDelegate{width: width}
	l := list.New(items, d, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Filter using /
	l.SetShowHelp(false)
	
	return l
}
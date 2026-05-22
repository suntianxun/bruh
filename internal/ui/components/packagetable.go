// internal/ui/components/packagetable.go
package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/user/brew-tui/internal/brew"
)

var (
	colorOutdated = lipgloss.Color("#f38ba8") // Mocha Red
	colorUpToDate = lipgloss.Color("#bac2de") // Mocha Subtext1 (softer)
	colorSelected = lipgloss.Color("#f9e2af") // Mocha Yellow (Golden)
	colorSelectedBg = lipgloss.Color("#313244") // Mocha Surface0 subtle background
	colorText     = lipgloss.Color("#cdd6f4")
	colorSurface  = lipgloss.Color("#313244")
)

type PackageItem brew.PackageInfo

func (i PackageItem) Title() string       { return i.Name }
func (i PackageItem) Description() string { return i.Desc }
func (i PackageItem) FilterValue() string { return i.Name }

type packageDelegate struct {
	width    int
	isSearch bool
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
		if d.isSearch {
			fg = lipgloss.Color("#a6e3a1") // Mocha Green
		}
		if i.IsOutdated {
			fg = colorOutdated
		}
	} else if i.IsOutdated {
		fg = colorOutdated
	}

	bg := lipgloss.Color("")
	bold := false

	if index == m.Index() {
		fg = colorSelected
		bg = colorSelectedBg
		bold = true
	}

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
	row = runewidth.FillRight(row, d.width)

	fullRow := lipgloss.NewStyle().Foreground(fg).Background(bg).Bold(bold).Render(row)
	
	fmt.Fprint(w, fullRow)
}

func applyGradient(text string, bold bool, bg lipgloss.Color) string {
	colors := []lipgloss.Color{
		lipgloss.Color("#cba6f7"),
		lipgloss.Color("#f5c2e7"),
		lipgloss.Color("#eba0ac"),
		lipgloss.Color("#fab387"),
		lipgloss.Color("#f9e2af"),
	}

	var coloredRow string
	for i, c := range text {
		colorIdx := (i * len(colors)) / len(text)
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		style := lipgloss.NewStyle().Foreground(colors[colorIdx]).Bold(bold)
		if bg != "" {
			style = style.Background(bg)
		}
		coloredRow += style.Render(string(c))
	}
	return coloredRow
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

	fmt1 := fmt.Sprintf("%%-%ds", col1W)
	fmt2 := fmt.Sprintf("%%-%ds", col2W)
	fmt3 := fmt.Sprintf("%%-%ds", col3W)
	fmt4 := fmt.Sprintf("%%-%ds", col4W)
	fmt5 := fmt.Sprintf("%%-%ds", col5W)

	col1 := fmt.Sprintf(fmt1, applyGradient("Package", true, ""))
	col2 := fmt.Sprintf(fmt2, applyGradient("Description", true, ""))
	col3 := fmt.Sprintf(fmt3, applyGradient("Version", true, ""))
	col4 := fmt.Sprintf(fmt4, applyGradient("Latest", true, ""))
	col5 := fmt.Sprintf(fmt5, applyGradient("Type", true, ""))

	// When using sprintf with lipgloss styled strings, the length calculations get messed up 
	// by the ANSI escape codes. Instead, we'll pad the strings correctly without ansi codes,
	// and then apply the gradient to just the words.

	col1Padded := runewidth.FillRight("Package", col1W)
	col2Padded := runewidth.FillRight("Description", col2W)
	col3Padded := runewidth.FillRight("Version", col3W)
	col4Padded := runewidth.FillRight("Latest", col4W)
	col5Padded := runewidth.FillRight("Type", col5W)
	
	// Apply gradient to the text, then append the spaces
	col1 = applyGradient("Package", true, "") + col1Padded[len("Package"):]
	col2 = applyGradient("Description", true, "") + col2Padded[len("Description"):]
	col3 = applyGradient("Version", true, "") + col3Padded[len("Version"):]
	col4 = applyGradient("Latest", true, "") + col4Padded[len("Latest"):]
	col5 = applyGradient("Type", true, "") + col5Padded[len("Type"):]

	rowText := fmt.Sprintf("%s %s %s %s %s", col1, col2, col3, col4, col5)
	
	// We can't use FillRight on rowText directly because it contains ansi codes.
	// Since we padded each column exactly, we just need to ensure the total width.
	// The total visible length should be col1W + col2W + col3W + col4W + col5W + 4 spaces
	visibleLen := col1W + col2W + col3W + col4W + col5W + 4
	if visibleLen < width {
		rowText += strings.Repeat(" ", width-visibleLen)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#585b70")). // Mocha Surface2
		PaddingBottom(0).
		Render(rowText)
}

func NewPackageTable(pkgs []brew.PackageInfo, width, height int, isSearch bool) list.Model {
	var items []list.Item
	for _, p := range pkgs {
		items = append(items, PackageItem(p))
	}

	d := packageDelegate{width: width, isSearch: isSearch}
	l := list.New(items, d, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Filter using /
	l.SetShowHelp(false)
	
	return l
}
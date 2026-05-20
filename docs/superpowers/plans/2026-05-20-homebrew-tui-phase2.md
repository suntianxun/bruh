# Homebrew TUI Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Tabbed Navigation, Remote Search View, and Action/Progress Overlay to complete the Homebrew TUI features defined in the design spec.

**Architecture:** Extend the existing Bubble Tea Model to manage state for multiple views (tabs). Extract the active view rendering logic and wire up the `textinput` bubble for the search functionality. Add a spinner/viewport overlay for blocking command execution.

**Tech Stack:** Go, Charm Bubble Tea, Lip Gloss, Bubbles (list, textinput, spinner, viewport), Brew CLI.

---

## File Structure

- Create: `internal/ui/components/tabs.go` - Tab bar rendering component.
- Create: `internal/ui/components/search.go` - Search input and remote results list.
- Create: `internal/ui/components/progress.go` - Spinner and status message overlay for blocking commands.
- Modify: `internal/ui/model.go` - Add tab state management and route messages to active components.
- Modify: `internal/brew/client.go` - Add remote search capabilities (`brew search`).

## Tasks

### Task 1: Tab Navigation UI

**Files:**
- Create: `internal/ui/components/tabs.go`
- Modify: `internal/ui/model.go`

- [ ] **Step 1: Create Tabs component**

```go
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
```

- [ ] **Step 2: Update Main Model to handle Tab switching**

We need to edit `internal/ui/model.go` to use the tabs. Since we can't reliably `sed` this file with high confidence in the exact lines, we will rewrite the Update and View methods. Use `Edit` tool or `Write` tool to replace the `Update` and `View` methods in `internal/ui/model.go`.

```go
// Replace Update and View in internal/ui/model.go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgList.SetSize(msg.Width, m.height-6) // account for header, tabs, footer
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.tab = (m.tab + 1) % 3
			return m, nil
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + 3) % 3
			return m, nil
		}
	case pkgsLoadedMsg:
		m.pkgList = components.NewList(msg, m.width, m.height-6)
	}
	
	if m.tab == tabInstalled {
		m.pkgList, cmd = m.pkgList.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	header := HeaderStyle.Render("Homebrew TUI")
	tabs := []string{"Installed", "Search", "Updates"}
	renderedTabs := components.RenderTabs(tabs, int(m.tab), m.width)
	
	var mainContent string
	switch m.tab {
	case tabInstalled:
		mainContent = m.pkgList.View()
	case tabSearch:
		mainContent = "Search functionality coming soon..."
	case tabUpdates:
		mainContent = "Updates functionality coming soon..."
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, renderedTabs, mainContent)
}
```

- [ ] **Step 3: Run and test**
Run: `go run .`
Expected: TUI displays tabs underneath the header. Pressing Tab/Right Arrow switches between views.

- [ ] **Step 4: Commit**

```bash
git add internal/ui/
git commit -m "feat: implement tabbed navigation"
```

### Task 2: Brew Client Search Functionality

**Files:**
- Modify: `internal/brew/client.go`

- [ ] **Step 1: Add Search Function**

```go
// Add to internal/brew/client.go
import "strings" // Add to imports if not there

func Search(query string) ([]PackageInfo, error) {
	if query == "" {
		return nil, nil
	}
	// Brew search returns simple text lines, not json
	cmd := exec.Command("brew", "search", query)
	out, err := cmd.Output()
	if err != nil {
		// brew search returns non-zero if no results found
		return nil, nil
	}

	var results []PackageInfo
	lines := strings.Split(string(out), "\n")
	
	isCaskSection := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "==> Formulae" {
			isCaskSection = false
			continue
		}
		if line == "==> Casks" {
			isCaskSection = true
			continue
		}
		
		// It's a package name
		results = append(results, PackageInfo{
			Name: line,
			Desc: "Press enter to view details", // We don't get descriptions from basic search
			IsCask: isCaskSection,
		})
	}
	return results, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/brew/client.go
git commit -m "feat: add remote package search to brew client"
```

### Task 3: Search Component UI

**Files:**
- Create: `internal/ui/components/search.go`
- Modify: `internal/ui/model.go`

- [ ] **Step 1: Get textinput bubble**

```bash
go get github.com/charmbracelet/bubbles/textinput
```

- [ ] **Step 2: Create Search Model**

```go
// internal/ui/components/search.go
package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
)

type SearchModel struct {
	Input   textinput.Model
	Results list.Model
	width   int
	height  int
}

func NewSearchModel(width, height int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search for a package..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	l := NewList(nil, width, height-3)
	l.Title = "Results"

	return SearchModel{
		Input:   ti,
		Results: l,
		width:   width,
		height:  height,
	}
}

type SearchResultsMsg []brew.PackageInfo

func SearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		res, err := brew.Search(query)
		if err != nil {
			return err // In a real app we'd wrap this
		}
		return SearchResultsMsg(res)
	}
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.Results.SetSize(msg.Width, msg.Height-3)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Input.Focused() {
				m.Input.Blur()
				return m, SearchCmd(m.Input.Value())
			}
		case tea.KeyEsc:
			m.Input.Focus()
		}
	case SearchResultsMsg:
		m.Results = NewList(msg, m.width, m.height-3)
		m.Results.Title = "Results for: " + m.Input.Value()
	}

	if m.Input.Focused() {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.Results, cmd = m.Results.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m SearchModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(1, 2).Render(m.Input.View()),
		m.Results.View(),
	)
}
```

- [ ] **Step 3: Integrate Search into Main Model**

Edit `internal/ui/model.go` to add `search  components.SearchModel` to the `Model` struct.
Initialize it in `InitialModel()`: `search: components.NewSearchModel(20, 10),`
Update `Update()` and `View()` to route to the search model when `m.tab == tabSearch`.

```go
// Replace struct and Init in internal/ui/model.go
type Model struct {
	tab      activeTab
	pkgList  list.Model
	search   components.SearchModel
	width    int
	height   int
}

func InitialModel() Model {
	return Model{
		tab:     tabInstalled,
		pkgList: components.NewList(nil, 20, 10),
		search:  components.NewSearchModel(20, 10),
	}
}
```

```go
// Modify Update in internal/ui/model.go
// Inside Update switch for WindowSizeMsg:
		m.search, _ = m.search.Update(msg) // pass window size down
// Inside Update at the end:
	if m.tab == tabInstalled {
		m.pkgList, cmd = m.pkgList.Update(msg)
	} else if m.tab == tabSearch {
		m.search, cmd = m.search.Update(msg)
	}
```

```go
// Modify View in internal/ui/model.go
	case tabSearch:
		mainContent = m.search.View()
```

- [ ] **Step 4: Run and Test**
Run `go run .`. Switch to Search tab, type something like "wget", press Enter. Wait a few seconds for results to load.

- [ ] **Step 5: Commit**

```bash
git add internal/ui/
git commit -m "feat: implement search UI component"
```

### Task 4: Command Progress Overlay (Spinner)

**Files:**
- Create: `internal/ui/components/progress.go`
- Modify: `internal/ui/model.go`

- [ ] **Step 1: Get spinner bubble**

```bash
go get github.com/charmbracelet/bubbles/spinner
```

- [ ] **Step 2: Create Progress Component**

```go
// internal/ui/components/progress.go
package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressModel struct {
	Spinner spinner.Model
	Active  bool
	Message string
}

func NewProgressModel() ProgressModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")) // MochaMauve
	return ProgressModel{
		Spinner: s,
	}
}

func (m ProgressModel) Update(msg tea.Msg) (ProgressModel, tea.Cmd) {
	if !m.Active {
		return m, nil
	}
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m ProgressModel) View() string {
	if !m.Active {
		return ""
	}
	overlay := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		Padding(1, 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, m.Spinner.View(), " ", m.Message))
		
	// For simplicity in this plan, we'll just render it as a block rather than true absolute positioning
	return overlay
}
```

- [ ] **Step 3: Integrate Progress into Main Model**

Edit `internal/ui/model.go`. Add `progress components.ProgressModel` to `Model`. Initialize it.
Update `fetchInstalledCmd` to start the spinner, and intercept `pkgsLoadedMsg` to stop it.

```go
// Replace fetchInstalledCmd in internal/ui/model.go
func fetchInstalledCmd() tea.Cmd {
	return tea.Batch(
		components.NewProgressModel().Spinner.Tick,
		func() tea.Msg {
			res, err := brew.GetInstalled()
			if err != nil {
				return errMsg(err)
			}
			return pkgsLoadedMsg(res)
		},
	)
}
```

Add `m.progress.Active = true` / `.Message = "Loading packages..."` in `InitialModel`.
In `Update`, under `pkgsLoadedMsg`, set `m.progress.Active = false`.
Also route generic messages to `m.progress.Update(msg)` at the top of `Update`.
In `View()`, prepend or append `m.progress.View()` to `mainContent`.

- [ ] **Step 4: Commit**

```bash
git add internal/ui/
git commit -m "feat: add progress overlay component"
```
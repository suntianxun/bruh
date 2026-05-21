// internal/ui/components/progress.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressModel struct {
	Spinner spinner.Model
	Active  bool
	Message string
	Logs    []string
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

	logLines := strings.Join(m.Logs, "\n")
	if logLines != "" {
		logLines = "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Render(logLines)
	}

	overlay := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		Padding(1, 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, m.Spinner.View(), " ", m.Message) + logLines)

	return overlay
}
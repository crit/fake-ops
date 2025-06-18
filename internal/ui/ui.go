package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crit/fake-ops/internal/app"
)

type QuitSignal struct{}

// DelayedQuit delays the QuitSignal so that any UI updates can
// happen before the binary is stopped.
func DelayedQuit() tea.Msg {
	time.Sleep(1 * time.Second)
	return QuitSignal{}
}

// UI handles creating all parts of the user interface.
type UI struct {
	ctx            *app.Context
	cancel         func()
	width          int
	height         int
	logs           *LogView
	services       *ServiceView
	containerStyle lipgloss.Style
	footerStyle    lipgloss.Style
}

// Init handles startup items.
func (ui *UI) Init() tea.Cmd {
	return func() tea.Msg {
		// Clear the terminal screen
		fmt.Print("\033[H\033[2J")
		return nil
	}
}

// Update UI with data needed by View.
func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for the keys "q" or "Ctrl+C" to quit
		switch msg.String() {
		case "q", "ctrl+c":
			ui.cancel()
			return ui, DelayedQuit
		}
		return ui, nil // Ignore all other keys

	case QuitSignal:
		return ui, tea.Quit

	case tea.WindowSizeMsg:
		colw, colh := (msg.Width-4)/2, msg.Height-6

		ui.width = msg.Width
		ui.height = msg.Height
		ui.services.Resize(colw, colh)
		ui.logs.Resize(colw, colh)
		return ui, nil

	case app.Message:
		ui.logs.Update(msg)
		return ui, nil

	case app.ServiceMessage:
		ui.services.Update(msg)
		return ui, nil

	case app.ServiceStatus:
		ui.services.UpdateStatus(msg)
		return ui, nil

	default:
		return ui, nil
	}
}

// View combines all of UI's data into a string for display.
func (ui *UI) View() string {
	container := ui.containerStyle.Width(ui.width - 2).Height(ui.height - 2)
	row := lipgloss.JoinHorizontal(lipgloss.Top, ui.services.View(), ui.logs.View())
	footer := ui.footerStyle.Width(ui.width - 4).Render("q or ctrl+c to quit")

	return container.Render(lipgloss.JoinVertical(lipgloss.Top, row, footer))
}

// SetContext makes the app.Context available to the UI.
func (ui *UI) SetContext(ctx *app.Context) {
	ui.ctx = ctx
}

// SetCancel makes the context.Context cancel function available to the UI.
func (ui *UI) SetCancel(cancel func()) {
	ui.cancel = cancel
}

// New creates the UI with appropriate defaults.
func New() *UI {
	var ui UI
	ui.logs = NewLogView()
	ui.services = NewServiceView()

	ui.containerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1)

	ui.footerStyle = lipgloss.NewStyle().
		Foreground(cSecondary).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(cBorder).
		Height(1).
		Align(lipgloss.Right)

	return &ui
}

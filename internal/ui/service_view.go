package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/crit/fake-ops/internal/app"
)

// ServiceView handles creating the services portion of the UI.
type ServiceView struct {
	position     map[string]int
	services     []app.ServiceMessage
	width        int
	height       int
	offlineStyle lipgloss.Style
	onlineStyle  lipgloss.Style
	errorStyle   lipgloss.Style
	gridStyle    lipgloss.Style
}

// NewServiceView creates a ServiceView with appropriate defaults.
func NewServiceView() *ServiceView {
	blockStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, true, true).
		Width(25).
		Height(1)

	return &ServiceView{
		position:     make(map[string]int),
		offlineStyle: blockStyle.Foreground(cSecondary).BorderForeground(cOffline),
		onlineStyle:  blockStyle.Foreground(cOnline).BorderForeground(cOnline),
		errorStyle:   blockStyle.Foreground(cDanger).BorderForeground(cDanger),
		gridStyle:    lipgloss.NewStyle().Foreground(cPrimary),
	}
}

// Update ServiceView with data needed for View.
func (v *ServiceView) Update(msg app.ServiceMessage) {
	if _, ok := v.position[msg.Name]; ok {
		return // already in the slice
	}

	v.services = append(v.services, msg)

	sort.Slice(v.services, func(i, j int) bool {
		return v.services[i].Port < v.services[j].Port
	})

	for i, service := range v.services {
		v.position[service.Name] = i
	}
}

// UpdateStatus handles changes in status based on app.ServiceStatus.
func (v *ServiceView) UpdateStatus(msg app.ServiceStatus) {
	pos, ok := v.position[msg.Name]
	if !ok {
		return // service not registered
	}

	service := v.services[pos]

	if msg.Sent.Before(service.LastStatus) {
		return
	}

	service.Status = msg.Status
	service.LastStatus = msg.Sent
	v.services[pos] = service
}

// Resize sets the width/height of the ServiceView.
func (v *ServiceView) Resize(width, height int) {
	v.width = width
	v.height = height
}

// View combines all of ServiceView's data into a string used by the parent View.
func (v *ServiceView) View() string {
	gridStyle := v.gridStyle.Height(v.height).Width(v.width)

	var blocks []string
	for _, svc := range v.services {
		var block lipgloss.Style
		switch svc.Status {
		case "online":
			block = v.onlineStyle
		case "error":
			block = v.errorStyle
		default:
			block = v.offlineStyle
		}

		var icon string
		switch svc.Kind {
		case "app":
			icon = iCommand
		case "http":
			icon = iCloud
		default:
			icon = iGlobe
		}

		blocks = append(blocks, block.Render(fmt.Sprintf(" %s  %s:%d", icon, svc.Name, svc.Port)))
	}

	return gridStyle.Render(lipgloss.JoinVertical(lipgloss.Top, blocks...))
}

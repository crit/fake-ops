package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/crit/fake-ops/internal/app"
)

type ServiceView struct {
	width    int
	height   int
	position map[string]int
	services []app.ServiceMessage
}

func NewServiceView() *ServiceView {
	return &ServiceView{
		position: make(map[string]int),
	}
}

func (g *ServiceView) Update(msg app.ServiceMessage) {
	if _, ok := g.position[msg.Name]; ok {
		return // already in the slice
	}

	g.services = append(g.services, msg)

	sort.Slice(g.services, func(i, j int) bool {
		return g.services[i].Port < g.services[j].Port
	})

	for i, service := range g.services {
		g.position[service.Name] = i
	}
}

func (g *ServiceView) UpdateStatus(msg app.ServiceStatus) {
	pos, ok := g.position[msg.Name]
	if !ok {
		return // service not registered
	}

	service := g.services[pos]

	if msg.Sent.Before(service.LastStatus) {
		return
	}

	service.Status = msg.Status
	service.LastStatus = msg.Sent
	g.services[pos] = service
}

func (g *ServiceView) Resize(width, height int) {
	g.width = width
	g.height = height
}

func (g *ServiceView) View() string {
	blockStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, true, true).
		Width(25).
		Height(1)

	offlineStyle := blockStyle.Foreground(cSecondary).BorderForeground(cOffline)
	onlineStyle := blockStyle.Foreground(cOnline).BorderForeground(cOnline)
	errorStyle := blockStyle.Foreground(cDanger).BorderForeground(cDanger)

	gridStyle := lipgloss.NewStyle().
		Foreground(cPrimary).
		Height(g.height).        // Use dynamic height
		Width((g.width - 4) / 2) // Use dynamic width

	var blocks []string
	for _, svc := range g.services {
		var block lipgloss.Style
		switch svc.Status {
		case "online":
			block = onlineStyle
		case "error":
			block = errorStyle
		default:
			block = offlineStyle
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

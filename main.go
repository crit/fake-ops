package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crit/fake-ops/internal/app"
	"github.com/crit/fake-ops/internal/services"
	"github.com/crit/fake-ops/internal/ui"
	"github.com/gin-gonic/gin"
)

// ./main --services=./services --results=./results
func main() {
	// silence gin's debug messages
	gin.SetMode(gin.ReleaseMode)

	// create the model and program
	model := ui.New()
	p := tea.NewProgram(model)

	// create the context and register its parts with the model
	// for use by the program.
	ctx, cancel := app.NewContext(p.Send)
	model.SetContext(ctx)
	model.SetCancel(cancel)

	// background all services
	go startServices(ctx)

	// Run blocks until some service or the UI calls the cancel function.
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}

func startServices(ctx *app.Context) {
	// tell me what services should exist
	list, err := services.List(ctx)
	if err != nil {
		ctx.PublishError("failed to list services: %s", err)
	}

	// start each service with the files it needs
	for _, service := range list {
		err := services.Run(ctx, service)
		if err != nil {
			ctx.PublishError("failed to start service %s: %s", service.Name, err)
		}
	}
}

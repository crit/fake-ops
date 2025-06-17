package services

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/crit/fake-ops/internal/app"
)

func StartApp(svc Service, ctx *app.Context) {
	if svc.Skip {
		ctx.PublishInfo("skipping %s", svc.Name)
		return
	}
	
	ctx.PublishInfo("starting service %s:%d", svc.Name, svc.Port)
	svc.Exec = strings.ReplaceAll(svc.Exec, "{port}", fmt.Sprintf("%d", svc.Port))

	parts := strings.Split(svc.Exec, " ")
	e := exec.Command(parts[0], parts[1:]...)

	if svc.Stdout {
		pipe, _ := e.StdoutPipe()
		go capture(pipe, svc, ctx.PublishInfo)
	}

	if svc.Stderr {
		pipe, _ := e.StderrPipe()
		go capture(pipe, svc, ctx.PublishError)
	}

	if err := e.Start(); err != nil {
		ctx.PublishServiceError(svc.Name)
		ctx.PublishError("error running %s: %s", svc.Name, err)
		return
	} else {
		ctx.PublishServiceOnline(svc.Name)
	}

	go func(ctx *app.Context) {
		<-ctx.Done()

		ctx.PublishServiceOffline(svc.Name)
		ctx.PublishInfo("stopping service %s", svc.Name)

		if err := e.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			ctx.PublishServiceError(svc.Name)
			ctx.PublishError("error killing process %s: %s", svc.Name, err)
		}
	}(ctx)
}

func capture(r io.ReadCloser, svc Service, publish func(msg string, args ...any)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		publish("%s: %s", svc.Name, line)
	}

	if err := scanner.Err(); err != nil {
		publish("%s capture err: %s", svc.Name, err)
	}
}

package services

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/crit/fake-ops/internal/app"
	"github.com/crit/fake-ops/internal/http_results"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

func StartHTTP(svc Service, ctx *app.Context) {
	if svc.Skip {
		ctx.PublishInfo("skipping %s", svc.Name)
		return
	}

	resultsPath := ctx.Flags.Results

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.PublishServiceError(svc.Name)
		ctx.PublishError("failed to create watcher for service %s: %s", svc.Name, err)
		return
	}
	defer watcher.Close()

	var mu sync.Mutex
	var server *http.Server

	// watch the directory for the service
	dirPath := filepath.Join(resultsPath, svc.Name)
	err = watcher.Add(dirPath)
	if err != nil {
		ctx.PublishServiceError(svc.Name)
		ctx.PublishError("failed to watch directory for service %s: %s", svc.Name, err)
		return
	}

	// Watch all existing files in the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		ctx.PublishServiceError(svc.Name)
		ctx.PublishError("failed to read directory %s: %s", dirPath, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		svc.Files = append(svc.Files, filePath)

		// Add the file to the watcher
		err = watcher.Add(filePath)
		if err != nil {
			ctx.PublishServiceError(svc.Name)
			ctx.PublishError("failed to watch file %s: %s", filePath, err)
		}
	}

	// Parse the initial responses based on current files
	parseResponses := func() {
		svc.Responses = nil // Clear previous responses

		for _, file := range svc.Files {
			data, err := os.ReadFile(file)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}

				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("failed to read file %s: %s", file, err)
				continue
			}

			result, err := http_results.Parse(data)
			if err != nil {
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("failed to parse file %s: %s", file, err)
				continue
			}

			svc.Responses = append(svc.Responses, result)
		}
	}

	parseResponses()

	// Function to start the server
	startServer := func() {
		mu.Lock()
		defer mu.Unlock()

		// Create a new Gin instance
		g := gin.New()
		g.GET("/", func(c *gin.Context) { c.String(http.StatusOK, svc.Name) })

		for _, result := range svc.Responses {
			handler := func(c *gin.Context) {
				c.Data(result.Code, result.ContentType, http_results.FillUUID(result.Data, len(c.Params)))
			}

			// Check if the route already exists
			existingRoutes := g.Routes()
			routeExists := false
			for _, r := range existingRoutes {
				if r.Method == result.Method && r.Path == result.Path {
					routeExists = true
					break
				}
			}

			if routeExists {
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("route already exists: %s %s", result.Method, result.Path)
				continue
			}

			switch result.Method {
			case "GET":
				g.GET(result.Path, handler)
			case "POST":
				g.POST(result.Path, handler)
			case "DELETE":
				g.DELETE(result.Path, handler)
			case "PUT":
				g.PUT(result.Path, handler)
			default:
				ctx.PublishServiceError(svc.Name)
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("%s unsupported method: %s", result.Path, result.Method)
			}
		}

		// Start HTTP server in a goroutine
		go func() {
			ctx.PublishServiceOnline(svc.Name)
			ctx.PublishInfo("starting service %s:%d", svc.Name, svc.Port)

			server = &http.Server{
				Addr:    ":" + strconv.Itoa(svc.Port),
				Handler: g,
			}

			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("server error: %s", err)
			}
		}()
	}

	// Function to gracefully stop the server
	stopCurrentServer := func() {
		mu.Lock()
		defer mu.Unlock()
		if server != nil {
			ctx.PublishInfo("stopping service %s", svc.Name)

			if err := server.Close(); err != nil {
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("error stopping server: %s", err)
			} else {
				ctx.PublishServiceOffline(svc.Name)
			}

			server = nil
		}
	}

	// Start the server initially
	startServer()

	// Watch for file changes and new files
	go func() {
		var timer *time.Timer // debounce timer

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Handle new files and changes
				if event.Has(fsnotify.Create | fsnotify.Write | fsnotify.Rename) {
					ctx.PublishInfo("%s: %s", event.Op, event.Name)

					if timer != nil {
						timer.Stop()
					}

					timer = time.AfterFunc(300*time.Millisecond, func() {
						// Handle directory changes (e.g., a new file added)
						if event.Op&fsnotify.Create != 0 {
							fileInfo, err := os.Stat(event.Name)
							if err == nil && !fileInfo.IsDir() {
								// Add the new file to svc.Files and watcher
								mu.Lock()
								svc.Files = append(svc.Files, event.Name)
								mu.Unlock()

								err = watcher.Add(event.Name)
								if err != nil {
									ctx.PublishServiceError(svc.Name)
									ctx.PublishError("failed to watch new file %s: %s", event.Name, err)
								}
							}
						}

						// Stop the current server and reload responses
						stopCurrentServer()
						parseResponses()
						startServer()
					})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				ctx.PublishServiceError(svc.Name)
				ctx.PublishError("file watcher error: %s", err)
			}
		}
	}()

	// wait for termination
	<-ctx.Done()
	stopCurrentServer()
}

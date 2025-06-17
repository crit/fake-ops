package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crit/fake-ops/internal/app"
	"github.com/crit/fake-ops/internal/http_results"
	"gopkg.in/yaml.v3"
)

type Type string

const (
	ServiceHTTP Type = "http"
	ServiceApp  Type = "app"
)

type Service struct {
	Skip   bool   `yaml:"skip"`
	Name   string `yaml:"name"`
	Port   int    `yaml:"port"`
	Type   Type   `yaml:"type"`
	Exec   string `yaml:"exec"`
	Stdout bool   `yaml:"stdout"`
	Stderr bool   `yaml:"stderr"`

	Files     []string
	Responses []*http_results.Result
}

func NewService(data []byte) (*Service, error) {
	var service Service

	err := yaml.Unmarshal(data, &service)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service: %s", err)
	}

	return &service, nil
}

func Run(ctx *app.Context, service Service) error {
	switch service.Type {
	case ServiceHTTP:
		ctx.PublishService("http", service.Name, service.Port)
		go StartHTTP(service, ctx)
	case ServiceApp:
		ctx.PublishService("app", service.Name, service.Port)
		go StartApp(service, ctx)
	default:
		return fmt.Errorf("unsupported service type: %s\n", service.Type)
	}

	return nil
}

func List(ctx *app.Context) ([]Service, error) {
	path := ctx.Flags.Services

	// get all files in dir path
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var list []Service

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		svc, err := NewService(data)
		if err != nil {
			return nil, err
		}

		list = append(list, *svc)
	}

	return list, nil
}

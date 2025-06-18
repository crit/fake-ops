package app

import (
	"flag"
	"sync"
)

// Flags contains flagged arguments to this program.
type Flags struct {
	Services string
	Results  string
}

// Parse handles loading flags passed to this program on startup.
func (f *Flags) Parse() {
	var once sync.Once

	once.Do(func() {
		svc := flag.String("services", "./services", "directory containing service definitions")
		res := flag.String("results", "./results", "directory containing service results")
		flag.Parse()

		f.Services = *svc
		f.Results = *res
	})
}

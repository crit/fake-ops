package app

import "flag"

type Flags struct {
	Services string `flag:"services" default:"./services" validate:"required"`
	Results  string `flag:"results" default:"./results" validate:"required"`
}

func (f *Flags) Parse() {
	svc := flag.String("services", "./services", "directory containing service definitions")
	res := flag.String("results", "./results", "directory containing service results")
	flag.Parse()

	f.Services = *svc
	f.Results = *res
}

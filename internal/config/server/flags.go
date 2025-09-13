package server

import "flag"

type flagValues struct {
	config  string
	address string
	verbose bool
}

func readFlags(values flagValues) flagValues {
	flag.StringVar(&values.config, "config", values.config, "path to config file")
	flag.StringVar(&values.address, "address", values.address, "server address")
	flag.BoolVar(&values.verbose, "verbose", values.verbose, "enable verbose logging")

	flag.Parse()

	return values
}

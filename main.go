package main

import (
	"github.com/codegangsta/cli"
	"github.com/euclid1990/go-bigquery/cmd"
	"github.com/euclid1990/go-bigquery/utilities"
	"os"
)

func main() {
	// Read env vars from .env file
	utilities.LoadEnv("")
	// Create cli application
	app := cli.NewApp()
	app.Name = "Go-BigQuery"
	app.Version = "1.0.0"
	app.Usage = "A small cli written in Go to queries, loading data, and exporting data in BigQuery."
	app.Action = cmd.Action
	app.Flags = cmd.Flags
	app.Run(os.Args)
}

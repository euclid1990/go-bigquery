package main

import (
    "github.com/codegangsta/cli"
    "github.com/euclid1990/go-bigquery/cmd"
    "os"
)

func main() {
    app := cli.NewApp()
    app.Name = "Go-BigQuery"
    app.Version = "1.0.0"
    app.Usage = "A small cli written in Go to queries, loading data, and exporting data in BigQuery."
    app.Action = cmd.Action
    app.Flags = cmd.Flags
    app.Run(os.Args)
}

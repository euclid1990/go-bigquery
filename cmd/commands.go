package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/euclid1990/go-bigquery/configs"
	"github.com/euclid1990/go-bigquery/utilities"
)

// List of options
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  "exec",
		Value: "all",
		Usage: "Execute action you want to do",
	},
}

// Action defines the main action for application
func Action(c *cli.Context) {
	exec := c.String("exec")
	utilities.Log(configs.LOG_INFO, fmt.Sprintf("Action: %v", exec))

	switch exec {
	case configs.ACTION_ALL:
		fmt.Printf("Run [All] command.\n")
	case configs.ACTION_INIT:
		fmt.Printf("Run [Init] command.\n")
	}
}

package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/euclid1990/go-bigquery/configs"
	s "github.com/euclid1990/go-bigquery/schemas"
	"github.com/euclid1990/go-bigquery/utilities"
	"github.com/euclid1990/go-bigquery/web"
	"golang.org/x/net/context"
	"os"
)

// List of options
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  "exec",
		Value: "all",
		Usage: "Execute action you want to do",
	},
}

// Instance of Google BigQuery
var (
	bigquery *utilities.BigQuery
	ctx      = context.Background()
)

// Action defines the main action for application
func Action(c *cli.Context) {
	exec := c.String("exec")
	utilities.Log(configs.LOG_INFO, fmt.Sprintf("Action: %v", exec))

	switch exec {
	case configs.ACTION_ALL:
		fmt.Printf("Run [All] command.\n")
		web.Init()

	case configs.ACTION_INIT:
		fmt.Printf("Run [Init] command.\n")
		fmt.Println("-------------------------------------")

		bigquery := utilities.NewBigQuery(ctx)
		datasetId := os.Getenv("DATASET_ID")
		bigquery.CreateDataset(ctx, datasetId)

		// Create & Load data into 'users' Table
		if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_USER, s.User{}); err == nil {
			userJsonFilePath := utilities.GetUserJsonFilePath()
			bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_USER, userJsonFilePath)
		}

		// Create & Load data into 'access' Table
		if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_ACCESS, s.Access{}); err == nil {
			accessJsonFilePath := utilities.GetAccessJsonFilePath()
			bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_ACCESS, accessJsonFilePath)
		}

		// List all tables
		fmt.Println("-------------------------------------")
		fmt.Printf("Start listing all tables in BigQuery:\n")
		bigquery.ListTable(ctx, datasetId)

	case configs.ACTION_FAKE:
		fmt.Printf("Run [Fake] command.\n")
		utilities.GenrateDummyData(configs.DATA_TYPE_JSON)

	case configs.ACTION_WEB:
		fmt.Printf("Run [Web] command.\n")
		web.Init()
	}
}

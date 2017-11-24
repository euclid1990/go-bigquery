package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/euclid1990/go-bigquery/configs"
	s "github.com/euclid1990/go-bigquery/schemas"
	"github.com/euclid1990/go-bigquery/utilities"
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
	case configs.ACTION_INIT:
		fmt.Printf("Run [Init] command.\n")
		bigquery := utilities.NewBigQuery(ctx)
		datasetId := os.Getenv("DATASET_ID")
		bigquery.CreateDataset(ctx, datasetId)
		bigquery.CreateTable(ctx, datasetId, s.TABLE_USER, s.User{})
		userJsonFilePath := fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_USER, configs.DATA_TYPE_JSON)
		bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_USER, userJsonFilePath)
		accessJsonFilePath := fmt.Sprintf(configs.DATA_FORMAT_FILE_NAME, configs.DATA_INPUT_PATH+configs.DATA_INPUT_ACCESSS, configs.DATA_TYPE_JSON)
		bigquery.CreateTable(ctx, datasetId, s.TABLE_ACCESS, s.Access{})
	case configs.ACTION_FAKE:
		fmt.Printf("Run [Fake] command.\n")
		utilities.GenrateDummyData(configs.DATA_TYPE_JSON)
	}
}

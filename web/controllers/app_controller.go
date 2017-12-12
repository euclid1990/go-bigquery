package controllers

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/euclid1990/go-bigquery/configs"
	s "github.com/euclid1990/go-bigquery/schemas"
	"github.com/euclid1990/go-bigquery/utilities"
	"github.com/kataras/iris/mvc"
	"golang.org/x/net/context"
	"os"
)

type AppController struct {
	mvc.C
}

var (
	ctx      = context.Background()
	bigquery *utilities.BigQuery
)

func (c *AppController) Get() mvc.View {
	return mvc.View{
		Name: "app/index.html",
		Data: map[string]interface{}{
			"Title": "BigQuery & Golang App",
		},
	}
}

func (c *AppController) PostFakeData() {
	total := utilities.GenrateDummyData(configs.DATA_TYPE_JSON)
	storage, _ := utilities.DirSize(configs.DATA_INPUT_PATH)
	msg := fmt.Sprintf("%s record is generated. Total: %s (MB).", humanize.Comma(total), humanize.FormatFloat("#,###.##", storage))
	c.Ctx.JSON(map[string]string{"msg": msg})
}

func (c *AppController) PostSeedTable() {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")

	// Create BigQuery Database
	bigquery.CreateDataset(ctx, datasetId)

	// Create & Load data into 'users' Table
	if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_USER, s.User{}); err != nil {
		userJsonFilePath := utilities.GetUserJsonFilePath()
		bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_USER, userJsonFilePath)
	}

	// Create & Load data into 'access' Table
	if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_ACCESS, s.Access{}); err != nil {
		accessJsonFilePath := utilities.GetAccessJsonFilePath()
		bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_ACCESS, accessJsonFilePath)
	}

	msg := fmt.Sprintf("Inserted data into table: '%s', '%s' on BigQuery.", s.TABLE_USER, s.TABLE_ACCESS)
	c.Ctx.JSON(map[string]string{"msg": msg})
}

func (c *AppController) PostDropTable() {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")

	// Create BigQuery Database
	bigquery.CreateDataset(ctx, datasetId)

	// Drop 'users' & 'access' table
	bigquery.DeleteTable(ctx, datasetId, s.TABLE_USER)
	bigquery.DeleteTable(ctx, datasetId, s.TABLE_ACCESS)

	msg := fmt.Sprintf("Table: '%s', '%s' have been dropped.", s.TABLE_USER, s.TABLE_ACCESS)
	c.Ctx.JSON(map[string]string{"msg": msg})
}

func (c *AppController) GetUser() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "User",
		},
	}
}

func (c *AppController) GetAccess() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "Access",
		},
	}
}

func (c *AppController) GetRetention() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "RetentionRate",
		},
	}
}

func (c *AppController) GetDevice() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "Device",
		},
	}
}

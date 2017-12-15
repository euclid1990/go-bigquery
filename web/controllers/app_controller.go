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
	"reflect"
	"strings"
	"time"
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
	if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_USER, s.User{}); err == nil {
		userJsonFilePath := utilities.GetUserJsonFilePath()
		bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_USER, userJsonFilePath)
	}

	// Create & Load data into 'access' Table
	if err := bigquery.CreateTable(ctx, datasetId, s.TABLE_ACCESS, s.Access{}); err == nil {
		accessJsonFilePath := utilities.GetAccessJsonFilePath()
		bigquery.InsertDataFromFile(ctx, datasetId, s.TABLE_ACCESS, accessJsonFilePath)
	}

	msg := fmt.Sprintf("Inserted data into table: '%s', '%s' on BigQuery.", s.TABLE_USER, s.TABLE_ACCESS)
	c.Ctx.JSON(map[string]string{"msg": msg})
}

func (c *AppController) PostDropTable() {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")

	// Drop 'users' & 'access' table
	bigquery.DeleteTable(ctx, datasetId, s.TABLE_USER)
	bigquery.DeleteTable(ctx, datasetId, s.TABLE_ACCESS)

	msg := fmt.Sprintf("Table: '%s', '%s' have been dropped.", s.TABLE_USER, s.TABLE_ACCESS)
	c.Ctx.JSON(map[string]string{"msg": msg})
}

func (c *AppController) GetUser() mvc.View {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")
	tableId := s.TABLE_USER
	q := fmt.Sprintf(`
        SELECT COUNT(id) total_user, FORMAT_DATETIME("%%Y", created_at) created_year, FORMAT_DATETIME("%%m", created_at) created_month, FORMAT_DATETIME("%%b", created_at) month_name
        FROM %s.%s
        GROUP BY created_year, created_month, month_name
        ORDER BY created_year, created_month ASC
    `, datasetId, tableId)
	r, _ := bigquery.Query(ctx, datasetId, q)

	// Create xAxis & series data for Highchart
	xCategories := make([]string, 0)
	series := make([]int64, 0)
	for i := 0; i < len(r); i++ {
		item := reflect.ValueOf(r[i])
		xCategories = append(xCategories, item.Index(1).Interface().(string)+" "+item.Index(3).Interface().(string))
		series = append(series, item.Index(0).Interface().(int64))
	}

	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":          "User Analytics",
			"Header":         "User",
			"NewUserByMonth": map[string]interface{}{"xCategories": xCategories, "series": series},
		},
	}
}

func (c *AppController) GetAccess() mvc.View {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")
	accessTableId := s.TABLE_ACCESS
	userTableId := s.TABLE_USER
	q := fmt.Sprintf(`
        WITH max_total AS (SELECT COUNT(access.id) total_access FROM %s.%s)
        SELECT users.address.country country, COUNT(access.id)*100/total_access access_by_country, total_access
        FROM %s.%s access, max_total
        JOIN %s.%s users ON access.user_id = users.id
        GROUP BY country, total_access
        ORDER BY country ASC
    `, datasetId, accessTableId, datasetId, accessTableId, datasetId, userTableId)
	r, _ := bigquery.Query(ctx, datasetId, q)

	// Retrieve total access of all users
	var total_access int64 = 0
	for i := 0; i < len(r); i++ {
		item := reflect.ValueOf(r[i])
		total_access = item.Index(2).Interface().(int64)
		break
	}

	// Calculate average access percent by country
	average_access := float64(total_access) / float64(len(r))
	average_access = average_access * 100 / float64(total_access)
	var other float64 = 0.0

	// Create series data for Highchart
	series := make([]map[string]interface{}, 0)
	for i := 0; i < len(r); i++ {
		item := reflect.ValueOf(r[i])
		country := item.Index(0).Interface().(string)
		access_by_country := item.Index(1).Interface().(float64)
		if access_by_country >= average_access {
			series = append(series, map[string]interface{}{"name": country, "y": access_by_country})
		} else {
			other = other + access_by_country
		}
	}

	// Generate other's total access
	series = append(series, map[string]interface{}{"name": "Other", "y": other})
	return mvc.View{
		Name: "app/access.html",
		Data: map[string]interface{}{
			"Title":           "Access Analytics",
			"Header":          "Access",
			"AccessByCountry": map[string]interface{}{"series": series},
		},
	}
}

func (c *AppController) GetRetention() mvc.View {
	bigquery := utilities.NewBigQuery(ctx)
	datasetId := os.Getenv("DATASET_ID")
	accessTableId := s.TABLE_ACCESS

	now := time.Now()
	oneMonthAgo := make([]string, 0)
	for i := 30; i >= 1; i-- {
		day := now.AddDate(0, 0, -1*i)
		oneMonthAgo = append(oneMonthAgo[:], day.Format("2006-01-02"))
	}
	oneMonthAgoStr := "\"" + strings.Join(oneMonthAgo, "\",\"") + "\""

	q1 := fmt.Sprintf(`
        WITH OneMonthAgo AS (
            SELECT * FROM UNNEST([%s]) AS each_day
        )

        SELECT COUNT(DISTINCT user_id) total_access_user, each_day
        FROM OneMonthAgo
        LEFT JOIN %s.%s ON FORMAT_DATETIME("%%Y-%%m-%%d", DATETIME_ADD(access_at, INTERVAL 30 DAY)) = OneMonthAgo.each_day
        GROUP BY each_day
        ORDER BY each_day ASC
    `, oneMonthAgoStr, datasetId, accessTableId)
	r1, _ := bigquery.Query(ctx, datasetId, q1)

	q2 := fmt.Sprintf(`
        WITH OneMonthAgo AS (
            SELECT * FROM UNNEST([%s]) AS each_day
        ),
        SubQ1 AS (
            SELECT DISTINCT user_id, each_day
            FROM OneMonthAgo
            LEFT JOIN %s.%s ON FORMAT_DATETIME("%%Y-%%m-%%d", access_at) = OneMonthAgo.each_day
            GROUP BY each_day, user_id
        ),
        SubQ2 AS (
            SELECT DISTINCT user_id, each_day
            FROM OneMonthAgo
            LEFT JOIN %s.%s ON FORMAT_DATETIME("%%Y-%%m-%%d", DATETIME_ADD(access_at, INTERVAL 30 DAY)) = OneMonthAgo.each_day
            GROUP BY each_day, user_id
        ),
        SubQ3 AS (
            SELECT DISTINCT SubQ1.user_id, SubQ1.each_day
            FROM SubQ1
            JOIN SubQ2 ON SubQ1.each_day = SubQ2.each_day
            WHERE SubQ1.user_id = SubQ2.user_id
        )

        SELECT COUNT(SubQ3.user_id) retention_user, OneMonthAgo.each_day
        FROM OneMonthAgo
        LEFT JOIN SubQ3 ON SubQ3.each_day = OneMonthAgo.each_day
        GROUP BY OneMonthAgo.each_day
        ORDER BY OneMonthAgo.each_day ASC
    `, oneMonthAgoStr, datasetId, accessTableId, datasetId, accessTableId)
	r2, _ := bigquery.Query(ctx, datasetId, q2)

	// Create series data for Highchart
	series := make(map[string][]interface{}, 0)
	series["days"] = make([]interface{}, 0)
	series["past30days"] = make([]interface{}, 0)
	series["rate"] = make([]interface{}, 0)
	for i := 0; i < len(r1); i++ {
		item1 := reflect.ValueOf(r1[i])
		item2 := reflect.ValueOf(r2[i])
		past30dayCount := item1.Index(0).Interface().(int64)
		dayCount := item2.Index(0).Interface().(int64)
		series["past30days"] = append(series["past30days"], past30dayCount)
		series["days"] = append(series["days"], dayCount)
		series["rate"] = append(series["rate"], float64(dayCount*100)/float64(past30dayCount))
	}

	return mvc.View{
		Name: "app/retention.html",
		Data: map[string]interface{}{
			"Title":     "Retention Rate",
			"Header":    "Retention Rate",
			"Retention": map[string]interface{}{"xCategories": oneMonthAgo, "series": series},
		},
	}
}

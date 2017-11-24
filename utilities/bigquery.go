package utilities

import (
	"cloud.google.com/go/bigquery"
	"github.com/euclid1990/go-bigquery/configs"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"os"
)

type BigQuery struct {
	ProjectId string             `json:"project_id"`
	Client    *(bigquery.Client) `json:"client"`
}

func NewBigQuery(ctx context.Context) *BigQuery {
	projectID := os.Getenv("GCP_PROJECT_ID")
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to create client: %v", err)
	}
	Log(configs.LOG_INFO, "BigQuery client created")
	return &BigQuery{ProjectId: projectID, Client: client}
}

func (b *BigQuery) CreateDataset(ctx context.Context, datasetId string) error {
	if _, err := b.Client.Dataset(datasetId).Metadata(ctx); err != nil {
		if err.(*googleapi.Error).Code == 404 {
			if err := b.Client.Dataset(datasetId).Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
				Logf(configs.LOG_ERROR, "Failed to create dataset: %v", err)
				return err
			}
			Log(configs.LOG_INFO, "Dataset client created")
			return nil
		} else {
			Logf(configs.LOG_ERROR, "Failed to fetch metadata: %v", err)
		}
	}
	return nil
}

func (b *BigQuery) CreateTable(ctx context.Context, datasetId string, tableId string, schemaStruct interface{}) error {
	schema, err := bigquery.InferSchema(schemaStruct)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to fetch metadata: %v", err)
		return err
	}
	table := b.Client.Dataset(datasetId).Table(tableId)
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		Logf(configs.LOG_ERROR, "Failed to create table %s: %v", tableId, err)
		return err
	}
	Logf(configs.LOG_INFO, "Dataset: %s - Table: %s is created", datasetId, tableId)
	return nil
}

package utilities

import (
	"cloud.google.com/go/bigquery"
	"github.com/euclid1990/go-bigquery/configs"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"os"
	"reflect"
	"time"
)

type BigQuery struct {
	ProjectId string             `json:"project_id"`
	Client    *(bigquery.Client) `json:"client"`
}

func NewBigQuery(ctx context.Context) *BigQuery {
	projectId := os.Getenv("GCP_PROJECT_ID")
	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to create client: %v", err)
	}
	Log(configs.LOG_INFO, "BigQuery client created")
	return &BigQuery{ProjectId: projectId, Client: client}
}

func (b *BigQuery) CreateDataset(ctx context.Context, datasetId string) error {
	if _, err := b.Client.Dataset(datasetId).Metadata(ctx); err != nil {
		if err.(*googleapi.Error).Code == 404 {
			if err := b.Client.Dataset(datasetId).Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
				Logf(configs.LOG_ERROR, "Failed to create dataset: %v", err)
				return err
			}
			Logf(configs.LOG_INFO, "Dataset client: '%s' created", datasetId)
			return nil
		} else {
			Logf(configs.LOG_ERROR, "Failed to fetch metadata: %v", err)
		}
	}
	Logf(configs.LOG_INFO, "Dataset: '%s' has already existed", datasetId)
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
		if err.(*googleapi.Error).Code == 409 {
			Logf(configs.LOG_INFO, "Table: '%s' has already existed", tableId)
			return err
		}
		Logf(configs.LOG_ERROR, "Failed to create table %s: %v", tableId, err)
		return err
	}
	Logf(configs.LOG_INFO, "Table: '%s' created", tableId)
	return nil
}

func (b *BigQuery) CreateView(ctx context.Context, datasetId string, viewId string, query string, expirationTime time.Time) error {
	table := b.Client.Dataset(datasetId).Table(viewId)
	// The time when this table expires. If not set, the table will persist indefinitely.
	if expirationTime.IsZero() {
		expirationTime = time.Now().Add(1 * time.Hour)
	}
	if err := table.Create(ctx, &bigquery.TableMetadata{
		Schema:         nil,
		ViewQuery:      query,
		ExpirationTime: expirationTime,
	}); err != nil {
		if err.(*googleapi.Error).Code == 409 {
			Logf(configs.LOG_INFO, "View: '%s' has already existed", viewId)
			return nil
		}
		Logf(configs.LOG_ERROR, "Failed to create view %s: %v", viewId, err)
		return err
	}
	Logf(configs.LOG_INFO, "View: '%s' created", viewId)
	return nil
}

func (b *BigQuery) InsertDataFromFile(ctx context.Context, datasetId, tableId, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	source := bigquery.NewReaderSource(f)
	// Config CSV options
	// source.AllowJaggedRows = true
	// source.SkipLeadingRows = 1
	// Using json/csv format: bigquery.JSON/bigquery.CSV
	source.SourceFormat = bigquery.JSON

	loader := b.Client.Dataset(datasetId).Table(tableId).LoaderFrom(source)
	loader.CreateDisposition = bigquery.CreateNever

	job, err := loader.Run(ctx)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to create job of table '%s': %v", tableId, err)
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to get job status of table '%s': %v", tableId, err)
		return err
	}
	if err := status.Err(); err != nil {
		Logf(configs.LOG_ERROR, "Failed to inserted data into table '%s': %v", tableId, err)
		return err
	}
	Logf(configs.LOG_INFO, "Inserted data into table: '%s'", tableId)
	return nil
}

func (b *BigQuery) ListTable(ctx context.Context, datasetId string) error {
	tables := b.Client.Dataset(datasetId).Tables(ctx)
	for {
		t, err := tables.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		Logf(configs.LOG_INFO, "Table: %s", t.TableID)
	}
	return nil
}

func (b *BigQuery) CopyTable(ctx context.Context, datasetId, srcId, dstId string) error {
	dataset := b.Client.Dataset(datasetId)
	copier := dataset.Table(dstId).CopierFrom(dataset.Table(srcId))
	copier.WriteDisposition = bigquery.WriteTruncate
	job, err := copier.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	Logf(configs.LOG_INFO, "Table: '%s' has been copied into Table: '%s'", srcId, dstId)
	return nil
}

func (b *BigQuery) DeleteTable(ctx context.Context, datasetId, tableId string) error {
	dataset := b.Client.Dataset(datasetId)
	if err := dataset.Table(tableId).Delete(ctx); err != nil {
		if err.(*googleapi.Error).Code == 404 {
			Logf(configs.LOG_INFO, "Table: '%s' is not existing", tableId)
			return err
		}
		Logf(configs.LOG_ERROR, "Failed to deleted table '%s': %v", tableId, err)
		return err
	}
	Logf(configs.LOG_INFO, "Table: '%s' has been deleted", tableId)
	return nil
}

// Type of insert data must be slice/array []*Struct
func (b *BigQuery) InsertRow(ctx context.Context, datasetId, tableId string, data interface{}) error {
	u := b.Client.Dataset(datasetId).Table(tableId).Uploader()
	if err := u.Put(ctx, data); err != nil {
		Logf(configs.LOG_ERROR, "Failed to insert data into table '%s': %v", tableId, err)
		return err
	}
	srcVal := reflect.ValueOf(data)
	Logf(configs.LOG_INFO, "Total %d record has been inserted into Table: '%s'", srcVal.Len(), tableId)
	return nil
}

func (b *BigQuery) Query(ctx context.Context, datasetId, sqlQuery string) ([]interface{}, error) {
	start := time.Now()
	q := b.Client.Query(sqlQuery)
	// Use standard SQL syntax for queries.
	// See: https://cloud.google.com/bigquery/sql-reference/
	q.QueryConfig.UseStandardSQL = true
	q.Dst = nil

	it, err := q.Read(ctx)
	if err != nil {
		Logf(configs.LOG_ERROR, "Failed to query ----------\n%s\n----------%v", sqlQuery, err)
		return nil, err
	}
	result := make([]interface{}, 0)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			Logf(configs.LOG_ERROR, "Failed to map result data from query ----------\n%s\n----------%v", sqlQuery, err)
			return nil, err
		}
		result = append(result, row)
	}
	Logf(configs.LOG_INFO, "Get total %v records, Query took %.3f s", len(result), TimeTrack(start))
	return result, nil
}

package bq

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
)

type Client struct {
	bqClient *bigquery.Client
}

// Query implements interfaces.BigQueryClient.
func (c *Client) Query(ctx context.Context, query string) (interfaces.BigQueryIterator, error) {
	q := c.bqClient.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return it, nil
}

func New(projectID string) (*Client, error) {
	ctx := context.Background()
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}
	return &Client{bqClient: bqClient}, nil
}

var _ interfaces.BigQueryClient = (*Client)(nil)

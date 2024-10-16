package cs

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/secmon-as-code/overseer/pkg/interfaces"
)

type client struct {
	storageClient *storage.Client
}

// NewClient creates a new Cloud Storage client.
func NewClient(ctx context.Context, options ...option.ClientOption) (*client, error) {
	storageClient, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	return &client{storageClient: storageClient}, nil
}

// GetObject implements interfaces.CloudStorageClient.
func (c *client) GetObject(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	bucket := c.storageClient.Bucket(bucketName)
	object := bucket.Object(objectName)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// PutObject implements interfaces.CloudStorageClient.
func (c *client) PutObject(ctx context.Context, bucketName string, objectName string) (io.WriteCloser, error) {
	bucket := c.storageClient.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	writer.ContentType = "application/octet-stream"
	writer.ChunkSize = 0 // Use default chunk size
	writer.Metadata = map[string]string{
		"created": time.Now().Format(time.RFC3339),
	}
	return writer, nil
}

var _ interfaces.CloudStorageClient = &client{}

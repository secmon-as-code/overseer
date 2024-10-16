package interfaces

import (
	"context"
	"io"

	"github.com/m-mizutani/opac"
	"github.com/open-policy-agent/opa/ast"
)

type CloudStorageClient interface {
	PutObject(ctx context.Context, bucketName, objectName string) (io.WriteCloser, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
}

type BigQueryClient interface {
	Query(ctx context.Context, query string) (BigQueryIterator, error)
}

type BigQueryIterator interface {
	Next(dst interface{}) error
}

type PubSubClient interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type PolicyClient interface {
	Query(ctx context.Context, query string, input, output any, options ...opac.QueryOption) error
	AnnotationSet() *ast.AnnotationSet
}

package usecase_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-as-code/overseer/pkg/adaptor"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
	"github.com/secmon-as-code/overseer/pkg/mock"
	"github.com/secmon-as-code/overseer/pkg/usecase"
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

type iterator struct {
	results []map[string]any
	index   int
}

func (it *iterator) Next(row interface{}) error {
	if it.index >= len(it.results) {
		return io.EOF
	}

	data := it.results[it.index]
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawData, row); err != nil {
		return err
	}

	it.index++
	return nil
}

//go:embed testdata/query1.sql
var query1 []byte

func TestExtract(t *testing.T) {
	mockBQ := mock.BigQueryClientMock{
		QueryFunc: func(ctx context.Context, query string) (interfaces.BigQueryIterator, error) {
			return &iterator{
				results: []map[string]any{
					{"key1": "value1"},
					{"key1": "value2"},
				},
			}, nil
		},
	}

	var buf buffer
	cache := &mock.CacheServiceMock{
		NewWriterFunc: func(ctx context.Context, ID model.QueryID) (io.WriteCloser, error) {
			return &buf, nil
		},
	}

	uc := usecase.New(adaptor.New(
		adaptor.WithBigQuery(&mockBQ),
	))

	gt.NoError(t, uc.Extract(context.Background(), model.Queries{
		model.MustNewQuery(query1),
	}, cache))

	var result []map[string]string
	gt.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	gt.Equal(t, len(result), 2)
	gt.Equal(t, result[0]["key1"], "value1")
	gt.Equal(t, result[1]["key1"], "value2")
}

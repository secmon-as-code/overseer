package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
)

type NewWriter func(ID model.QueryID) (io.WriteCloser, error)

func (x *UseCase) Extract(ctx context.Context, queries model.Queries, cache interfaces.CacheService) error {
	if err := queries.Validate(); err != nil {
		return err
	}

	for _, query := range queries {
		if err := queryAndDump(ctx, x.clients.BigQuery(), query, cache); err != nil {
			return goerr.Wrap(err, "fail to extract query")
		}
	}

	return nil
}

func queryAndDump(ctx context.Context, bq interfaces.BigQueryClient, query *model.Query, cache interfaces.CacheService) error {
	it, err := bq.Query(ctx, query.String())
	if err != nil {
		return err
	}

	w, err := cache.NewWriter(ctx, query.ID())
	if err != nil {
		return err
	}
	defer safeClose(ctx, w)

	if _, err := w.Write([]byte("[")); err != nil {
		return goerr.Wrap(err, "fail to write header bracket")
	}

	isFirst := true
	for {
		var row any
		if err := it.Next(&row); err == io.EOF {
			break
		} else if err != nil {
			return goerr.Wrap(err, "fail to get next row")
		}

		if !isFirst {
			if _, err := w.Write([]byte(",")); err != nil {
				return goerr.Wrap(err, "fail to write separator")
			}
		}
		isFirst = false

		data, err := json.Marshal(row)
		if err != nil {
			return goerr.Wrap(err, "fail to marshal row")
		}

		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("]")); err != nil {
		return goerr.Wrap(err, "fail to write footer bracket")
	}

	return nil
}

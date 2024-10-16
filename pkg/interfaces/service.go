package interfaces

import (
	"context"
	"io"

	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

type CacheService interface {
	NewWriter(ctx context.Context, ID model.QueryID) (io.WriteCloser, error)
	NewReader(ctx context.Context, ID model.QueryID) (io.ReadCloser, error)
}

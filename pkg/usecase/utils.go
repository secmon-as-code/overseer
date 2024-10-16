package usecase

import (
	"context"
	"io"

	"github.com/secmon-as-code/overseer/pkg/logging"
)

func safeClose(ctx context.Context, c io.Closer) {
	if err := c.Close(); err != nil {
		logging.FromCtx(ctx).Error("fail to close", "error", err)
	}
}

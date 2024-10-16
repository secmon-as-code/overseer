package model

import (
	"github.com/google/uuid"
	"github.com/secmon-as-code/overseer/pkg/logging"
)

type JobID string

func NewJobID() JobID {
	id, err := uuid.NewV7()
	if err != nil {
		logging.Default().Error("fail to generate new JobID", "err", err)
		panic(err)
	}
	return JobID(id.String())
}

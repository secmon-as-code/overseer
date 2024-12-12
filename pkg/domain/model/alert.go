package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
)

type AlertID string

func NewAlertID() AlertID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return AlertID(id.String())
}

const (
	AlertSchemaVersion = "v0"
)

// Alert is a notification of overseer. It's generated by eval rules and sent to NotifyService. It's based on AlertBody and has additional metadata fields.
type Alert struct {
	// ID is unique identifier of alert. It would be used for aggregation. If ID is empty, it will be generated automatically.
	ID AlertID `json:"id"`

	// Version is schema version of alert. It should be overwritten with AlertSchemaVersion.
	Version string `json:"version"`

	// JobID is identifier of job that generates alert.
	JobID JobID `json:"job_id"`

	// Timestamp is time when alert is generated. If Timestamp is zero, it will be set to current time.
	Timestamp time.Time `json:"timestamp"`

	AlertBody
}

// AlertBody is a body of Alert. It contains title, description, timestamp and additional attributes.
type AlertBody struct {
	// Title is short description of alert. It's required.
	Title string `json:"title"`

	// Description is detailed description of alert.
	Description string `json:"description"`

	// Timestamp is time when alert is generated. If Timestamp is nil, it will be set to current time. It allows not only string time format (RFC3339) but also integer (Unix timestamp) and float (Unix timestamp with nano seconds).
	Timestamp any `json:"timestamp"`

	// Attrs is key-value pairs of additional information of alert.
	Attrs Attrs `json:"attrs"`
}

type Attrs map[string]any

func NewAlert(ctx context.Context, body AlertBody) (*Alert, error) {
	var ts time.Time

	if body.Timestamp == nil {
		ts = time.Now()
	} else {
		switch v := body.Timestamp.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, goerr.Wrap(err, "fail to parse timestamp").With("timestamp", v)
			}
			ts = t

		case int:
			ts = time.Unix(int64(v), 0)

		case float64:
			sec := int64(v)
			nsec := int64((v - float64(sec)) * 1e9)
			ts = time.Unix(sec, nsec)

		default:
			return nil, goerr.New("unsupported timestamp type").With("timestamp", v)
		}
	}

	jobID := JobIDFromCtx(ctx)

	x := &Alert{
		ID:        NewAlertID(),
		Version:   AlertSchemaVersion,
		JobID:     jobID,
		Timestamp: ts,

		AlertBody: body,
	}

	if x.Timestamp.IsZero() {
		x.Timestamp = time.Now()
	}

	if err := body.Validate(); err != nil {
		return nil, err
	}

	return x, nil
}

func (x AlertBody) Validate() error {
	if x.Title == "" {
		return goerr.New("title is required")
	}

	return nil
}

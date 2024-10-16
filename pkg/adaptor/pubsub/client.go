package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
	"github.com/secmon-as-code/overseer/pkg/logging"
)

type Client struct {
	client *pubsub.Client
}

// Publish implements interfaces.PubSubClient.
func (c *Client) Publish(ctx context.Context, topic string, data []byte) error {
	t := c.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Block until the result is returned and a server-generated ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	logging.FromCtx(ctx).Debug("Published message", "id", id)
	return nil
}

var _ interfaces.PubSubClient = (*Client)(nil)

func New(projectID string) (*Client, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	return &Client{client: client}, nil
}

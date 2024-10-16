package adaptor

import "github.com/secmon-as-code/overseer/pkg/interfaces"

type Clients struct {
	cloudStorage interfaces.CloudStorageClient
	bigQuery     interfaces.BigQueryClient
	pubSub       interfaces.PubSubClient
	policy       interfaces.PolicyClient
}

func (x *Clients) CloudStorage() interfaces.CloudStorageClient {
	return x.cloudStorage
}

func (x *Clients) BigQuery() interfaces.BigQueryClient {
	return x.bigQuery
}

func (x *Clients) PubSub() interfaces.PubSubClient {
	return x.pubSub
}

func (x *Clients) Policy() interfaces.PolicyClient {
	return x.policy
}

func New(options ...Option) *Clients {
	c := &Clients{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

type Option func(*Clients)

func WithCloudStorage(client interfaces.CloudStorageClient) Option {
	return func(c *Clients) {
		c.cloudStorage = client
	}
}

func WithBigQuery(client interfaces.BigQueryClient) Option {
	return func(c *Clients) {
		c.bigQuery = client
	}
}

func WithPubSub(client interfaces.PubSubClient) Option {
	return func(c *Clients) {
		c.pubSub = client
	}
}

func WithPolicy(client interfaces.PolicyClient) Option {
	return func(c *Clients) {
		c.policy = client
	}
}

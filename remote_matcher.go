package synapse

import (
	"context"
	"time"
)

type RemoteMatcher struct {
	config ConfigMatcher
	client *Client
}

func newRemoteMatcher(config ConfigMatcher) (Matcher, error) {
	client, err := newClient(config.Host)
	if err != nil {
		return nil, err
	}

	return &RemoteMatcher{
		client: client,
		config: config,
	}, nil
}

func (m RemoteMatcher) Name() string {
	return m.config.Name
}

func (m RemoteMatcher) Match(params Context) (Hits, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return m.client.Match(ctx, m.Name(), params)
}

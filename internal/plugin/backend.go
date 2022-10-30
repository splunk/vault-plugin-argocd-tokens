package plugin

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend is the backend for the argo cd tokens plugin
type backend struct {
	*framework.Backend
	logger hclog.Logger
}

// Factory is the factory that produces the backend.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := getBackend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

// getBackend returns a configured backend
func getBackend(conf *logical.BackendConfig) *backend {
	backend := &backend{logger: conf.Logger}
	backend.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			pathProjectToken(backend),
			pathAccountToken(backend),
			pathConfig(backend),
		),
		Secrets: []*framework.Secret{
			secretProjectToken(backend),
			secretAccountToken(backend),
		},
		Help: trimHelp(helpBackend),
	}
	return backend
}

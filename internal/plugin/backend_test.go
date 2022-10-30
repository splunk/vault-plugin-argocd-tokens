package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestCreateBackend(t *testing.T) {
	getTestBackend(t) // ensure nothing panics or errors out
}

func getTestBackend(t *testing.T) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	require.NoError(t, err)
	return b.(*backend), config.StorageView
}

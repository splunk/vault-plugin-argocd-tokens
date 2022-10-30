package plugin

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
)

func saveToStorage[V any](ctx context.Context, storage logical.Storage, key string, data *V) error {
	entry, err := logical.StorageEntryJSON(key, data)
	if err != nil {
		return fmt.Errorf("error while creating the storage entry: %s", err)
	}

	if err := storage.Put(ctx, entry); err != nil {
		return fmt.Errorf("error while saving the storage entry: %s", err)
	}

	return nil
}

func readFromStorage[V any](ctx context.Context, storage logical.Storage, key string) (value V, err error) {
	rawValue, err := storage.Get(ctx, key)

	if rawValue == nil || err != nil {
		return value, fmt.Errorf("error while reading the storage entry: %s", err)
	}

	if err := rawValue.DecodeJSON(&value); err != nil {
		return value, fmt.Errorf("error while decoding the storage entry: %s", err)
	}

	return value, nil
}

// Ignores missing key in the storage, still returns an error if the json decoding fails
func tryReadFromStorage[V any](ctx context.Context, storage logical.Storage, key string) (value V, err error) {
	rawValue, err := storage.Get(ctx, key)

	if rawValue == nil || err != nil {
		return value, nil
	}

	if err := rawValue.DecodeJSON(&value); err != nil {
		return value, fmt.Errorf("error while decoding the storage entry: %s", err)
	}

	return value, nil
}

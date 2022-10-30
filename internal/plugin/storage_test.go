package plugin

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getTestStorage() logical.Storage {
	return &logical.InmemStorage{}
}

func TestSaveAndRead(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "save and read - string",
			fn: func(t *testing.T) {
				v := "value"
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[string](c, s, k, &v))
				actual, err := readFromStorage[string](c, s, k)
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(v, actual)
			},
		},
		{
			name: "save and read incorrect type",
			fn: func(t *testing.T) {
				v := "value"
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[string](c, s, k, &v))
				actual, err := readFromStorage[int](c, s, k)
				require.ErrorContains(t, err, "error while decoding")
				a := assert.New(t)
				a.EqualValues(0, actual)
			},
		},
		{
			name: "read incorrect property - string",
			fn: func(t *testing.T) {
				c := context.Background()
				s := getTestStorage()
				actual, err := readFromStorage[string](c, s, "some_property")
				require.ErrorContains(t, err, "error while reading")
				a := assert.New(t)
				a.EqualValues("", actual)
			},
		},
		{
			name: "save and read - struct",
			fn: func(t *testing.T) {
				v := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[accountTokenMetadata](c, s, k, &v))
				actual, err := readFromStorage[accountTokenMetadata](c, s, k)
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(v, actual)
			},
		},
		{
			name: "save and read incorrect type - struct",
			fn: func(t *testing.T) {
				v := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[accountTokenMetadata](c, s, k, &v))
				actual, err := readFromStorage[projectTokenMetadata](c, s, k)
				expected := projectTokenMetadata{
					Id:  "i1",
					TTL: 1 * time.Hour,
				}
				// decoding to wrong type does not return an error. Only fields with the same names are assigned
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "read incorrect property - struct",
			fn: func(t *testing.T) {
				c := context.Background()
				s := getTestStorage()
				actual, err := readFromStorage[accountTokenMetadata](c, s, "some_property")
				require.ErrorContains(t, err, "error while reading")
				a := assert.New(t)
				a.EqualValues(accountTokenMetadata{}, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestTryReadFromStorage(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "save and read - string",
			fn: func(t *testing.T) {
				v := "value"
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[string](c, s, k, &v))
				actual, err := tryReadFromStorage[string](c, s, k)
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(v, actual)
			},
		},
		{
			name: "save and read incorrect type",
			fn: func(t *testing.T) {
				v := "value"
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[string](c, s, k, &v))
				actual, err := tryReadFromStorage[int](c, s, k)
				require.ErrorContains(t, err, "error while decoding")
				a := assert.New(t)
				a.EqualValues(0, actual)
			},
		},
		{
			name: "read incorrect property - string",
			fn: func(t *testing.T) {
				c := context.Background()
				s := getTestStorage()
				actual, err := tryReadFromStorage[string](c, s, "some_property")
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues("", actual)
			},
		},
		{
			name: "save and read - struct",
			fn: func(t *testing.T) {
				v := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[accountTokenMetadata](c, s, k, &v))
				actual, err := tryReadFromStorage[accountTokenMetadata](c, s, k)
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(v, actual)
			},
		},
		{
			name: "save and read incorrect type - struct",
			fn: func(t *testing.T) {
				v := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				k := "key"
				c := context.Background()
				s := getTestStorage()
				require.NoError(t, saveToStorage[accountTokenMetadata](c, s, k, &v))
				actual, err := tryReadFromStorage[projectTokenMetadata](c, s, k)
				expected := projectTokenMetadata{
					Id:  "i1",
					TTL: 1 * time.Hour,
				}
				// decoding to wrong type does not return an error. Only fields with the same names are assigned
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "read incorrect property - struct",
			fn: func(t *testing.T) {
				c := context.Background()
				s := getTestStorage()
				actual, err := tryReadFromStorage[accountTokenMetadata](c, s, "some_property")
				require.NoError(t, err)
				a := assert.New(t)
				a.EqualValues(accountTokenMetadata{}, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

package plugin

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func updateConfig(b logical.Backend, r *logical.Request, d map[string]interface{}) error {
	r.Operation = logical.UpdateOperation
	r.Path = "config"
	r.Data = d
	res, err := b.HandleRequest(context.Background(), r)
	if err != nil {
		return err
	}
	if res != nil && res.IsError() {
		return res.Error()
	}
	return nil
}

func updateConfigSuccess(t *testing.T, b logical.Backend, r *logical.Request, d map[string]interface{}) {
	err := updateConfig(b, r, d)
	require.NoError(t, err)
}

func updateConfigError(t *testing.T, b logical.Backend, r *logical.Request, d map[string]interface{}, errContains string) {
	err := updateConfig(b, r, d)
	require.ErrorContains(t, err, errContains)
}

func readConfigSuccess(t *testing.T, r *logical.Request) configEntry {
	c, err := getConfig(context.Background(), r)
	require.NoError(t, err)
	return c
}

func readConfigError(t *testing.T, r *logical.Request) {
	_, err := getConfig(context.Background(), r)
	require.Error(t, err, "error while reading the storage entry")
}

func TestConfig(t *testing.T) {
	b, s := getTestBackend(t)
	r := &logical.Request{Storage: s}
	var expected configEntry
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "empty",
			fn: func(t *testing.T) {
				updateConfigError(
					t,
					b,
					r,
					map[string]interface{}{"some_garbage_value": 10},
					"missing data")
				readConfigError(t, r)
			},
		},
		{
			name: "invalid_url_with_protocol",
			fn: func(t *testing.T) {
				updateConfigError(
					t,
					b,
					r,
					map[string]interface{}{"argo_cd_url": "https://argocd.wfecd.splunk.lol", "admin_token": "some-dummy-token"},
					"invalid argo cd url")
				readConfigError(t, r)
			},
		},
		{
			name: "invalid_url_address",
			fn: func(t *testing.T) {
				updateConfigError(
					t,
					b,
					r,
					map[string]interface{}{"argo_cd_url": "argocdurl", "admin_token": "some-dummy-token"},
					"invalid argo cd url")
				readConfigError(t, r)
			},
		},
		{
			name: "missing_token",
			fn: func(t *testing.T) {
				expected.AccountTokenMaxTTL = 10 * time.Hour
				expected.ProjectTokenMaxTTL = 6 * time.Hour
				updateConfigError(
					t,
					b,
					r,
					map[string]interface{}{"argo_cd_url": "https://argocd.wfecd.splunk.lol"},
					"missing data")
				readConfigError(t, r)
			},
		},
		{
			name: "missing_url",
			fn: func(t *testing.T) {
				expected.AccountTokenMaxTTL = 10 * time.Hour
				expected.ProjectTokenMaxTTL = 6 * time.Hour
				updateConfigError(
					t,
					b,
					r,
					map[string]interface{}{"admin_token": "some-dummy-token"},
					"missing data")
				readConfigError(t, r)
			},
		},
		{
			name: "min_valid_config",
			fn: func(t *testing.T) {
				expected.ArgoCDUrl = "argocd.wfecd.splunk.lol"
				expected.AdminToken = "some-dummy-token"
				expected.AccountTokenMaxTTL = 6 * time.Hour
				expected.ProjectTokenMaxTTL = 6 * time.Hour
				updateConfigSuccess(t, b, r, map[string]interface{}{"argo_cd_url": "argocd.wfecd.splunk.lol", "admin_token": "some-dummy-token"})
				c := readConfigSuccess(t, r)
				require.EqualValues(t, expected, c)
			},
		},
		{
			name: "valid_config",
			fn: func(t *testing.T) {
				expected.ArgoCDUrl = "argocd.wfecd.splunk.lol"
				expected.AdminToken = "some-dummy-token"
				expected.AccountTokenMaxTTL = 10 * time.Hour
				expected.ProjectTokenMaxTTL = 11 * time.Hour
				updateConfigSuccess(
					t,
					b,
					r,
					map[string]interface{}{
						"argo_cd_url":           "argocd.wfecd.splunk.lol",
						"admin_token":           "some-dummy-token",
						"account_token_max_ttl": "10h",
						"project_token_max_ttl": "11h",
					})
				c := readConfigSuccess(t, r)
				require.EqualValues(t, expected, c)
			},
		},
		{
			name: "ttl_cap",
			fn: func(t *testing.T) {
				expected.ArgoCDUrl = "argocd.wfecd.splunk.lol"
				expected.AdminToken = "some-dummy-token"
				expected.AccountTokenMaxTTL = 12 * time.Hour
				expected.ProjectTokenMaxTTL = 12 * time.Hour
				updateConfigSuccess(
					t,
					b,
					r,
					map[string]interface{}{
						"argo_cd_url":           "argocd.wfecd.splunk.lol",
						"admin_token":           "some-dummy-token",
						"account_token_max_ttl": "50h",
						"project_token_max_ttl": "40h",
					})
				c := readConfigSuccess(t, r)
				require.EqualValues(t, expected, c)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}

	// now do a read op and make sure we get the correct config out without an API key
	res, err := b.HandleRequest(context.Background(), r)
	require.NoError(t, err)
	require.False(t, res.IsError())
	data := res.Data
	a := assert.New(t)
	a.EqualValues(expected.ArgoCDUrl, data["argo_cd_url"])
	a.EqualValues(expected.ProjectTokenMaxTTL.String(), data["project_token_max_ttl"])
	a.EqualValues(expected.AccountTokenMaxTTL.String(), data["account_token_max_ttl"])
	_, ok := data["admin_token"]
	a.False(ok)
}

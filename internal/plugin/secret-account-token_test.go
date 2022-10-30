package plugin

import (
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/account"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func deleteAccountTokenSuccess(b *backend, id string, accountName string, dtr *account.EmptyResponse) (*logical.Response, error) {
	accountClient := testAccountClient{
		createTokenResponse: nil,
		deleteTokenResponse: dtr,
		createTokenError:    nil,
		DeleteTokenError:    nil,
	}
	return b.deleteAccountToken(getTestAccountClientContext(&accountClient), id, accountName)
}

func deleteAccountTokenFailure(b *backend, id string, accountName string, error string) (*logical.Response, error) {
	accountClient := testAccountClient{
		createTokenResponse: nil,
		deleteTokenResponse: nil,
		createTokenError:    nil,
		DeleteTokenError:    fmt.Errorf(error),
	}
	return b.deleteAccountToken(getTestAccountClientContext(&accountClient), id, accountName)
}

func TestDeleteAccountToken(t *testing.T) {
	b, _ := getTestBackend(t)
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "success",
			fn: func(t *testing.T) {
				dtr := &account.EmptyResponse{}
				res, err := deleteAccountTokenSuccess(b, "some-id", "some-account", dtr)
				a := assert.New(t)
				require.NoError(t, err)
				a.True(res == nil)
			},
		},
		{
			name: "failure",
			fn: func(t *testing.T) {
				res, err := deleteAccountTokenFailure(b, "some-id", "some-account", "token does not exist")
				require.ErrorContains(t, err, "token does not exist")
				require.ErrorContains(t, res.Error(), "token does not exist")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

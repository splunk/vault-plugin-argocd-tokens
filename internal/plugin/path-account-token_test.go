package plugin

import (
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/account"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateAccountTokenSuccess(b *backend, accountName string, ctr *account.CreateTokenResponse, ttl time.Duration) (*logical.Response, error) {
	accountClient := testAccountClient{
		createTokenResponse: ctr,
		deleteTokenResponse: nil,
		createTokenError:    nil,
		DeleteTokenError:    nil,
	}
	return b.getAccountToken(getTestAccountClientContext(&accountClient), accountName, ttl)
}

func generateAccountTokenFailure(b *backend, accountName string, error string, ttl time.Duration) (*logical.Response, error) {
	accountClient := testAccountClient{
		createTokenResponse: nil,
		deleteTokenResponse: nil,
		createTokenError:    fmt.Errorf(error),
		DeleteTokenError:    nil,
	}
	return b.getAccountToken(getTestAccountClientContext(&accountClient), accountName, ttl)
}

func TestGenerateAccountToken(t *testing.T) {
	b, _ := getTestBackend(t)
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "success",
			fn: func(t *testing.T) {
				ttl := 1 * time.Hour
				ctr := &account.CreateTokenResponse{
					Token: "some-dummy-token",
				}
				res, err := generateAccountTokenSuccess(b, "some-account", ctr, ttl)
				a := assert.New(t)
				a.EqualValues("some-account", res.Data["account_name"])
				a.EqualValues("some-dummy-token", res.Data["token"])
				a.EqualValues("some-account", res.Secret.InternalData["account_name"])
				a.EqualValues(res.Data["id"], res.Secret.InternalData["id"])

				a.EqualValues(ttl.String(), res.Secret.LeaseOptions.TTL.String())
				a.EqualValues(false, res.Secret.LeaseOptions.Renewable)
				require.NoError(t, err)
				require.NoError(t, res.Error())

				// token is not saved in the lease data
				tokenInLeaseData, ok := res.Secret.InternalData["token"]
				a.True(tokenInLeaseData == nil)
				a.False(ok)
			},
		},
		{
			name: "failure",
			fn: func(t *testing.T) {
				ttl := 1 * time.Hour
				res, err := generateAccountTokenFailure(b, "some-account", "account does not exist", ttl)
				require.ErrorContains(t, err, "account does not exist")
				require.ErrorContains(t, res.Error(), "account does not exist")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

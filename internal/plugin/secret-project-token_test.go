package plugin

import (
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func deleteProjectTokenSuccess(b *backend, id string, projectName string, projectRoleName string, dtr *project.EmptyResponse) (*logical.Response, error) {
	projectClient := testProjectClient{
		createTokenResponse: nil,
		deleteTokenResponse: dtr,
		createTokenError:    nil,
		DeleteTokenError:    nil,
	}
	return b.deleteProjectToken(getTestProjectClientContext(&projectClient), id, projectName, projectRoleName)
}

func deleteProjectTokenFailure(b *backend, id string, projectName string, projectRoleName string, error string) (*logical.Response, error) {
	projectClient := testProjectClient{
		createTokenResponse: nil,
		deleteTokenResponse: nil,
		createTokenError:    nil,
		DeleteTokenError:    fmt.Errorf(error),
	}
	return b.deleteProjectToken(getTestProjectClientContext(&projectClient), id, projectName, projectRoleName)
}

func TestDeleteProjectToken(t *testing.T) {
	b, _ := getTestBackend(t)
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "success",
			fn: func(t *testing.T) {
				dtr := &project.EmptyResponse{}
				res, err := deleteProjectTokenSuccess(b, "some-id", "some-project", "some-role", dtr)
				a := assert.New(t)
				require.NoError(t, err)
				a.True(res == nil)
			},
		},
		{
			name: "failure",
			fn: func(t *testing.T) {
				res, err := deleteProjectTokenFailure(b, "some-id", "some-project", "some-role", "token does not exist")
				require.ErrorContains(t, err, "token does not exist")
				require.ErrorContains(t, res.Error(), "token does not exist")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

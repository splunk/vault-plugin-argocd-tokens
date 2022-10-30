package plugin

import (
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateProjectTokenSuccess(b *backend, projectName string, projectRoleName string, ctr *project.ProjectTokenResponse, ttl time.Duration) (*logical.Response, error) {
	projectClient := testProjectClient{
		createTokenResponse: ctr,
		deleteTokenResponse: nil,
		createTokenError:    nil,
		DeleteTokenError:    nil,
	}
	return b.getProjectToken(getTestProjectClientContext(&projectClient), projectName, projectRoleName, ttl)
}

func generateProjectTokenFailure(b *backend, projectName string, projectRoleName string, error string, ttl time.Duration) (*logical.Response, error) {
	projectClient := testProjectClient{
		createTokenResponse: nil,
		deleteTokenResponse: nil,
		createTokenError:    fmt.Errorf(error),
		DeleteTokenError:    nil,
	}
	return b.getProjectToken(getTestProjectClientContext(&projectClient), projectName, projectRoleName, ttl)
}

func TestGenerateProjectToken(t *testing.T) {
	b, _ := getTestBackend(t)
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "success",
			fn: func(t *testing.T) {
				ttl := 1 * time.Hour
				ctr := &project.ProjectTokenResponse{
					Token: "some-dummy-token",
				}
				res, err := generateProjectTokenSuccess(b, "some-project", "some-role", ctr, ttl)
				a := assert.New(t)
				a.EqualValues("some-project", res.Data["project_name"])
				a.EqualValues("some-role", res.Data["project_role_name"])
				a.EqualValues("some-dummy-token", res.Data["token"])
				a.EqualValues("some-project", res.Secret.InternalData["project_name"])
				a.EqualValues("some-role", res.Secret.InternalData["project_role_name"])
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
				res, err := generateProjectTokenFailure(b, "some-project", "some-role", "project does not exist", ttl)
				require.ErrorContains(t, err, "project does not exist")
				require.ErrorContains(t, res.Error(), "project does not exist")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

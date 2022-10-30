package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var getProjectTokenSchema = map[string]*framework.FieldSchema{
	fldProjectName: {
		Type:        framework.TypeString,
		Description: `ArgoCD Project name`,
	},
	fldProjectRoleName: {
		Type:        framework.TypeString,
		Description: `ArgoCD Project Role name`,
	},
	fldTTL: {
		Type:        framework.TypeDurationSecond,
		Description: `Expires in (default: 1h, max: 12h)`,
	},
}

func pathProjectToken(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Sprintf("project/%s/role/%s", framework.GenericNameRegex(fldProjectName), framework.GenericNameRegex(fldProjectRoleName)),
			Fields:  getProjectTokenSchema,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.getProjectTokenCallback,
					Summary:  "gets a token for a project role",
				},
			},
			HelpSynopsis:    trimHelp(helpPathProjectSynopsis),
			HelpDescription: trimHelp(helpPathProjectDescription),
		},
	}
}

func (b *backend) getProjectTokenCallback(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	projectName, err := getFromFieldData[string](data, fldProjectName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	projectRoleName, err := getFromFieldData[string](data, fldProjectRoleName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	config, err := getConfig(ctx, req)
	if err != nil {
		errMsg := fmt.Sprintf("error while reading config: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	ttl := getTTLFromFieldData(data, fldTTL, 1*time.Hour, config.ProjectTokenMaxTTL)

	clientCtx, err := NewProjectClient(ctx, &config)
	if err != nil {
		errMsg := fmt.Sprintf("error while creating a new project client: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	return b.getProjectToken(clientCtx, projectName, projectRoleName, ttl)
}

func (b *backend) getProjectToken(
	clientCtx *projectClientContext,
	projectName string,
	projectRoleName string,
	ttl time.Duration) (*logical.Response, error) {
	token, err := clientCtx.GenerateToken(projectName, projectRoleName, ttl)
	if err != nil {
		errMsg := fmt.Sprintf("error while creating a new token for project role(%s/%s): %s", projectName, projectRoleName, err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	err = clientCtx.closer.Close()
	if err != nil {
		b.logger.Error(err.Error())
	}

	response := newTokenSecret(projectTokenSecretType, token.metadata.TTL).Response(token.toResponseData(), token.toLeaseData())

	return response, nil
}

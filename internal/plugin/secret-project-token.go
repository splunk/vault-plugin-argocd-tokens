package plugin

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func secretProjectToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type:   projectTokenSecretType,
		Fields: map[string]*framework.FieldSchema{},
		Revoke: b.deleteProjectTokenCallback,
	}
}

// Project Tokens need to be deleted from argo cd as
// -- argo cd saves token metadata for all projects in their specific k8s apprpoj resource
// -- argo cd does not clear metadata from the appproj resource after the tokens expire
// -- So the yaml manifest for the apprpoj should be able to hold the metadata for all the tokens for a given project
// -- If we don't clear expired tokens from the apprpoj resource, then the ephemeral token approach can make argo cd perform slower or bring it down completely
// -- Calling this delete endpoint on a schedule or from the backend.PeriodicFunc would ensure that appproj resource only keeps required tokens in the manifest
func (b *backend) deleteProjectTokenCallback(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id, err := getFromData[string](req.Secret.InternalData, fldID)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	projectName, err := getFromData[string](req.Secret.InternalData, fldProjectName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	projectRoleName, err := getFromData[string](req.Secret.InternalData, fldProjectRoleName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	config, err := getConfig(ctx, req)
	if err != nil {
		errMsg := fmt.Sprintf("error while reading config: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), fmt.Errorf(errMsg)
	}

	clientCtx, err := NewProjectClient(ctx, &config)
	if err != nil {
		errMsg := fmt.Sprintf("error while creating a new account client: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	return b.deleteProjectToken(clientCtx, id, projectName, projectRoleName)
}

func (b *backend) deleteProjectToken(clientCtx *projectClientContext, id string, projectName string, projectRoleName string) (*logical.Response, error) {
	if err := clientCtx.DeleteToken(id, projectName, projectRoleName); err != nil {
		errMsg := fmt.Sprintf("error while deleting token(%s) for project/role(%s/%s): %s", id, projectName, projectRoleName, err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	err := clientCtx.closer.Close()
	if err != nil {
		b.logger.Error(err.Error())
	}

	return nil, nil
}

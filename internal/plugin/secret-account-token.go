package plugin

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func secretAccountToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type:   accountTokenSecretType,
		Fields: map[string]*framework.FieldSchema{},
		Revoke: b.deleteAccountTokenCallback,
	}
}

// Account Tokens need to be deleted from argo cd as
// -- argo cd saves token metadata for all accounts in a k8s secret: argocd-secret
// -- argo cd does not clear metadata from the k8s secret after the tokens expire
// -- So the yaml manifest for the argocd-secret should be able to hold the metadata for all the tokens for all the accounts
// -- If we don't clear expired tokens from the k8s secret, then the ephemeral token approach can make argo cd perform slower or bring it down completely
// -- Calling this delete endpoint on a schedule or from the backend.PeriodicFunc would ensure that argocd-secret only keeps required tokens in the manifest
func (b *backend) deleteAccountTokenCallback(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountName, err := getFromData[string](req.Secret.InternalData, fldAccountName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	id, err := getFromData[string](req.Secret.InternalData, fldID)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	config, err := getConfig(ctx, req)
	if err != nil {
		errMsg := fmt.Sprintf("error while reading config: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), fmt.Errorf(errMsg)
	}

	clientCtx, err := NewAccountClient(ctx, &config)
	if err != nil {
		errMsg := fmt.Sprintf("error while creating a new account client: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	return b.deleteAccountToken(clientCtx, id, accountName)
}

func (b *backend) deleteAccountToken(clientCtx *accountClientContext, id string, accountName string) (*logical.Response, error) {
	if err := clientCtx.DeleteToken(id, accountName); err != nil {
		errMsg := fmt.Sprintf("error while deleting token(%s) for account(%s): %s", id, accountName, err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	err := clientCtx.closer.Close()
	if err != nil {
		b.logger.Error(err.Error())
	}

	return nil, nil
}

package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var getAccountTokenSchema = map[string]*framework.FieldSchema{
	fldAccountName: {
		Type:        framework.TypeString,
		Description: `ArgoCD Account name`,
	},
	fldTTL: {
		Type:        framework.TypeDurationSecond,
		Description: `Expires in (default: 1h, max: 180d)`,
	},
}

func pathAccountToken(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Sprintf("account/%s", framework.GenericNameRegex(fldAccountName)),
			Fields:  getAccountTokenSchema,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.getAccountTokenCallback,
					Summary:  "gets a token for an argo cd account",
				},
			},
			HelpSynopsis:    trimHelp(helpPathAccountSynopsis),
			HelpDescription: trimHelp(helpPathAccountDescription),
		},
	}
}

func (b *backend) getAccountTokenCallback(
	ctx context.Context,
	req *logical.Request,
	data *framework.FieldData) (*logical.Response, error) {

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

	accountName, err := getFromFieldData[string](data, fldAccountName)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error while getting account name from data: %s", err)), err
	}

	ttl := getTTLFromFieldData(data, fldTTL, 1*time.Hour, config.AccountTokenMaxTTL)

	return b.getAccountToken(clientCtx, accountName, ttl)
}

func (b *backend) getAccountToken(
	clientCtx *accountClientContext,
	accountName string,
	ttl time.Duration) (*logical.Response, error) {
	token, err := clientCtx.GenerateToken(accountName, ttl)
	if err != nil {
		b.logger.Error(err.Error())
		return logical.ErrorResponse(err.Error()), err
	}

	err = clientCtx.closer.Close()
	if err != nil {
		b.logger.Error(err.Error())
	}
	
	response := newTokenSecret(accountTokenSecretType, token.metadata.TTL).Response(token.toResponseData(), token.toLeaseData())

	return response, nil
}

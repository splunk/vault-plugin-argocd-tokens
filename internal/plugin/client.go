package plugin

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/account"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/google/uuid"
)

const (
	totalRetries = 4
)

var retryWaitSeconds = []time.Duration{0, 3, 5, 10} // first value should remain 0

type projectClientContext struct {
	client        project.ProjectServiceClient
	clientContext context.Context
	closer        io.Closer
}

type accountClientContext struct {
	client        account.AccountServiceClient
	clientContext context.Context
	closer        io.Closer
}

type accountTokenMetadata struct {
	Id          string        `json:"id" structs:"id" mapstructure:"id"`
	AccountName string        `json:"account_name" structs:"account_name" mapstructure:"account_name"`
	TTL         time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
}

type projectTokenMetadata struct {
	Id          string        `json:"id" structs:"id" mapstructure:"id"`
	ProjectName string        `json:"project_name" structs:"project_name" mapstructure:"project_name"`
	RoleName    string        `json:"role_name" structs:"role_name" mapstructure:"role_name"`
	TTL         time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
}

type accountToken struct {
	metadata accountTokenMetadata
	token    string
}

type projectToken struct {
	metadata projectTokenMetadata
	token    string
}

func (c *configEntry) toClientOptions() *apiclient.ClientOptions {
	clientOptions := apiclient.ClientOptions{
		ServerAddr:   c.ArgoCDUrl,
		AuthToken:    c.AdminToken,
		GRPCWeb:      true,
		HttpRetryMax: 3,
		Insecure:     c.Insecure,
		PlainText:    c.Plaintext,
	}

	return &clientOptions
}

func NewProjectClient(ctx context.Context, config *configEntry) (*projectClientContext, error) {
	clientOptions := config.toClientOptions()

	client, err := apiclient.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error while creating new apiClient: %s", err)
	}

	closer, projectClient, err := client.NewProjectClient()

	if err != nil {
		return nil, fmt.Errorf("error while creating new projectClient: %s", err)
	}

	clientContext := projectClientContext{
		client:        projectClient,
		clientContext: ctx,
		closer:        closer,
	}

	return &clientContext, nil
}

func NewAccountClient(ctx context.Context, config *configEntry) (*accountClientContext, error) {
	clientOptions := config.toClientOptions()

	client, err := apiclient.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error while creating new apiClient: %s", err)
	}

	closer, accountClient, err := client.NewAccountClient()
	if err != nil {
		return nil, fmt.Errorf("error while creating new accountClient: %s", err)
	}

	clientContext := accountClientContext{
		client:        accountClient,
		clientContext: ctx,
		closer:        closer,
	}

	return &clientContext, nil
}

func (clientCtx *projectClientContext) GenerateToken(projectName string, projectRoleName string, expiresIn time.Duration) (*projectToken, error) {
	retries := 0
	var response *project.ProjectTokenResponse
	var err error

	for retries < totalRetries {
		id := uuid.New().String()
		time.Sleep(retryWaitSeconds[retries] * time.Second)
		createTokenRequest := &project.ProjectTokenCreateRequest{
			Project:   projectName,
			Role:      projectRoleName,
			ExpiresIn: toDurationSeconds(expiresIn),
			Id:        id,
		}

		projectClient := clientCtx.client
		response, err = projectClient.CreateToken(clientCtx.clientContext, createTokenRequest)

		if err == nil {
			token := projectToken{
				metadata: projectTokenMetadata{
					Id:          id,
					ProjectName: projectName,
					RoleName:    projectRoleName,
					TTL:         expiresIn,
				},
				token: response.Token,
			}
		
			return &token, nil
		}
		retries++
	}
	return nil, fmt.Errorf("Error in Generate token for projectClient: %s", err)

}

func (clientCtx *accountClientContext) GenerateToken(accountName string, expiresIn time.Duration) (*accountToken, error) {
	retries := 0
	var response *account.CreateTokenResponse
	var err error

	for retries < totalRetries {
		id := uuid.New().String()
		time.Sleep(retryWaitSeconds[retries] * time.Second)
		createTokenRequest := &account.CreateTokenRequest{
			Name:      accountName,
			ExpiresIn: toDurationSeconds(expiresIn),
			Id:        id,
		}

		accountClient := clientCtx.client
		response, err = accountClient.CreateToken(clientCtx.clientContext, createTokenRequest)
		if err == nil {
			token := accountToken{
				metadata: accountTokenMetadata{
					Id:          id,
					AccountName: accountName,
					TTL:         expiresIn,
				},
				token: response.Token,
			}
		
			return &token, nil
		}
		retries++
	}
	return nil, fmt.Errorf("Error in Generate token for accountClient: %s", err)

}

func (clientCtx *accountClientContext) DeleteToken(tokenId string, accountName string) error {
	deleteTokenRequest := &account.DeleteTokenRequest{
		Name: accountName,
		Id:   tokenId,
	}

	accountClient := clientCtx.client
	_, err := accountClient.DeleteToken(clientCtx.clientContext, deleteTokenRequest)

	if err != nil {
		return fmt.Errorf("error in delete token for accountClient: %s", err)
	}

	return nil
}

func (clientCtx *projectClientContext) DeleteToken(tokenId string, projectName string, roleName string) error {
	deleteTokenRequest := &project.ProjectTokenDeleteRequest{
		Project: projectName,
		Role:    roleName,
		Id:      tokenId,
	}

	projectClient := clientCtx.client
	_, err := projectClient.DeleteToken(clientCtx.clientContext, deleteTokenRequest)

	if err != nil {
		return fmt.Errorf("error in delete token for projectClient: %s", err)
	}

	return nil
}

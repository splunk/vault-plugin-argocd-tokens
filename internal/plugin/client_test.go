package plugin

import (
	"context"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/account"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
)

type testCloser struct{}

func (c testCloser) Close() error {
	return nil
}

type testAccountClient struct {
	createTokenResponse *account.CreateTokenResponse
	deleteTokenResponse *account.EmptyResponse
	createTokenError    error
	DeleteTokenError    error
}

func (client *testAccountClient) CreateToken(ctx context.Context, in *account.CreateTokenRequest, opts ...grpc.CallOption) (*account.CreateTokenResponse, error) {
	return client.createTokenResponse, client.createTokenError
}
func (client *testAccountClient) DeleteToken(ctx context.Context, in *account.DeleteTokenRequest, opts ...grpc.CallOption) (*account.EmptyResponse, error) {
	return client.deleteTokenResponse, client.DeleteTokenError
}
func (client *testAccountClient) CanI(ctx context.Context, in *account.CanIRequest, opts ...grpc.CallOption) (*account.CanIResponse, error) {
	return nil, nil
}
func (client *testAccountClient) UpdatePassword(ctx context.Context, in *account.UpdatePasswordRequest, opts ...grpc.CallOption) (*account.UpdatePasswordResponse, error) {
	return nil, nil
}
func (client *testAccountClient) ListAccounts(ctx context.Context, in *account.ListAccountRequest, opts ...grpc.CallOption) (*account.AccountsList, error) {
	return nil, nil
}
func (client *testAccountClient) GetAccount(ctx context.Context, in *account.GetAccountRequest, opts ...grpc.CallOption) (*account.Account, error) {
	return nil, nil
}

func getTestAccountClientContext(accountClient *testAccountClient) *accountClientContext {
	return &accountClientContext{
		client:        accountClient,
		clientContext: context.Background(),
		closer:        testCloser{},
	}
}

type testProjectClient struct {
	createTokenResponse *project.ProjectTokenResponse
	deleteTokenResponse *project.EmptyResponse
	createTokenError    error
	DeleteTokenError    error
}

func (client *testProjectClient) CreateToken(ctx context.Context, in *project.ProjectTokenCreateRequest, opts ...grpc.CallOption) (*project.ProjectTokenResponse, error) {
	return client.createTokenResponse, client.createTokenError
}
func (client *testProjectClient) DeleteToken(ctx context.Context, in *project.ProjectTokenDeleteRequest, opts ...grpc.CallOption) (*project.EmptyResponse, error) {
	return client.deleteTokenResponse, client.DeleteTokenError
}
func (client *testProjectClient) Create(ctx context.Context, in *project.ProjectCreateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	return nil, nil
}
func (client *testProjectClient) List(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProjectList, error) {
	return nil, nil
}
func (client *testProjectClient) GetDetailedProject(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*project.DetailedProjectsResponse, error) {
	return nil, nil
}
func (client *testProjectClient) Get(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	return nil, nil
}
func (client *testProjectClient) GetGlobalProjects(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*project.GlobalProjectsResponse, error) {
	return nil, nil
}
func (client *testProjectClient) Update(ctx context.Context, in *project.ProjectUpdateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	return nil, nil
}
func (client *testProjectClient) Delete(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*project.EmptyResponse, error) {
	return nil, nil
}
func (client *testProjectClient) ListEvents(ctx context.Context, in *project.ProjectQuery, opts ...grpc.CallOption) (*v1.EventList, error) {
	return nil, nil
}
func (client *testProjectClient) GetSyncWindowsState(ctx context.Context, in *project.SyncWindowsQuery, opts ...grpc.CallOption) (*project.SyncWindowsResponse, error) {
	return nil, nil
}

func (client *testProjectClient) ListLinks(ctx context.Context, in *project.ListProjectLinksRequest, opts ...grpc.CallOption) (*application.LinksResponse, error) {
	return nil, nil
}

func getTestProjectClientContext(projectClient *testProjectClient) *projectClientContext {
	return &projectClientContext{
		client:        projectClient,
		clientContext: context.Background(),
		closer:        testCloser{},
	}
}

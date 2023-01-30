package plugin

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"net/url"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"strings"
	"time"
)

const (
	cfgStorageKey            = "config"
	cfgFldArgoCdUrl          = "argo_cd_url"
	cfgFldAdminToken         = "admin_token"
	cfgFldAccountTokenMaxTTL = "account_token_max_ttl"
	cfgFldProjectTokenMaxTTL = "project_token_max_ttl"
	cfgFldInsecure           = "insecure"
	cfgFldPlaintext          = "plaintext"
	fldAccountName           = "account_name"
	fldProjectName           = "project_name"
	fldProjectRoleName       = "project_role_name"
	fldTTL                   = "ttl"
	fldID                    = "id"
	fldToken                 = "token"
	accountTokenSecretType   = "account_token_secret"
	projectTokenSecretType   = "project_token_secret"
)

// configEntry represents the vault config
type configEntry struct {
	ArgoCDUrl          string        `json:"argo_cd_url" structs:"argo_cd_url" mapstructure:"argo_cd_url"`
	AdminToken         string        `json:"admin_token" structs:"admin_token" mapstructure:"admin_token"`
	AccountTokenMaxTTL time.Duration `json:"account_token_max_ttl" structs:"account_token_max_ttl" mapstructure:"account_token_max_ttl"`
	ProjectTokenMaxTTL time.Duration `json:"project_token_max_ttl" structs:"project_token_max_ttl" mapstructure:"project_token_max_ttl"`
	Insecure           bool          `json:"insecure" structs:"insecure" mapstructure:"insecure"`
	Plaintext          bool          `json:"plaintext" structs:"plaintext" mapstructure:"plaintext"`
}

// toResponse returns the logical response corresponding to the config entry, ensuring that the Admin Token is not exposed
func (c *configEntry) toResponse() *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			cfgFldArgoCdUrl:          c.ArgoCDUrl,
			cfgFldAccountTokenMaxTTL: c.AccountTokenMaxTTL.String(),
			cfgFldProjectTokenMaxTTL: c.ProjectTokenMaxTTL.String(),
			cfgFldInsecure:           c.Insecure,
			cfgFldPlaintext:          c.Plaintext,
		},
	}
}

var configSchema = map[string]*framework.FieldSchema{
	cfgFldArgoCdUrl: {
		Type:        framework.TypeString,
		Description: `Argo CD Instance Url`,
	},
	cfgFldAdminToken: {
		Type:        framework.TypeString,
		Description: `Argo CD Instance Account Token with admin role`,
	},
	cfgFldAccountTokenMaxTTL: {
		Type:        framework.TypeDurationSecond,
		Description: `Max TTL for account tokens`,
	},
	cfgFldProjectTokenMaxTTL: {
		Type:        framework.TypeDurationSecond,
		Description: `Max TTL for project tokens`,
	},
	cfgFldInsecure: {
		Type:        framework.TypeBool,
		Description: `Argo CD insecure connection (This should not be used in production environment)`,
	},
	cfgFldPlaintext: {
		Type:        framework.TypeBool,
		Description: `Argo CD plaintext communication (This should not be used in production environments)`,
	},
}

// initFromInputs initializes the entry from partial input data
func (c *configEntry) initFromInputs(data *framework.FieldData) error {
	var allErorrs error
	argoCDURL, err := getFromFieldData[string](data, cfgFldArgoCdUrl)
	if err != nil {
		allErorrs = errors.Wrap(err)
	}

	adminToken, err := getFromFieldData[string](data, cfgFldAdminToken)
	if err != nil {
		allErorrs = errors.Wrap(err)
	}

	//Explicitly set insecure to false by default if not provided
	insecure, insecureErr := getFromFieldData[bool](data, cfgFldInsecure)
	if insecureErr != nil {
		insecure = false
	}

	//Explicitly set plaintext to false by default if not provided
	plaintext, plaintextErr := getFromFieldData[bool](data, cfgFldPlaintext)
	if plaintextErr != nil {
		plaintext = false
	}

	c.AccountTokenMaxTTL = getTTLFromFieldData(data, cfgFldAccountTokenMaxTTL, 6*time.Hour, 12*time.Hour)
	c.ProjectTokenMaxTTL = getTTLFromFieldData(data, cfgFldProjectTokenMaxTTL, 6*time.Hour, 12*time.Hour)

	if err != nil {
		allErorrs = errors.Wrap(err)
	}

	if allErorrs != nil {
		return allErorrs
	}

	c.AdminToken = adminToken
	c.ArgoCDUrl = argoCDURL
	c.Insecure = insecure
	c.Plaintext = plaintext

	return c.assertValid()
}

// getConfig returns the configuration from storage or an empty object if not found
func getConfig(ctx context.Context, req *logical.Request) (configEntry, error) {
	return readFromStorage[configEntry](ctx, req.Storage, cfgStorageKey)
}

// pathConfigRead implements read on the /config path
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	cfg, err := getConfig(ctx, req)

	if err != nil {
		errMsg := fmt.Sprintf("error while reading config from storage: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}
	return cfg.toResponse(), nil
}

// pathConfigWrite implements write on the /config path
func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := tryReadFromStorage[configEntry](ctx, req.Storage, cfgStorageKey)
	if err != nil {
		errMsg := fmt.Sprintf("error while reading config from storage: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	if err := cfg.initFromInputs(data); err != nil {
		errMsg := fmt.Sprintf("error while init in config: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	if cfg.Insecure {
		b.logger.Warn(fmt.Sprintf("ArgoCD server (%s) configured with insecure connection. This should NOT be used in a production environment!", cfg.ArgoCDUrl))
	}
	if cfg.Plaintext {
		b.logger.Warn(fmt.Sprintf("ArgoCD server (%s) configured with plaintext communication. This should NOT be used in a production environment!", cfg.ArgoCDUrl))
	}

	if err := saveToStorage[configEntry](ctx, req.Storage, cfgStorageKey, &cfg); err != nil {
		errMsg := fmt.Sprintf("error while writing config to storage: %s", err)
		b.logger.Error(errMsg)
		return logical.ErrorResponse(errMsg), err
	}

	return cfg.toResponse(), nil
}

// pathConfig configures operations on the /config path
func pathConfig(b *backend) []*framework.Path {
	paths := []*framework.Path{
		{
			Pattern: "config$",
			Fields:  configSchema,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:    b.pathConfigRead,
					Summary:     "retrieves the argo cd token plugin configuration",
					Description: `returns the configuration for the specific argo cd tokens plugin mount. Does not expose the API key`,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback:    b.pathConfigWrite,
					Summary:     "updates the argo cd tokens plugin configuration",
					Description: `updates the configuration for the specific argo cd tokens plugin mount`,
				},
			},
			HelpSynopsis:    trimHelp(helpPathConfigSynopsis),
			HelpDescription: trimHelp(helpPathConfigDescription),
		},
	}
	return paths
}

func (c *configEntry) assertValid() error {
	_, err := url.Parse(c.ArgoCDUrl)
	if err != nil ||
		strings.Contains(c.ArgoCDUrl, "http") ||
		strings.Contains(c.ArgoCDUrl, "tcp") ||
		!strings.Contains(c.ArgoCDUrl, ".") {
		return fmt.Errorf("invalid argo cd url: argo cd url(%s) should only contain the address without protocol", c.ArgoCDUrl)
	}

	return nil
}

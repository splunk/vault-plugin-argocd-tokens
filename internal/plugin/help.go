package plugin

const helpBackend = `
This backend can create argo cd account tokens and project tokens using the config provided.
For each argo cd instance that the backend needs to connect to should be enabled to a different path,
and the write engine-path/config endpoint should be called first to setup the config.
Once the config is set, the account and project paths can be used to create the ephemeral tokens.
`

const helpPathConfigSynopsis = `
Configures argo cd connection for the given path
`

const helpPathConfigDescription = `
config properties:
vault write engine-path/config "key1=value1" "key2=value2"
keys:
argo_cd_url: URL for the argo cd instance (do not add https in front of the URL)
admin_token: Token for an account that has admin access for the given argo cd instance
account_token_max_ttl: Max TTL for the account tokens created from this plugin
project_token_max_ttl: Max TTL for the project tokens created from this plugin
`

const helpPathAccountSynopsis = `
Create tokens for the given argo cd account
`

const helpPathAccountDescription = `
- vault write engine-path/account/account-name expires_in=2h
-- creates a token for the specified account
-- Default value for expires_in=1h
-- returns created token
-- when the token expires, it is removed from argo cd
`
const helpPathProjectSynopsis = `
Create tokens for the given argo cd project role
`

const helpPathProjectDescription = `
- vault write engine-path/project/project_name/role/role_name expires_in=2h
-- creates a token for the specified role in an argo cd project
-- Default value for expires_in=1h
-- returns created token
-- when the token expires, it is removed from argo cd
`

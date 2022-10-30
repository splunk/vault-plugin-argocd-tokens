pushd $(dirname "$0")
vault login root
source ./common.sh

echo -- Enable wfecd engine ---
vault secrets enable -path=wfe vault-plugin-argocd-tokens

echo -- List secrets engine ---
vault secrets list

echo --- Write config ---
vault write wfe/config \
  "argo_cd_url=${WFECDSTG_SERVER}" \
  "admin_token=${WFECDSTG_TOKEN}" \
  "account_token_max_ttl=1h" \
  "project_token_max_ttl=1h"

echo -- Read config --
vault read wfe/config
popd
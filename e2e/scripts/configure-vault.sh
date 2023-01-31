sleep 20

password=$(kubectl get secret -n argocd-tokens-vault-plugin-testing argocd-initial-admin-secret --output=json | jq -r '.data.password' | base64 --decode)

argocd login --insecure --username=admin --password=$password $ARGOCD_SERVER

ARGOCD_TOKEN=$(argocd --insecure account generate-token -a argocd-tokens-plugin -e 720h)

vault login root 

vault secrets enable -path="${ENGINE_PATH}" vault-plugin-argocd-tokens

vault secrets list

vault write "${ENGINE_PATH}"/config "argo_cd_url=${ARGOCD_SERVER}" "admin_token=${ARGOCD_TOKEN}" "account_token_max_ttl=1h" "project_token_max_ttl=1h" "insecure=true" "plaintext=true"

vault read "${ENGINE_PATH}"/config
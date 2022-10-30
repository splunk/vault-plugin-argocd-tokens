#Build the plugin
rm -rf tmp/plugins
mkdir -p tmp/plugins
go build -o tmp/plugins/ ./cmd/vault-plugin-argocd-tokens
vault server -dev -dev-root-token-id=root -dev-plugin-dir=./tmp/plugins -log-level=debug

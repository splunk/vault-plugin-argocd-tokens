.DEFAULT_GOAL := all
.PHONY: all
all: build test lint

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./plugins/ ./cmd/vault-plugin-argocd-tokens

.PHONY: install
install:
	go install ./...

.PHONY: test
test:
	go test ./... -coverprofile cover.out

.PHONY: lint
lint:
	golangci-lint run --timeout=3m

.PHONY: deploy-argocd
deploy-argocd:
	echo deploying Argocd
	kubectl create namespace argocd-tokens-vault-plugin-testing --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./e2e/manifests/argocd/crd.yaml
	kustomize build ./e2e/manifests/argocd > ./build.yaml
	kubectl -n argocd-tokens-vault-plugin-testing apply -f ./build.yaml
	kubectl wait --timeout=2m -n argocd-tokens-vault-plugin-testing --all --for=jsonpath='{.status.phase}'=Running pod

.PHONY: deploy-vault
deploy-vault:
	echo Deploying vault
	docker pull ghcr.io/splunk/workflow-engine-base:2.0.12
	kubectl create namespace argocd-tokens-vault-plugin-testing --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./e2e/manifests/vault
	kubectl wait --timeout=5m -n argocd-tokens-vault-plugin-testing --for=jsonpath='{.status.phase}'=Running pod vault
	kubectl cp ./plugins/ argocd-tokens-vault-plugin-testing/vault:/
	kubectl cp ./e2e/scripts argocd-tokens-vault-plugin-testing/vault:/
	kubectl cp ./e2e/scenarios argocd-tokens-vault-plugin-testing/vault:/
	kubectl cp ~/.kube argocd-tokens-vault-plugin-testing/vault:/root
	if [ -d ~/.minikube ]; then \
		kubectl exec -n argocd-tokens-vault-plugin-testing vault -- mkdir -p ${HOME}; \
		kubectl cp ~/.minikube/ argocd-tokens-vault-plugin-testing/vault:${HOME}; \
	fi

.PHONY: e2e
e2e: build deploy-argocd deploy-vault
	while ! kubectl exec -n argocd-tokens-vault-plugin-testing vault -- bash -c 'curl -Ss $$ARGOCD_SERVER 2>&1 > /dev/null'; do \
		echo "Waiting for ArgoCD to be ready"; \
		sleep 1; \
	done
	while ! kubectl exec -n argocd-tokens-vault-plugin-testing vault -- bash -c 'curl -Ss $$VAULT_ADDR 2>&1 > /dev/null'; do \
		echo "Waiting for Vault to be ready"; \
		sleep 1; \
	done
	mkdir -p ./e2e/logs
	kubectl exec -n argocd-tokens-vault-plugin-testing vault -- bash ./scripts/configure-vault.sh
	kubectl exec -n argocd-tokens-vault-plugin-testing vault -- bash ./scripts/run-scenarios.sh
	kubectl logs -n argocd-tokens-vault-plugin-testing -l app.kubernetes.io/name=argocd-server > ./e2e/logs/argocd-server.log
	kubectl logs -n argocd-tokens-vault-plugin-testing vault > ./e2e/logs/vault.log

.PHONY: destroy
destroy:
	kubectl delete --force --grace-period=0 -f ./e2e/manifests/vault --ignore-not-found=true
	kubectl delete --force --grace-period=0 -f ./e2e/manifests/argocd/crd.yaml --ignore-not-found=true
	kubectl delete -n argocd-tokens-vault-plugin-testing --force --grace-period=0 -f ./build.yaml --ignore-not-found=true
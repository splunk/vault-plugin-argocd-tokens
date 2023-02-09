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
	kubectl create namespace vault-plugin-argocd-tokens-testing --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./e2e/manifests/argocd/crd.yaml
	kustomize build ./e2e/manifests/argocd > ./build.yaml
	kubectl -n vault-plugin-argocd-tokens-testing apply -f ./build.yaml
	kubectl wait --timeout=2m -n vault-plugin-argocd-tokens-testing --all --for=jsonpath='{.status.phase}'=Running pod

.PHONY: deploy-vault
deploy-vault:
	docker pull ghcr.io/splunk/workflow-engine-base:2.0.12
	kubectl create namespace vault-plugin-argocd-tokens-testing --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./e2e/manifests/vault
	kubectl wait --timeout=5m -n vault-plugin-argocd-tokens-testing --for=jsonpath='{.status.phase}'=Running pod vault
	kubectl cp ./plugins/ vault-plugin-argocd-tokens-testing/vault:/
	kubectl cp ./e2e/scripts vault-plugin-argocd-tokens-testing/vault:/
	kubectl cp ./e2e/scenarios vault-plugin-argocd-tokens-testing/vault:/
	kubectl cp ~/.kube vault-plugin-argocd-tokens-testing/vault:/root
	if [ -d ~/.minikube ]; then \
		kubectl exec -n vault-plugin-argocd-tokens-testing vault -- mkdir -p ${HOME}; \
		kubectl cp ~/.minikube/ vault-plugin-argocd-tokens-testing/vault:${HOME}; \
	fi

.PHONY: e2e
e2e: build deploy-argocd deploy-vault
	while ! kubectl exec -n vault-plugin-argocd-tokens-testing vault -- bash -c 'curl -Ss $$ARGOCD_SERVER 2>&1 > /dev/null'; do \
		echo "Waiting for ArgoCD to be ready"; \
		sleep 1; \
	done
	while ! kubectl exec -n vault-plugin-argocd-tokens-testing vault -- bash -c 'curl -Ss $$VAULT_ADDR 2>&1 > /dev/null'; do \
		echo "Waiting for Vault to be ready"; \
		sleep 1; \
	done
	mkdir -p ./e2e/logs
	kubectl exec -n vault-plugin-argocd-tokens-testing vault -- bash ./scripts/configure-vault.sh
	kubectl exec -n vault-plugin-argocd-tokens-testing vault -- bash ./scripts/run-scenarios.sh
	kubectl logs -n vault-plugin-argocd-tokens-testing -l app.kubernetes.io/name=argocd-server > ./e2e/logs/argocd-server.log
	kubectl logs -n vault-plugin-argocd-tokens-testing vault > ./e2e/logs/vault.log

.PHONY: destroy
destroy:
	kubectl delete --force --grace-period=0 -f ./e2e/manifests/vault --ignore-not-found=true
	kubectl delete --force --grace-period=0 -f ./e2e/manifests/argocd/crd.yaml --ignore-not-found=true
	kubectl delete -n vault-plugin-argocd-tokens-testing --force --grace-period=0 -f ./build.yaml --ignore-not-found=true
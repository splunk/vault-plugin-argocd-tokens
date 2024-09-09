#!/bin/bash
# Source this script in the same session where you want to run vault-config or vault-test

pushd $(dirname "$0")
source ./common.sh

kubectl config set-context ${KUBE_CONTEXT} --namespace ${KUBE_NS}
argocd login "${WFECDSTG_SERVER}"  --core
export WFECDSTG_TOKEN=$(argocd account generate-token -a argocd-tokens-plugin -e 720h)

if [[ "${WFECDSTG_TOKEN:-}" == "" ]]
then
  echo "WFECDSTG_TOKEN environment variable not set" >&2
  exit 1
fi

popd
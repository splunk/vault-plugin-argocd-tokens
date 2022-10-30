#!/bin/bash
export WFECDSTG_SERVER="argocd.wfecd-stg.optimus.prime.us-west-2.splunk8s.io"
export WFECDSTG_ADDR="https://${WFECDSTG_SERVER}"
export VAULT_ADDR="http://127.0.0.1:${VAULT_PORT:-8200}"
export ENGINE_PATH="wfecd-stg"
#!/bin/bash
export WFECDSTG_SERVER="wfecd-stg.prod.svc.splunk8s.io"
export WFECDSTG_ADDR="https://${WFECDSTG_SERVER}"
export KUBE_CONTEXT="cyclops-pdx10-prod-diffie.splunk8s.io"
export KUBE_NS="wfecd-stg"
export VAULT_ADDR="http://127.0.0.1:${VAULT_PORT:-8200}"
export ENGINE_PATH="wfecd-stg"
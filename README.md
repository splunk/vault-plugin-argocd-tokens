# License Notice

Copyright 2021 Splunk, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. 
See the License for the specific language governing permissions and
limitations under the License.

# How to run a dev environment

## Prerequiresites

- argocd cli
- vault cli
- go 1.18

## Setup

- Run `./devel/run-test-vault.sh` in a terminal. Keep this open. This will print useful logs from the vault server
- Run `source ./devel/argocd-login.sh` in a new session. It is important to source this script as the exported variables are needed later.
  - You will likely need to run `okta-kube-token` and create an escalation request like the following first:
```
kubectl escalation -n $KUBE_NS --context $KUBE_CONTEXT request "Vault plugin test" --approve
```
- Run `./devel/vault-config.sh` to configure the vault instance (enable the plugin, write the config etc.)
- Run `./devel/vault-tests.sh` to run the tests

## Dev changes

- You can deactivate, deregister, re-register and reactivate the plugin while making changes in the plugin
- Or simply stop and re-run the `./devel/run-test-vault.sh` to reload vault with the new build of the plugin

## Integration Tests 

- All json files in `e2e/scenarios` are scenarios for integration tests.
- The file structure is as follows 
```json
{
  "scenario": "String scenario outline",
  "command": "Vault CLI command to test scenario, must include -format=json flag",
  "assert": { 
    # Block of values to assert in response. e.g.
    "A": "X", # will assert that the key '.A' takes the value "X" and
    "B": {
      "C": "Y" # will assert that '.B.C' takes the value "Y"
    }
  }
}
```
- To add a new test, simply create a new json in `e2e/scenarios/account` or `e2e/scenarios/project`.
- If configuration of the test Argocd instance is required, add it to `/e2e/manifests/patch` or `/e2e/manifests/resources` as required.
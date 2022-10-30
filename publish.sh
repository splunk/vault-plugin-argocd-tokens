#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [[ "${ARTIFACTORY_PASSWORD}" == "" ]]
then
    echo ARTIFACTORY_PASSWORD not set >&2
    exit 1
fi

if [[ "${ARTIFACTORY_USERNAME}" == "" ]]
then
    echo ARTIFACTORY_USERNAME not set >&2
    exit 1
fi

VERSION=$1
ROOT_URL="https://repo.splunk.com/artifactory"
REPO=splunk8s
BASE_URL="${ROOT_URL}/${REPO}"

EXES="vault-plugin-argocd-tokens"
PUBLISH_RESULTS=/tmp/publish.log


echo publish ${EXES} version ${VERSION}
curl -u "${ARTIFACTORY_USERNAME}:${ARTIFACTORY_PASSWORD}" -X PUT -sSi "${BASE_URL}/vault-plugins/${EXES}/${VERSION}/${EXES}" -T "${GOPATH}/bin/${EXES}" >>${PUBLISH_RESULTS} 2>&1

echo Publish results
cat ${PUBLISH_RESULTS}
echo

curl --head -i "${BASE_URL}/vault-plugins/${EXES}-plugin/${VERSION}/${EXES}"
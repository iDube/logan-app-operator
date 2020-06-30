#!/usr/bin/env bash
# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u

export REPO="logancloud/logan-app-operator"
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
docker tag ${REPO}:latest ${REPO}:${TAG}
docker push ${REPO}:${TAG}

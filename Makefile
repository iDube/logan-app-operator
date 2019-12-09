# Image URL to use all building/pushing image targets
IMG ?= logancloud/logan-app-operator:latest

all: test

# Run tests
test:
	ginkgo -r pkg/

dingding:
	bash ./scripts/dingding.sh

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	operator-sdk up local --namespace=logan --operator-flags "--config=configs/config_local.yaml --zap-devel --zap-level info"

rundebug: fmt vet
	operator-sdk up local --namespace=logan --operator-flags "--config=configs/config_local.yaml --zap-devel"

rundev:
	LOGAN_ENV=dev WATCH_NAMESPACE=logan-dev operator-sdk up local --namespace=logan-dev --operator-flags "--config=configs/config_local.yaml --zap-devel"

runauto:
	LOGAN_ENV=auto WATCH_NAMESPACE=logan-auto operator-sdk up local --namespace=logan-auto --operator-flags "--config=configs/config_local.yaml --zap-devel"

# Install CRDs into a cluster
install:
	kubectl apply -f deploy

# Install webhook into a cluster
initwebhook: initwebhook-test initwebhook-dev initwebhook-auto

initwebhook-test:
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook --namespace logan --secret logan-app-operator-webhook
	cat deploy/webhook.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -

initwebhook-dev:
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook-dev --namespace logan --secret logan-app-operator-webhook-dev
	cat deploy/webhook-dev.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -

initwebhook-auto:
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook-auto --namespace logan --secret logan-app-operator-webhook-auto
	cat deploy/webhook-auto.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -

# Re Install webhook into a cluster
rewebhook:
	oc delete -f deploy/webhook.yaml --ignore-not-found=true
	oc delete secret logan-app-operator-webhook --ignore-not-found=true
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook --namespace logan --secret logan-app-operator-webhook
	cat deploy/webhook.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -
	oc delete -f deploy/webhook-dev.yaml --ignore-not-found=true
	oc delete secret logan-app-operator-webhook-dev --ignore-not-found=true
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook-dev --namespace logan --secret logan-app-operator-webhook-dev
	cat deploy/webhook-dev.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -
	oc delete -f deploy/webhook-auto.yaml --ignore-not-found=true
	oc delete secret logan-app-operator-webhook-auto --ignore-not-found=true
	scripts/webhook-create-signed-cert.sh --service logan-app-webhook-auto --namespace logan --secret logan-app-operator-webhook-auto
	cat deploy/webhook-auto.yaml | scripts/webhook-patch-ca-bundle.sh | kubectl create -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	kubectl apply -f deploy/crds

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/... ./test/...

# Run go vet against code
vet:
	export GO111MODULE=on && go vet ./pkg/... ./cmd/... ./test/...

# Run generate k8s
gen-k8s:
	operator-sdk generate k8s

# Build
build: docker-build docker-push

# Build the docker image
docker-build:
	export GO111MODULE=on && operator-sdk build ${IMG}

travis-docker-build:
	bash ./scripts/travis-build.sh ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

travis-build:
	bash ./scripts/travis-push-docker-image.sh

test-e2e:
	bash ./scripts/travis-e2e.sh

test-e2e-local: docker-build
	bash ./scripts/travis-e2e.sh local

# Init Operator
initdeploy: initcm initrole initcrd
	oc create -f deploy/operator-test.yaml -f deploy/operator-dev.yaml -n logan -f deploy/operator-auto.yaml -n logan

initcm:
	oc create configmap logan-app-operator-config --from-file=configs/config.yaml
	oc create configmap logan-app-operator-config-auto --from-file=configs/config.yaml
	oc create configmap logan-app-operator-config-dev --from-file=configs/config.yaml

initrole:
	oc apply -f deploy/role.yaml
	oc apply -f deploy/role_binding.yaml
	oc apply -f deploy/role_operator.yaml
	oc apply -f deploy/service_account.yaml

initcrd:
	oc apply -f deploy/crds/app.logancloud.com_javaboots_crd.yaml
	oc apply -f deploy/crds/app.logancloud.com_phpboots_crd.yaml
	oc apply -f deploy/crds/app.logancloud.com_pythonboots_crd.yaml
	oc apply -f deploy/crds/app.logancloud.com_nodejsboots_crd.yaml
	oc apply -f deploy/crds/app.logancloud.com_webboots_crd.yaml
	oc apply -f deploy/crds/app.logancloud.com_bootrevisions_crd.yaml

# Redeploy Operator
redeploy: recm rerole recrd
	oc replace -f deploy/operator-test.yaml -f deploy/operator-dev.yaml -f deploy/operator-auto.yaml -n logan

recm:
	oc delete configmap logan-app-operator-config --ignore-not-found=true
	oc create configmap logan-app-operator-config --from-file=configs/config.yaml
	oc delete configmap logan-app-operator-config-dev --ignore-not-found=true
	oc create configmap logan-app-operator-config-dev --from-file=configs/config.yaml
	oc delete configmap logan-app-operator-config-auto --ignore-not-found=true
	oc create configmap logan-app-operator-config-auto --from-file=configs/config.yaml

rerole:
	oc replace -f deploy/role.yaml
	oc replace -f deploy/role_binding.yaml
	oc replace -f deploy/role_operator.yaml
	oc replace -f deploy/service_account.yaml

recrd:
	oc replace -f deploy/crds/app.logancloud.com_javaboots_crd.yaml
	oc replace -f deploy/crds/app.logancloud.com_phpboots_crd.yaml
	oc replace -f deploy/crds/app.logancloud.com_pythonboots_crd.yaml
	oc replace -f deploy/crds/app.logancloud.com_nodejsboots_crd.yaml
	oc replace -f deploy/crds/app.logancloud.com_webboots_crd.yaml
	oc replace -f deploy/crds/app.logancloud.com_bootrevisions_crd.yaml

# test java
test-java:
	oc delete -f examples/test-java.yaml --ignore-not-found=true
	oc create -f examples/test-java.yaml

# test php
test-php:
	oc delete -f examples/test-php.yaml --ignore-not-found=true
	oc create -f examples/test-php.yaml

# test python
test-python:
	oc delete -f examples/test-python.yaml --ignore-not-found=true
	oc create -f examples/test-python.yaml

# test nodejs
test-nodejs:
	oc delete -f examples/test-nodejs.yaml --ignore-not-found=true
	oc create -f examples/test-nodejs.yaml

# test web
test-web:
	oc delete -f examples/test-web.yaml --ignore-not-found=true
	oc create -f examples/test-web.yaml

test-all: test-java test-php test-python test-nodejs test-web

test-deleteall:
	oc delete -f examples/test-java.yaml --ignore-not-found=true
	oc delete -f examples/test-php.yaml --ignore-not-found=true
	oc delete -f examples/test-python.yaml --ignore-not-found=true
	oc delete -f examples/test-nodejs.yaml --ignore-not-found=true
	oc delete -f examples/test-web.yaml --ignore-not-found=true

test-createall:
	oc create -f examples/crds/test_java.yaml
	oc create -f examples/crds/test_php.yaml
	oc create -f examples/crds/test_python.yaml
	oc create -f examples/crds/test_node.yaml
	oc create -f examples/crds/test_web.yaml

#  test recreate 100 times
test-batch:
	scripts/all.sh
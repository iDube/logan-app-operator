# Pre-requisites
Before you get started using logan-app-operator, you need to have a running Kubernetes cluster setup. We use Minikube for testing purposes.

Setup a Kubernetes Cluster in version 1.11 or later
Install kubectl in version 1.11 or later

# Install

```bash
# clone project
git clone https://github.com/logancloud/logan-app-operator.git
cd logan-app-operator

# create logan namespace if not exists
kubectl create namespace logan

# init
make initwebhook
make initdeploy

# start logan-app-operator
kubectl scale deploy logan-app-operator -n logan --replicas=1

# check status
kubectl get pods -n logan
```

# Quick Start

After install the logan-app-operator, you can deploy a javaboot

```bash
# deploy a java boot
make test-java

# check status
kubectl get java -n logan
kubectl get pods -n logan
```

# Kubernetes Custom Resource Definition

## What is this?

An example of a custom Kubernetes controller that's only purpose is to watch for the creation, updating, or deletion of all custom resource of type `Health` (in the all namespaces). It switch on/off `/health` **PUT** method.

## Running

```
$ git clone https://github.com/hsiaoairplane/k8s-crd
$ cd k8s-crd
$ GO111MODULE=on go mod vendor
$ GO111MODULE=on go mod verify
$ make all
```

# Compile
1. Run code generate
2. Build from source code

# Test
1. kubectl apply -f deploy/templates/crd.yaml
2. kubectl apply -f deploy/charts/health-off.yaml
3. kubectl apply -f deploy/charts/health-on.yaml

# Deployment
1. helm install --name k8s-crd deployment
2. helm delete --purge k8s-crd

# Port-forward
1. export POD_NAME=$(kubectl get pods --namespace default -l "app=k8s-crd,release=k8s-crd" -o jsonpath="{.items[0].metadata.name}")
2. kubectl port-forward $POD_NAME 8080:80
3. curl -XGET http://127.0.0.1:8080/health
4. curl -XPUT http://127.0.0.1:8080/health

# Draft
1. draft init
2. draft up
3. draft connect

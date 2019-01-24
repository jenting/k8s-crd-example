# Kubernetes Custom Resource Definition

## What is this?

An example of a custom Kubernetes controller that with HTTP router endpoint `/health`, by default `GET` method supports. But the HTTP method could be changed via Custom Resource Definition to enable/disable other HTTP method on endpoint `/health`.

## Development environment

### Golang

Ensure you got Go 1.11 installed with `go mod` support.

```sh
brew install go
export GO111MODULE=on
go mod vendor
```

### Docker

Docker is required. You may download and install [the installation package](https://store.docker.com/editions/community/docker-ce-desktop-mac), or install it via [Homebrew Cask](https://brew.sh/).

```sh
brew cask install docker
```

### Kubernetes cluster

kubernetes is required. You may setup a K8s cluster by [docker-for-desktop](https://codefresh.io/kubernetes-tutorial/local-kubernetes-mac-minikube-vs-docker-desktop/) or [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

### Kubectl

kubectl is required. You could install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/), or install it via [Homebrew](https://brew.sh/).

```sh
brew install kubernetes-cli
```

### Helm

Helm is required. You may download and install [the binary releases](https://github.com/helm/helm/releases), or install it via [Homebrew](https://brew.sh/).

```sh
brew install kubernetes-helm
```

Tiller is required. You can install it via command.

```sh
kubectl -n kube-system create sa tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller
```

### Draft

Draft is required. You may download and install [the binary releases](https://github.com/Azure/draft/releases), or install it via [Homebrew](https://brew.sh/).

```sh
brew tap azure/draft && brew install draft
```

## Deploy container and helm chart via Draft

1. Set up draft (after Development environment are prepared)

```sh
draft init
```

2. To deploy the application to a Kubernetes dev sandbox, accessible using draft connect over a secured tunnel

```sh
draft up
```

3. Show information

```sh
helm list
kubectl get all
```

## Test CRD behavior

1. Apply CRD deployment

```sh
kubectl apply -f patch/deploy.yaml
```

2. Test default CRD behavior (GET enabled; PUT disabled)

```sh
curl -XGET -i 'localhost:8888/health'
curl -XPUT -i 'localhost:8888/health'
```

3. Enable PUT method

```sh
kubectl replace -f patch/put-on.yaml
```

4. Test CRD behavior (GET and enabled)

```sh
curl -XGET -i 'localhost:8888/health'
curl -XPUT -i 'localhost:8888/health'
```

5. Disable PUT method

```sh
kubectl replace -f patch/put-off.yaml
```

6. Test CRD behavior (GET enabled; PUT disabled)

```sh
curl -XGET -i 'localhost:8888/health'
curl -XPUT -i 'localhost:8888/health'
```

7. You could try other HTTP method on/off behavior with `patch/xxx-on.yaml` or `patch/xxx-off.yaml`

## Clean k8s enviroment

```sh
draft delete
```

#!/usr/bin/make -f

all: codegen
	go build

clean:
	go clean

codegen:
	go get -u k8s.io/code-generator/...
	${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "github.com/hsiaoairplane/k8s-crd/pkg/client" "github.com/hsiaoairplane/k8s-crd/pkg/apis" "health:v1"

.PHONY: all clean codegen

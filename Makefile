#!/usr/bin/make -f

all:
	go build

clean:
	go clean

codegen:
	go get -u k8s.io/code-generator/...
	${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "github.com/jenting/k8s-crd-example/pkg/client" "github.com/jenting/k8s-crd-example/pkg/apis" "health:v1"

.PHONY: all clean codegen

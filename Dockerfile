FROM golang:1.11.4-alpine as builder

ENV GO111MODULE=on

RUN apk add --update alpine-sdk
RUN mkdir -p /go/src/github.com/jenting/k8s-crd-example
WORKDIR /go/src/github.com/jenting/k8s-crd-example

COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .

RUN go build -ldflags "-s -w" -o k8s-crd-example -v

FROM alpine:3.8

RUN apk add --update ca-certificates

COPY --from=builder /go/src/github.com/jenting/k8s-crd-example/k8s-crd-example /app/

ENV PATH=/app:$PATH

WORKDIR /app

EXPOSE 80

CMD ["k8s-crd-example"]

FROM golang:1.12.5-stretch

ARG GITHUB_ACCESS_TOKEN
WORKDIR /go/src/github.com/unhookd/unctl

RUN mkdir -p /tmp/helm && cd /tmp/helm && curl https://storage.googleapis.com/kubernetes-helm/helm-v2.14.0-linux-amd64.tar.gz | tar -xvz && mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm-v2.14.0 && rm -rf /tmp/helm
COPY . .

RUN GO111MODULE=on go test -v \
    -ldflags "-X github.com/unhookd/unctl/lookup.EncodedConfigLookup=$(cat zero-trust-test.yaml | base64 -w0)" \
    ./...

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags "-X github.com/unhookd/unctl/lookup.EncodedConfigLookup=$(cat zero-trust.yaml | base64 -w0)" \
    -o ./unhookd-linux-amd64 main.go

FROM alpine:3.9.4
RUN apk --no-cache add curl git make bash

COPY --from=0 /usr/local/bin/helm-v2.14.0 /usr/local/bin/helm-v2.14.0
RUN ln -s /usr/local/bin/helm-v2.14.0 /usr/local/bin/helm
COPY --from=0 /go/src/github.com/unhookd/unctl/unhookd-linux-amd64 /usr/bin/unhookd
RUN helm init --client-only

CMD ["/usr/bin/unhookd"]

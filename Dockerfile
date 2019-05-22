FROM golang:1.11

ARG GITHUB_ACCESS_TOKEN
WORKDIR /go/src/github.com/org/unhookd

RUN mkdir -p /tmp/helm && cd /tmp/helm && curl https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz | tar -xvz && mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm-v2.9.1 && rm -rf /tmp/helm
RUN mkdir -p /tmp/helm && cd /tmp/helm && curl https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz | tar -xvz && mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm-v2.10.0 && rm -rf /tmp/helm
RUN mkdir -p /tmp/helm && cd /tmp/helm && curl https://storage.googleapis.com/kubernetes-helm/helm-v2.11.0-linux-amd64.tar.gz | tar -xvz && mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm-v2.11.0 && rm -rf /tmp/helm
RUN ln -s /usr/local/bin/helm-v2.11.0 /usr/local/bin/helm
COPY . .

RUN GO111MODULE=on go test -v \
    -ldflags "-X github.com/org/unhookd/lookup.EncodedConfigLookup=$(cat zero-trust-test.yaml | base64 -w0)" \
    ./...

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags "-X github.com/org/unhookd/lookup.EncodedConfigLookup=$(cat zero-trust.yaml | base64 -w0)" \
    -o ./unhookd-linux-amd64 main.go

FROM alpine:3.7
ENV HELM_S3_VERSION 0.7.0
RUN apk --no-cache add curl git make bash

COPY --from=0 /usr/local/bin/helm-v2.9.1 /usr/local/bin/helm-v2.9.1
COPY --from=0 /usr/local/bin/helm-v2.10.0 /usr/local/bin/helm-v2.10.0
COPY --from=0 /usr/local/bin/helm-v2.11.0 /usr/local/bin/helm-v2.11.0
RUN ln -s /usr/local/bin/helm-v2.11.0 /usr/local/bin/helm
COPY --from=0 /go/src/github.com/org/unhookd/unhookd-linux-amd64 /usr/bin/unhookd
RUN helm init --client-only
RUN helm plugin install https://github.com/hypnoglow/helm-s3.git --version ${HELM_S3_VERSION}

CMD ["/usr/bin/unhookd"]

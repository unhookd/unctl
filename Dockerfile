FROM golang:1.12.5-stretch

WORKDIR /go/src/github.com/org/unhookd

RUN mkdir -p /tmp/helm && cd /tmp/helm && curl https://storage.googleapis.com/kubernetes-helm/helm-v2.14.0-linux-amd64.tar.gz | tar -xvz && mv /tmp/helm/linux-amd64/helm /usr/local/bin/helm-v2.14.0 && rm -rf /tmp/helm
COPY . .

RUN GO111MODULE=on go test -v ./...

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./unhookd-linux-amd64 main.go

FROM alpine:3.7
RUN apk --no-cache add curl git make bash

RUN ln -s /usr/local/bin/helm-v2.14.0 /usr/local/bin/helm
COPY --from=0 /go/src/github.com/org/unhookd/unhookd-linux-amd64 /usr/bin/unhookd
RUN helm init --client-only

CMD ["/usr/bin/unhookd"]

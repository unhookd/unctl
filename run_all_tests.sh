#!/bin/bash

go test -v -ldflags "-X github.com/unhookd/unctl/config.EncodedConfigLookup=$(cat zero-trust-test.yaml | base64)"  ./...

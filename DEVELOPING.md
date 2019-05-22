# Developing on Unhookd

- [Quickstart](#quick-start)
- [Setting up your machine](#setting-up-your-machine)
- [Running Unhookd](#running-unhookd)
- [Installing to a Cluster](#installing-to-a-cluster)

## Quick Start

### Setting up your machine

#### Ensuring a correct Go version
Install Go on your machine but you should make sure that you're running on Go 1.11

```
go version
```

If not, upgrade it:

```
brew upgrade golang
```

#### Installing Goland
Goland is a Jetbrains IDE for Golang! If you've used Rubymine, you should feel right at home

## Running Unhookd
```
go run -ldflags "-X github.com/org/unhookd/lookup.EncodedConfigLookup=$(cat zero-trust.yaml | base64)" main.go
```

You can also build the project with `go build` and get dependencies with `go get`

## Deploying to a Cluster


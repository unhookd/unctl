# build stage
FROM golang:alpine AS build-env
ENV GO111MODULE=on
RUN apk update && apk add git
ADD . /src
RUN cd /src && go build -o unhookd

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/unhookd /app/
ENTRYPOINT ./unhookd

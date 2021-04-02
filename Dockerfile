# Build
FROM golang:1.16.2 AS build
ENV GO111MODULE=on
WORKDIR /go/src/github.com/healthcheck-exporter
COPY cmd ./cmd
COPY go.mod .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/

# Release
FROM alpine:3.13.3

RUN apk --no-cache add ca-certificates=20191127-r4

WORKDIR /service

ARG RUN_USER=service

RUN adduser -S -D -H -u 1001 -s /sbin/nologin -G root -g $RUN_USER $RUN_USER

COPY --from=build /go/src/github.com/healthcheck-exporter/service .

RUN chgrp -R 0 /service && chmod -R g+rX /service

USER $RUN_USER

CMD ./service
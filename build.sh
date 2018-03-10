#!/bin/sh

set -e

glide install

go vet $(glide novendor)
go test -race -v $(glide novendor)

sed -i '491i\\t"run",    // github.com/lucaslorentz/caddy-run' \
 vendor/github.com/mholt/caddy/caddyhttp/httpserver/plugin.go

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o caddy

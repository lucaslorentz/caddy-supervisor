#!/bin/sh

set -e

glide install

sed -i '464i\\t"run",    // github.com/lucaslorentz/caddy-run' \
 vendor/github.com/mholt/caddy/caddyhttp/httpserver/plugin.go

go build -o caddy

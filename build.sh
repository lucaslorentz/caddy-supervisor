#!/bin/sh

set -e

glide install

go vet $(glide novendor)
go test -race -v $(glide novendor)

CGO_ENABLED=0 go build -o caddy

echo ""
echo ""
echo ""
echo ==Starting caddy with servertype==
./caddy -type supervisor -conf ./examples/Supervisorfile -log stdout &
CADDY_PID=$!
sleep 5
kill $CADDY_PID
wait $CADDY_PID || true
echo ==Killed caddy with servertype==

sleep 5

echo ""
echo ""
echo ""
echo ==Starting caddy with httpplugin==
./caddy -conf ./examples/Caddyfile -log stdout &
CADDY_PID=$!
sleep 5
kill $CADDY_PID
wait $CADDY_PID || true
echo ==Killed caddy with httpplugin==

sleep 5

echo "Success!"
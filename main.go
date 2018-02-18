package main

import (
	// Caddy
	"github.com/mholt/caddy/caddy/caddymain"

	// Plugins
	_ "github.com/lucaslorentz/caddy-run/plugin"
	_ "github.com/lucaslorentz/caddy-service"
)

func main() {
	caddymain.Run()
}

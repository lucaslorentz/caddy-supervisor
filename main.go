package main

import (
	// Caddy
	"github.com/mholt/caddy/caddy/caddymain"

	// Plugins
	_ "github.com/lucaslorentz/caddy-service"
	_ "github.com/lucaslorentz/caddy-supervisor/httpplugin"
	_ "github.com/lucaslorentz/caddy-supervisor/servertype"
)

func main() {
	caddymain.Run()
}

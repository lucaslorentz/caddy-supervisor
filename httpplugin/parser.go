package httpplugin

import (
	"github.com/caddyserver/caddy"
	"github.com/lucaslorentz/caddy-supervisor/supervisor"
)

func parseHTTPDirectives(c *caddy.Controller) ([]*supervisor.Options, error) {
	var optionsList []*supervisor.Options

	for c.Next() {
		options, err := parseHTTPDirective(c)
		if err != nil {
			return optionsList, err
		}
		optionsList = append(optionsList, options)
	}

	return optionsList, nil
}

func parseHTTPDirective(c *caddy.Controller) (*supervisor.Options, error) {
	var options = supervisor.CreateOptions()

	args := c.RemainingArgs()
	if len(args) > 0 {
		options.Command = args[0]
		options.Args = args[1:]
	}

	for c.NextBlock() {
		supervisor.ParseOption(c, options)
	}

	return options, nil
}

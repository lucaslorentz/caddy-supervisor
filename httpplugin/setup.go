package httpplugin

import (
	"github.com/caddyserver/caddy"
	"github.com/lucaslorentz/caddy-supervisor/supervisor"
)

func init() {
	caddy.RegisterPlugin("supervisor", caddy.Plugin{
		ServerType: "http",
		Action:     setupDirective,
	})
}

func setupDirective(c *caddy.Controller) error {
	return c.OncePerServerBlock(func() error {
		optionsList, err := parseHTTPDirectives(c)
		if err != nil {
			return err
		}

		for _, options := range optionsList {
			supervisors := supervisor.CreateSupervisors(options)
			for _, sup := range supervisors {
				func(s *supervisor.Supervisor) {
					c.OnStartup(func() error {
						go s.Run()
						return nil
					})
					// Use OnRestart to shutdown supervisors before new instances starts
					c.OnRestart(func() error {
						s.Stop()
						return nil
					})
					c.OnFinalShutdown(func() error {
						s.Stop()
						return nil
					})
				}(sup)
			}
		}
		return nil
	})
}

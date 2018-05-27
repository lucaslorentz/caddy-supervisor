package httpplugin

import (
	"sync"

	"github.com/lucaslorentz/caddy-supervisor/supervisor"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("supervisor", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

var supervisors []*supervisor.Supervisor

func setup(c *caddy.Controller) error {
	setupEventsOnlyOnce(c)

	return c.OncePerServerBlock(func() error {
		optionsList, err := parseHTTPDirectives(c)
		if err != nil {
			return err
		}

		for _, options := range optionsList {
			newSupervisors := supervisor.CreateSupervisors(options)
			supervisors = append(supervisors, newSupervisors...)
			for _, supervisor := range newSupervisors {
				go supervisor.Start()
			}
		}
		return nil
	})
}

var didSetupEvents = false

func setupEventsOnlyOnce(c *caddy.Controller) {
	if didSetupEvents {
		return
	}
	c.OnShutdown(shutdownExecutions)
	didSetupEvents = true
}

func shutdownExecutions() error {
	var wg sync.WaitGroup

	for _, s := range supervisors {
		wg.Add(1)
		go func(s *supervisor.Supervisor) {
			defer wg.Done()
			s.Stop()
		}(s)
	}

	wg.Wait()

	return nil
}

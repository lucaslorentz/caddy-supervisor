package plugin

import (
	"log"
	"sync"

	"github.com/mholt/caddy"
)

func init() {
	log.Println("Registering plugin")

	caddy.RegisterPlugin("run", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

var executors []*executor

// setup used internally by Caddy to set up this middleware
func setup(c *caddy.Controller) error {
	setupEventsOnlyOnce(c)

	return c.OncePerServerBlock(func() error {
		optionsList, err := parseOptionsList(c)
		if err != nil {
			return err
		}

		for _, options := range optionsList {
			executor := createExecutor(options)
			executors = append(executors, executor)
			go executor.run()
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

	for _, e := range executors {
		wg.Add(1)
		go func(e *executor) {
			defer wg.Done()
			e.cancel()
		}(e)
	}

	wg.Wait()

	return nil
}

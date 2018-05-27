package servertype

import (
	"strings"
	"sync"

	"github.com/lucaslorentz/caddy-supervisor/supervisor"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

func init() {
	var directives = []string{
		"command",
		"args",
		"dir",
		"redirect_stdout",
		"redirect_stderr",
		"restart_policy",
		"termination_grace_period",
		"env",
		"replicas",
	}

	caddy.RegisterServerType("supervisor", caddy.ServerType{
		Directives: func() []string {
			return directives
		},
		NewContext: newContext,
	})

	for _, directive := range directives {
		caddy.RegisterPlugin(directive, caddy.Plugin{
			ServerType: "supervisor",
			Action:     setup,
		})
	}
}

func newContext(inst *caddy.Instance) caddy.Context {
	return &supervisorContext{
		instance: inst,
		options:  make(map[string]*supervisor.Options),
	}
}

type supervisorContext struct {
	instance *caddy.Instance
	options  map[string]*supervisor.Options
}

func (n *supervisorContext) InspectServerBlocks(sourceFile string, serverBlocks []caddyfile.ServerBlock) ([]caddyfile.ServerBlock, error) {
	for _, sb := range serverBlocks {
		key := mergeKeys(sb.Keys)

		n.options[key] = supervisor.CreateOptions()
	}

	return serverBlocks, nil
}

var supervisors []*supervisor.Supervisor

// MakeServers uses the newly-created configs to create and return a list of server instances.
func (n *supervisorContext) MakeServers() ([]caddy.Server, error) {
	for _, options := range n.options {
		newSupervisors := supervisor.CreateSupervisors(options)
		supervisors = append(supervisors, newSupervisors...)
		for _, supervisor := range newSupervisors {
			go supervisor.Start()
		}
	}
	return nil, nil
}

func setup(c *caddy.Controller) error {
	key := mergeKeys(c.ServerBlockKeys)

	ctx := c.Context().(*supervisorContext)

	setupEventsOnlyOnce(c)

	return c.OncePerServerBlock(func() error {
		options := ctx.options[key]
		for c.Next() {
			supervisor.ParseOption(c, options)
		}
		return nil
	})
}

func mergeKeys(keys []string) string {
	return strings.Join(keys, " ")
}

var didSetupEvents = false

func setupEventsOnlyOnce(c *caddy.Controller) {
	if didSetupEvents {
		return
	}
	c.OnShutdown(shutdownSupervisors)
	didSetupEvents = true
}

func shutdownSupervisors() error {
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

package servertype

import (
	"net"
	"strings"
	"sync"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyfile"
	"github.com/lucaslorentz/caddy-supervisor/supervisor"
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
			Action:     setupDirective,
		})
	}
}

func newContext(inst *caddy.Instance) caddy.Context {
	context := &supervisorContext{
		instance:    inst,
		options:     make(map[string]*supervisor.Options),
		supervisors: []*supervisor.Supervisor{},
	}

	inst.OnShutdown = append(inst.OnShutdown, func() error {
		shutdownSupervisors(context.supervisors)
		return nil
	})

	return context
}

type supervisorContext struct {
	instance    *caddy.Instance
	options     map[string]*supervisor.Options
	supervisors []*supervisor.Supervisor
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
	servers := []caddy.Server{}
	for _, options := range n.options {
		newSupervisors := supervisor.CreateSupervisors(options)
		supervisors = append(supervisors, newSupervisors...)
		for _, supervisor := range newSupervisors {
			n.supervisors = append(n.supervisors, supervisor)
			servers = append(servers, &supervisorServer{supervisor: supervisor})
		}
	}
	return servers, nil
}

type supervisorServer struct {
	supervisor *supervisor.Supervisor
}

func (server *supervisorServer) Listen() (net.Listener, error) {
	return nil, nil
}
func (server *supervisorServer) Serve(net.Listener) error {
	server.supervisor.Run()
	return nil
}
func (server *supervisorServer) ListenPacket() (net.PacketConn, error) {
	return nil, nil
}
func (server *supervisorServer) ServePacket(net.PacketConn) error {
	return nil
}

func setupDirective(c *caddy.Controller) error {
	key := mergeKeys(c.ServerBlockKeys)

	ctx := c.Context().(*supervisorContext)

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

func shutdownSupervisors(supervisors []*supervisor.Supervisor) error {
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

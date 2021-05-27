package supervisor

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

// Interface guards
var (
	_ caddy.App         = (*App)(nil)
	_ caddy.Module      = (*App)(nil)
	_ caddy.Provisioner = (*App)(nil)
)

func init() {
	caddy.RegisterModule(App{})
}

type App struct {
	Supervise   []Definition `json:"supervise,omitempty"`
	log         *zap.Logger
	supervisors []*Supervisor
}

// CaddyModule implements caddy.Module
func (a App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "supervisor",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision implements caddy.Provisioner
func (a *App) Provision(context caddy.Context) error {
	a.log = context.Logger(a)

	for _, definition := range a.Supervise {

		supervisors, err := definition.ToSupervisors(a.log)

		if err != nil {
			return err
		}

		a.supervisors = append(a.supervisors, supervisors...)
	}

	a.log.Debug("module provisioned", zap.Any("supervisors", a.supervisors))

	return nil
}

// Start implements caddy.App
func (a *App) Start() error {
	for _, s := range a.supervisors {
		go s.Run()
	}

	a.log.Debug("module started")

	return nil
}

// Stop implements caddy.App
func (a *App) Stop() error {
	for _, s := range a.supervisors {
		s.Stop()
	}

	return nil
}

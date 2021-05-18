package supervisor

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Definition is the configuration for process to supervise
type Definition struct {
	// Command to start and supervise. First item is the program to start, others are arguments.
	// Supports template.
	Command []string `json:"command"`
	// Replicas control how many instances of Command should run.
	Replicas int `json:"replicas,omitempty"`
	// Dir defines the working directory the command should be executed in.
	// Supports template.
	// Default: current working dir
	Dir string `json:"dir,omitempty"`
	// Env declares environment variables that should be passed to command.
	// Supports template.
	Env map[string]string `json:"env,omitempty"`
	// RedirectStdout is the file where Command stdout is written. Use "stdout" to redirect to caddy stdout.
	RedirectStdout string `json:"redirect_stdout,omitempty"`
	// RedirectStderr is the file where Command stderr is written. Use "stderr" to redirect to caddy stderr.
	RedirectStderr string `json:"redirect_stderr,omitempty"`
	// RestartPolicy define under which conditions the command should be restarted after exit.
	// Valid values:
	//  - **never**: do not restart the command
	//  - **on_failure**: restart if exit code is not 0
	//  - **always**: always restart
	RestartPolicy RestartPolicy `json:"restart_policy,omitempty"`
	// TerminationGracePeriod defines the amount of time to wait for Command graceful termination before killing it. Ex: 10s
	TerminationGracePeriod string `json:"termination_grace_period,omitempty"`
}

// ToSupervisors creates supervisors from the Definition (one per replica) and applies templates where needed
func (d Definition) ToSupervisors(logger *zap.Logger) ([]*Supervisor, error) {
	var supervisors []*Supervisor

	opts := &Options{
		Command:        d.Command[0],
		Args:           d.Command[1:],
		Dir:            d.Dir,
		Env:            d.envToCmdArg(),
		RedirectStdout: d.RedirectStdout,
		RedirectStderr: d.RedirectStderr,
		RestartPolicy:  d.RestartPolicy,
	}

	replicas := d.Replicas

	if replicas == 0 {
		replicas = 1
	}

	if opts.RestartPolicy == "" {
		opts.RestartPolicy = RestartAlways
	}

	if opts.RedirectStdout == "" {
		opts.RedirectStdout = "stdout"
	}

	if opts.RedirectStderr == "" {
		opts.RedirectStderr = "stderr"
	}

	if d.TerminationGracePeriod == "" {
		opts.TerminationGracePeriod = 10 * time.Second
	} else {
		var err error
		opts.TerminationGracePeriod, err = time.ParseDuration(d.TerminationGracePeriod)

		if err != nil {
			return supervisors, fmt.Errorf("cannot parse termination grace period of supervisor '%s'", strings.Join(d.Command, " "))
		}
	}

	for i := 0; i < replicas; i++ {
		opts.Replica = i

		templatedOpts, err := opts.processTemplates()

		if err != nil {
			return supervisors, err
		}

		supervisor := &Supervisor{
			Options: templatedOpts,
			logger: logger.
				With(zap.Strings("command", d.Command)).
				With(zap.Int("replica", templatedOpts.Replica)),
		}

		supervisors = append(supervisors, supervisor)
	}

	return supervisors, nil
}

func (d Definition) envToCmdArg() []string {
	env := make([]string, len(d.Env))
	i := 0

	for key, value := range d.Env {
		env[i] = fmt.Sprintf("%s=%s", key, value)
		i++
	}

	return env
}

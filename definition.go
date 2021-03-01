package supervisor

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Definition is the configuration for process to supervise
type Definition struct {
	Command                []string `json:"command,omitempty"`
	Replicas               int
	KeepRunning            bool `json:"keep_running"`
	Dir                    string
	Env                    []string
	RedirectStdout         string        `json:"redirect_stdout"`
	RedirectStderr         string        `json:"redirect_stderr"`
	RestartPolicy          RestartPolicy `json:"restart_policy"`
	TerminationGracePeriod string        `json:"termination_grace_period"`
}

// ToSupervisors creates supervisors from the Definition (one per replica) and applies templates where needed
func (d Definition) ToSupervisors(logger *zap.Logger) ([]*Supervisor, error) {
	var supervisors []*Supervisor

	cmd := strings.Join(d.Command, " ")

	opts := &Options{
		Command:                d.Command[0],
		Args:                   d.Command[1:],
		Dir:                    d.Dir,
		Env:                    d.Env,
		RedirectStdout:         d.RedirectStdout,
		RedirectStderr:         d.RedirectStderr,
		RestartPolicy:          d.RestartPolicy,
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
				With(zap.String("command", cmd)).
				With(zap.Int("replica", templatedOpts.Replica)),
		}

		supervisors = append(supervisors, supervisor)
	}

	return supervisors, nil
}
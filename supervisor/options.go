package supervisor

import (
	"time"
)

// Options exposes settings to create a process supervisor
type Options struct {
	Command                string
	Args                   []string
	Dir                    string
	Env                    []string
	RedirectStdout         string
	RedirectStderr         string
	RestartPolicy          RestartPolicy
	TerminationGracePeriod time.Duration
}

// CreateOptions createnew SupervisorOptions with default settings
func CreateOptions() *Options {
	return &Options{
		TerminationGracePeriod: 10 * time.Second,
	}
}

// RestartPolicy determines when a supervised process should be restarted
type RestartPolicy string

const (
	// RestartNever indicates to never restart the process
	RestartNever = RestartPolicy("never")
	// RestartOnFailure indicates to only restart the process after failures
	RestartOnFailure = RestartPolicy("on_failure")
	// RestartAlways indicates to always restart the process
	RestartAlways = RestartPolicy("always")
)

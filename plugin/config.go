package plugin

import (
	"time"
)

type runOptions struct {
	command                string
	args                   []string
	dir                    string
	env                    []string
	redirectStdout         string
	redirectStderr         string
	restartPolicy          restartPolicy
	terminationGracePeriod time.Duration
}

func createRunOptions() *runOptions {
	return &runOptions{
		terminationGracePeriod: 10 * time.Second,
	}
}

type restartPolicy string

const (
	restartNever     = restartPolicy("never")
	restartOnFailure = restartPolicy("on_failure")
	restartAlways    = restartPolicy("always")
)

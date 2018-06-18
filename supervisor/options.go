package supervisor

import (
	"bytes"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
)

// Options exposes settings to create a process supervisor
type Options struct {
	Replicas               int
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
		Replicas:               1,
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

func (options *Options) processTemplates(data *TemplateData) *Options {
	result := Options{}
	result.Command = processTemplates(options.Command, data)

	result.Args = make([]string, len(options.Args))
	for i, arg := range options.Args {
		result.Args[i] = processTemplates(arg, data)
	}

	result.Dir = processTemplates(options.Dir, data)

	result.Env = make([]string, len(options.Env))
	for i, env := range options.Env {
		result.Env[i] = processTemplates(env, data)
	}

	result.RedirectStdout = options.RedirectStdout
	result.RedirectStderr = options.RedirectStderr
	result.RestartPolicy = options.RestartPolicy
	result.TerminationGracePeriod = options.TerminationGracePeriod

	return &result
}

// TemplateData contains data to be accessed from templates
type TemplateData struct {
	Replica int
}

func processTemplates(text string, data interface{}) string {
	tmpl, err := template.New("test").Funcs(sprig.TxtFuncMap()).Parse(text)
	if err != nil {
		return err.Error()
	}

	var writer bytes.Buffer
	err = tmpl.Execute(&writer, data)
	if err != nil {
		return err.Error()
	}
	return writer.String()
}

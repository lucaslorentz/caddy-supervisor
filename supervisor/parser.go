package supervisor

import (
	"log"
	"strconv"
	"time"

	"github.com/caddyserver/caddy"
)

// ParseOption parses supervisor options
func ParseOption(c *caddy.Controller, options *Options) {
	v := c.Val()
	switch v {
	case "command":
		args := c.RemainingArgs()
		if len(args) == 1 {
			options.Command = args[0]
		} else {
			log.Printf("Option 'command' expects 1 argument\n")
		}
		break
	case "args":
		options.Args = c.RemainingArgs()
		break
	case "dir":
		args := c.RemainingArgs()
		if len(args) == 1 {
			options.Dir = args[0]
		} else {
			log.Printf("Option 'dir' expects 1 argument\n")
		}
		break
	case "redirect_stdout":
		if c.NextArg() {
			options.RedirectStdout = c.Val()
		} else {
			options.RedirectStdout = "stdout"
		}
		break
	case "redirect_stderr":
		if c.NextArg() {
			options.RedirectStderr = c.Val()
		} else {
			options.RedirectStderr = "stderr"
		}
		break
	case "restart_policy":
		args := c.RemainingArgs()
		if len(args) == 1 {
			switch args[0] {
			case "always":
				options.RestartPolicy = RestartAlways
				break
			case "on_failure":
				options.RestartPolicy = RestartOnFailure
				break
			case "never":
				options.RestartPolicy = RestartNever
				break
			default:
				options.RestartPolicy = RestartNever
				log.Printf("Invalid 'restart' option %v\n", options.RestartPolicy)
				break
			}
		} else {
			log.Printf("Option 'restart' expects 1 argument\n")
		}
	case "termination_grace_period":
		args := c.RemainingArgs()
		if len(args) == 1 {
			period, err := time.ParseDuration(args[0])
			if err == nil {
				options.TerminationGracePeriod = period
			} else {
				log.Printf("Invalid 'termination_grace_period' value %v\n", args[0])
			}
		} else {
			log.Printf("Option 'termination_grace_period' expects 1 argument\n")
		}
	case "env":
		args := c.RemainingArgs()
		if len(args) == 1 {
			options.Env = append(options.Env, args[0])
		} else {
			log.Printf("Option 'env' expects 1 argument in format KEY=VALUE\n")
		}
		break
	case "replicas":
		args := c.RemainingArgs()
		if len(args) == 1 {
			replicas, err := strconv.Atoi(args[0])
			if err == nil {
				options.Replicas = replicas
			} else {
				log.Printf("Invalid 'replicas' value %v\n", args[0])
			}
		} else {
			log.Printf("Option 'replicas' expects 1 argument\n")
		}
		break
	}
}

package plugin

import (
	"log"
	"time"

	"github.com/mholt/caddy"
)

func parseOptionsList(c *caddy.Controller) ([]runOptions, error) {
	var optionsList []runOptions

	for c.Next() {
		options, err := parseOptions(c)
		if err != nil {
			return optionsList, err
		}
		optionsList = append(optionsList, options)
	}

	return optionsList, nil
}

func parseOptions(c *caddy.Controller) (runOptions, error) {
	var options = createRunOptions()

	args := c.RemainingArgs()
	if len(args) > 0 {
		options.command = args[0]
		options.args = args[1:]
	}

	for c.NextBlock() {
		switch c.Val() {
		case "command":
			args := c.RemainingArgs()
			if len(args) == 1 {
				options.command = args[0]
			} else {
				log.Printf("Option 'command' expects 1 argument\n")
			}
			break
		case "args":
			options.args = c.RemainingArgs()
			break
		case "dir":
			args := c.RemainingArgs()
			if len(args) == 1 {
				options.dir = args[0]
			} else {
				log.Printf("Option 'dir' expects 1 argument\n")
			}
			break
		case "redirect_stdout":
			if c.NextArg() {
				options.redirectStdout = c.Val()
			} else {
				options.redirectStdout = "stdout"
			}
			break
		case "redirect_stderr":
			if c.NextArg() {
				options.redirectStderr = c.Val()
			} else {
				options.redirectStderr = "stderr"
			}
			break
		case "restart_policy":
			args := c.RemainingArgs()
			if len(args) == 1 {
				switch args[0] {
				case "always":
					options.restartPolicy = restartAlways
					break
				case "on_failure":
					options.restartPolicy = restartOnFailure
					break
				case "never":
					options.restartPolicy = restartNever
					break
				default:
					options.restartPolicy = restartNever
					log.Printf("Invalid 'restart' option %v\n", options.restartPolicy)
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
					options.terminationGracePeriod = period
				} else {
					log.Printf("Invalid 'termination_grace_period' value %v\n", args[0])
				}
			} else {
				log.Printf("Option 'termination_grace_period' expects 1 argument\n")
			}
		case "env":
			args := c.RemainingArgs()
			if len(args) == 1 {
				options.env = append(options.env, args[0])
			} else {
				log.Printf("Option 'env' expects 1 argument in format KEY=VALUE\n")
			}
			break
		}
	}

	return options, nil
}

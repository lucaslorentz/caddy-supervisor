package supervisor

import (
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	maxRestartDelay = 5 * time.Minute
	minRestartDelay = 10 * time.Second
)

// Supervisor provides functionality to start and supervise a background process
type Supervisor struct {
	options     *Options
	cmd         *exec.Cmd
	keepRunning bool
}

// CreateSupervisor creates a new process supervisor
func CreateSupervisor(options *Options) *Supervisor {
	return &Supervisor{
		options: options,
	}
}

// Start a process and supervise
func (s *Supervisor) Start() {
	s.keepRunning = true
	go s.supervise()
}

func (s *Supervisor) supervise() {
	restartDelay := minRestartDelay
	durationToResetRestartDelay := 10 * time.Minute

	for s.keepRunning {
		s.cmd = exec.Command(s.options.Command, s.options.Args...)

		s.cmd.Env = append(os.Environ(), s.options.Env...)

		if s.options.Dir != "" {
			s.cmd.Dir = s.options.Dir
		}

		if stdoutFile := getFile(s.options.RedirectStdout); stdoutFile != nil {
			s.cmd.Stdout = stdoutFile
			defer stdoutFile.Close()
		}

		if stderrFile := getFile(s.options.RedirectStderr); stderrFile != nil {
			s.cmd.Stderr = stderrFile
			defer stderrFile.Close()
		}

		start := time.Now()
		err := s.cmd.Run()
		duration := time.Now().Sub(start)

		if err != nil {
			log.Printf("Process error: %v\n", err)
		} else {
			log.Printf("Process exited after: %v\n", duration)
		}

		if !s.keepRunning {
			break
		}

		if duration > durationToResetRestartDelay {
			log.Printf("Resetting restart delay to %v\n", minRestartDelay)
			restartDelay = minRestartDelay
		}

		mustRestart := false

		switch s.options.RestartPolicy {
		case RestartAlways:
			mustRestart = true
			break
		case RestartOnFailure:
			mustRestart = err != nil
			break
		}

		if mustRestart {
			log.Printf("Restarting in %v\n", restartDelay)
			time.Sleep(restartDelay)

			restartDelay = restartDelay * 2
			if restartDelay > maxRestartDelay {
				restartDelay = maxRestartDelay
			}
		} else {
			s.keepRunning = false
		}
	}
}

// Stop the supervised process
func (s *Supervisor) Stop() {
	s.keepRunning = false

	if cmdIsRunning(s.cmd) {
		err := s.cmd.Process.Signal(os.Interrupt)
		if err == nil {
			go func() {
				time.Sleep(s.options.TerminationGracePeriod)
				if cmdIsRunning(s.cmd) {
					s.cmd.Process.Kill()
				}
			}()
			s.cmd.Process.Wait()
		} else {
			s.cmd.Process.Kill()
		}
	}
}

func cmdIsRunning(cmd *exec.Cmd) bool {
	return cmd != nil && cmd.Process != nil && (cmd.ProcessState == nil || !cmd.ProcessState.Exited())
}

func getFile(value string) *os.File {
	if value == "" {
		return nil
	}

	switch value {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		outFile, err := os.OpenFile(value, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil
		}
		return outFile
	}
}

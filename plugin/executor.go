package plugin

import (
	"log"
	"os"
	"os/exec"
	"time"
)

type executor struct {
	run    func()
	cancel func()
}

func createExecutor(options runOptions) executor {
	var cmd *exec.Cmd
	maxRestartDelay := 5 * time.Minute
	minRestartDelay := 10 * time.Second
	restartDelay := minRestartDelay
	durationToResetRestartDelay := 10 * time.Minute

	keepRunning := true

	run := func() {
		for keepRunning {
			cmd = exec.Command(options.command, options.args...)

			cmd.Env = append(os.Environ(), options.env...)

			if options.dir != "" {
				cmd.Dir = options.dir
			}

			if stdoutFile := getFile(options.redirectStdout); stdoutFile != nil {
				cmd.Stdout = stdoutFile
				defer stdoutFile.Close()
			}

			if stderrFile := getFile(options.redirectStderr); stderrFile != nil {
				cmd.Stderr = stderrFile
				defer stderrFile.Close()
			}

			start := time.Now()
			err := cmd.Run()
			duration := time.Now().Sub(start)

			if err != nil {
				log.Printf("Process error: %v\n", err)
			} else {
				log.Printf("Process exited after: %v\n", duration)
			}

			if !keepRunning {
				break
			}

			if duration > durationToResetRestartDelay {
				log.Printf("Resetting restart delay to %v\n", minRestartDelay)
				restartDelay = minRestartDelay
			}

			mustRestart := false

			switch options.restartPolicy {
			case restartAlways:
				mustRestart = true
				break
			case restartOnFailure:
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
				keepRunning = false
			}
		}
	}

	cancel := func() {
		keepRunning = false

		if isRunning(cmd) {
			err := cmd.Process.Signal(os.Interrupt)
			if err == nil {
				go func() {
					time.Sleep(options.terminationGracePeriod)
					if isRunning(cmd) {
						cmd.Process.Kill()
					}
				}()
				cmd.Process.Wait()
			} else {
				cmd.Process.Kill()
			}
		}
	}

	return executor{
		run:    run,
		cancel: cancel,
	}
}

func isRunning(cmd *exec.Cmd) bool {
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

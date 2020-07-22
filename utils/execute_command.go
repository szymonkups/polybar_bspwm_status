package utils

import (
	"bufio"
	"fmt"
	"os/exec"
)

func ExecuteCommand(callback func(), command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)

	// Get a pipe to read from standard out
	r, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("cannot create stdout pipe for command %s", command)
	}

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan []string)

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		var lines []string

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			if callback != nil {
				callback()
			}
			lines = append(lines, line)
		}

		// We're all done, unblock the channel
		done <- lines

	}()

	// Start the command and check for errors
	err = cmd.Start()

	if err != nil {
		return nil, fmt.Errorf("cannot start the command %s", command)
	}

	// Wait for all output to be processed
	output := <-done

	// Wait for the command to finish
	err = cmd.Wait()

	if err != nil {
		return nil, fmt.Errorf("command %s exited with error", command)
	}

	return output, nil
}

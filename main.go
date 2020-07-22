package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type Client struct {
	InstanceName string
	ClassName string
}

type Node struct {
	Id          int
	FirstChild  *Node
	SecondChild *Node
	Client Client
}

type DesktopInfo struct {
	Name string
	Id   int
	Root *Node
	FocusedNodeId int
}

type MonitorInfo struct {
	Name             string
	Id               int
	FocusedDesktopId int
	Desktops         []DesktopInfo
}

func main() {
	monitor, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	monitorId, err := getMonitorId(monitor)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot find monitor %d", monitor))
	}

	_, err = executeCommand(outputFunction(monitorId), "bspc", "subscribe")

	if err != nil {
		log.Fatal(err)
	}
}

func outputFunction(monitorId string) func() {
	return func() {
		monitorInfo, err := getMonitorInfo(monitorId)

		if err != nil {
			log.Fatal(err)
		}

		focusedDesktopId := monitorInfo.FocusedDesktopId
		var focusedDesktop DesktopInfo
		isMonFocused, err := isMonitorFocused(monitorId)

		if err != nil {
			log.Fatal(err)
		}

		output := ""

		// Desktops info
		for _, desktop := range monitorInfo.Desktops {
			isDesktopFocused := desktop.Id == focusedDesktopId

			if isDesktopFocused {
				focusedDesktop = desktop
			}

			character := "\ufc64"
			if desktop.Root != nil {
				character = "\ufc63"
			} else {
				if isDesktopFocused {
					character = "%{T2}\ufb66%{T-}"
				}
			}

			color := "#4C566A"

			if isMonFocused {
				color = "#D8DEE9"

				if isDesktopFocused {
					color = "#EBCB8B"
				}
			}

			output += fmt.Sprintf(" %%{F%s}%s%%{F-} ", color, character)
		}

		// Leafs on current desktop info
		leafs := getLeafNodesOnDesktop(focusedDesktopId)
		output += "  "

		focusedLeafIndex := 0
		focusedLeafId := fmt.Sprintf("0x%08X", focusedDesktop.FocusedNodeId)
		for i, leaf := range leafs {
			if focusedLeafId == leaf {
				focusedLeafIndex = i +1
			}
		}

		color := "#4C566A"

		if isMonFocused {
			color = "#D8DEE9"
		}
		
		output += fmt.Sprintf("%%{F%s}%02d/%02d%%{F-}", color, focusedLeafIndex, len(leafs))

		fmt.Println(output)
	}
}

func getMonitorInfo(monitorId string) (*MonitorInfo, error) {
	output, err := executeCommand(nil, "bspc", "query", "-T", "-m", monitorId)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot get info of monitor %s", monitorId))
	}

	monitorData := []byte(output[0])

	var i MonitorInfo
	err = json.Unmarshal(monitorData, &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func getMonitorId(monitor int) (string, error) {
	monitors, err := executeCommand(nil, "bspc", "query", "-M")

	if err != nil {
		return "", fmt.Errorf("could not execute bspc query -M")
	}

	if monitor > len(monitors) {
		return "", fmt.Errorf("monitor index not found: %d", monitor)
	}

	return monitors[monitor-1], nil
}

func isMonitorFocused(monitorId string) (bool, error) {
	monitors, err := executeCommand(nil, "bspc", "query", "-M", "-m", "focused")

	if err != nil {
		return false, fmt.Errorf("could not execute bspc query -M -m focused")
	}

	if len(monitors) < 1 {
		return false, nil
	}

	return monitors[0] == monitorId, nil
}

func getLeafNodesOnDesktop(desktopId int) []string {
	leafs, err := executeCommand(nil, "bspc", "query", "-N", "-d", fmt.Sprintf("0x%08X", desktopId), "-n", ".leaf")

	if err != nil {
		return []string{}
	}

	return leafs
}

func executeCommand(callback func(), command string, args ...string) ([]string, error) {
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

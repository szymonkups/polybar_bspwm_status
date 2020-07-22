package main

import (
	"encoding/json"
	"fmt"
	"github.com/szymonkups/polybar_bspwm_status/types"
	"github.com/szymonkups/polybar_bspwm_status/utils"
	"log"
	"os"
	"strconv"
)


func main() {
	monitor, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	monitorId, err := getMonitorId(monitor)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot find monitor %d", monitor))
	}

	_, err = utils.ExecuteCommand(outputFunction(monitorId), "bspc", "subscribe")

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
		var focusedDesktop types.DesktopInfo
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

func getMonitorInfo(monitorId string) (*types.MonitorInfo, error) {
	output, err := utils.ExecuteCommand(nil, "bspc", "query", "-T", "-m", monitorId)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot get info of monitor %s", monitorId))
	}

	monitorData := []byte(output[0])

	var i types.MonitorInfo
	err = json.Unmarshal(monitorData, &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func getMonitorId(monitor int) (string, error) {
	monitors, err := utils.ExecuteCommand(nil, "bspc", "query", "-M")

	if err != nil {
		return "", fmt.Errorf("could not execute bspc query -M")
	}

	if monitor > len(monitors) {
		return "", fmt.Errorf("monitor index not found: %d", monitor)
	}

	return monitors[monitor-1], nil
}

func isMonitorFocused(monitorId string) (bool, error) {
	monitors, err := utils.ExecuteCommand(nil, "bspc", "query", "-M", "-m", "focused")

	if err != nil {
		return false, fmt.Errorf("could not execute bspc query -M -m focused")
	}

	if len(monitors) < 1 {
		return false, nil
	}

	return monitors[0] == monitorId, nil
}

func getLeafNodesOnDesktop(desktopId int) []string {
	leafs, err := utils.ExecuteCommand(nil, "bspc", "query", "-N", "-d", fmt.Sprintf("0x%08X", desktopId), "-n", ".leaf")

	if err != nil {
		return []string{}
	}

	return leafs
}



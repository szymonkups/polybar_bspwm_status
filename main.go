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
	logger := log.New(os.Stderr, "", 0)

	if len(os.Args) < 2 {
		logger.Println("Please provide a correct monitor index")
		logger.Println()
		logger.Println("polybar_bspwm_status <monitor_index>")
		os.Exit(1)
	}

	monitorIndex, err := strconv.Atoi(os.Args[1])

	if err != nil {
		logger.Fatal(err)
	}

	_, err = utils.ExecuteCommand(outputFunction(monitorIndex), "bspc", "subscribe")

	if err != nil {
		logger.Fatal(err)
	}
}

func outputFunction(monitorIndex int) func() {
	return func() {
		monitorInfo, err := getMonitorInfo(monitorIndex)

		if err != nil {
			log.Fatal(err)
		}

		focusedDesktopId := monitorInfo.FocusedDesktopId
		var focusedDesktop types.DesktopInfo
		isMonFocused := isMonitorFocused(monitorInfo.Id)

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


func getMonitorInfo(monitorIndex int) (*types.MonitorInfo, error) {
	args := []string{"query", "-T", "-m"}

	if monitorIndex > 0 {
		args = append(args, "-m", fmt.Sprintf("^%d", monitorIndex))
	} else {
		args = append(args, "-m", "focused")
	}

	output, err := utils.ExecuteCommand(nil, "bspc", args...)


	if err != nil {
		log.Fatal(fmt.Errorf("cannot get info of monitor %d", monitorIndex))
	}

	monitorData := []byte(output[0])

	var i types.MonitorInfo
	err = json.Unmarshal(monitorData, &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}


func isMonitorFocused(monitorId int) bool {
	focusedMonitorInfo, err := getMonitorInfo(-1)

	if err != nil {
		log.Fatal(err)
	}

	return focusedMonitorInfo.Id == monitorId
}

func getLeafNodesOnDesktop(desktopId int) []string {
	leafs, err := utils.ExecuteCommand(nil, "bspc", "query", "-N", "-d", fmt.Sprintf("0x%08X", desktopId), "-n", ".leaf")

	if err != nil {
		return []string{}
	}

	return leafs
}



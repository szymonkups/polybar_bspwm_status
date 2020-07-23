package main

import (
	"fmt"
	"github.com/szymonkups/polybar_bspwm_status/bspwm"
	"github.com/szymonkups/polybar_bspwm_status/utils"
	"log"
	"os"
	"strconv"
)


func main() {
	logger := log.New(os.Stderr, "", 0)
	monitorIndex := 1

	if len(os.Args) > 1 {
		var err error
		monitorIndex, err = strconv.Atoi(os.Args[1])

		if err != nil {
			logger.Fatalf("Cannot parse monitor_index: %s", os.Args[1])
		}
	}



	_, err := utils.ExecuteCommand(outputFunction(monitorIndex), "bspc", "subscribe")

	if err != nil {
		logger.Fatal("Cannot execute \"bspc subscribe\"")
	}
}

func outputFunction(monitorIndex int) func() {
	return func() {
		monitorInfo, err := bspwm.GetMonitorInfo(monitorIndex)

		if err != nil {
			log.Fatal(err)
		}

		focusedDesktopId := monitorInfo.FocusedDesktopId
		var focusedDesktop bspwm.DesktopInfo
		isMonFocused := bspwm.IsMonitorFocused(monitorInfo.Id)

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
			}
			
			color := "#4C566A"

			if isMonFocused {
				color = "#D8DEE9"

				if isDesktopFocused {
					color = "#EBCB8B"
				}
			}

			if isMonFocused && isDesktopFocused {
				output += fmt.Sprintf("%%{u%s}%%{+u} %%{F%s}%s%%{F-} %%{-u}%%{u-}", color, color, character)
			} else {
				output += fmt.Sprintf(" %%{F%s}%s%%{F-} ", color, character)
			}

		}

		// Leafs on current desktop info
		leafs := bspwm.GetLeafNodesOnDesktop(focusedDesktopId)
		output += "  "

		focusedLeafIndex := 0
		for i, leaf := range leafs {
			if focusedDesktop.FocusedNodeId == leaf {
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





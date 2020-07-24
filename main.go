package main

import (
	"flag"
	"fmt"
	"github.com/szymonkups/polybar_bspwm_status/bspwm"
	"github.com/szymonkups/polybar_bspwm_status/utils"
	"log"
	"os"
)



func main() {
	logger := log.New(os.Stderr, "", 0)

	config, err := utils.LoadConfig()

	if err != nil {
		logger.Fatal("Cannot load config.json", err)
	}

	flag.IntVar(&config.MonitorIndex, "m", 1, "monitor index")
	flag.Parse()

	_, err = utils.ExecuteCommand(outputFunction(config), "bspc", "subscribe")

	if err != nil {
		logger.Fatal("Cannot execute \"bspc subscribe\"")
	}
}



func outputFunction(config *utils.Config) func() {
	return func() {
		monitorInfo, err := bspwm.GetMonitorInfo(config.MonitorIndex)

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

			isDesktopEmpty := desktop.Root == nil

			if isDesktopEmpty {
				if isMonFocused {
					if isDesktopFocused {
						output += config.Labels.EmptyDesktop.Focused
					} else {
						output += config.Labels.EmptyDesktop.Blurred
					}
				} else {
					output += config.Labels.EmptyDesktop.Inactive
				}
			} else {
				if isMonFocused {
					if isDesktopFocused {
						output += config.Labels.OccupiedDesktop.Focused
					} else {
						output += config.Labels.OccupiedDesktop.Blurred
					}
				} else {
					output += config.Labels.OccupiedDesktop.Inactive
				}
			}


			//color := "#4C566A"
			//
			//if isMonFocused {
			//	color = "#D8DEE9"
			//
			//	if isDesktopFocused {
			//		color = "#EBCB8B"
			//	}
			//}
			//
			//if isMonFocused && isDesktopFocused {
			//	output += fmt.Sprintf("%%{u%s}%%{+u} %%{F%s}%s%%{F-} %%{-u}%%{u-}", color, color, character)
			//} else {
			//	output += fmt.Sprintf(" %%{F%s}%s%%{F-} ", color, character)
			//}

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





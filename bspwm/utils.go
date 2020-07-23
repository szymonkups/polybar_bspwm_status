package bspwm

import (
	"encoding/json"
	"fmt"
	"github.com/szymonkups/polybar_bspwm_status/utils"
	"log"
	"strconv"
)

func GetLeafNodesOnDesktop(desktopId int) []int {
	leafs, err := utils.ExecuteCommand(nil, "bspc", "query", "-N", "-d", fmt.Sprintf("0x%08X", desktopId), "-n", ".leaf")

	if err != nil {
		return []int{}
	}

	var ids []int

	for _, leaf := range leafs {
		i, _ := strconv.ParseInt(leaf, 0, 64)
		ids = append(ids, int(i))
	}

	return ids
}

func GetMonitorInfo(monitorIndex int) (*MonitorInfo, error) {
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

	var i MonitorInfo
	err = json.Unmarshal(monitorData, &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}


func IsMonitorFocused(monitorId int) bool {
	focusedMonitorInfo, err := GetMonitorInfo(-1)

	if err != nil {
		log.Fatal(err)
	}

	return focusedMonitorInfo.Id == monitorId
}
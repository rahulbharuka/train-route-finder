package repository

import (
	"fmt"
)

// preparePath preapares path to passed destination.
func (h *handlerImpl) preparePath(dst int, prev []int, route *[]int) error {
	if prev[dst] >= 0 {
		h.preparePath(prev[dst], prev, route)
	}
	*route = append(*route, dst)
	return nil
}

// prepareResponse prepares and returns detailed route(s) in string format.
func (h *handlerImpl) prepareRouteSteps(route []int) string {
	var resp string
	prevStationName := stationIndexMap[route[0]].name
	newStationName := stationIndexMap[route[0]].name
	prevTrainLine := ""
	for i := 0; i+1 < len(route); i++ {
		newTrainLine := h.getTrainLine(route[i], route[i+1])

		if prevTrainLine != "" && newTrainLine != prevTrainLine {
			resp = resp + fmt.Sprintf("Take %v line from %v to %v. ", prevTrainLine, prevStationName, newStationName)
			resp = resp + fmt.Sprintf("Change from %v line to %v line. ", prevTrainLine, newTrainLine)
			prevStationName = newStationName
		}
		prevTrainLine = newTrainLine
		newStationName = stationIndexMap[route[i+1]].name
	}
	resp = resp + fmt.Sprintf("Take %v line from %v to %v.", prevTrainLine, prevStationName, newStationName)
	return resp
}

// getTrainLine returns the train line code between two stations.
func (h *handlerImpl) getTrainLine(i, j int) string {
	if len(stationIndexMap[i].codes) == 1 {
		return stationIndexMap[i].codes[0][:2]
	}
	if len(stationIndexMap[j].codes) == 1 {
		return stationIndexMap[j].codes[0][:2]
	}

	for _, iCode := range stationIndexMap[i].codes {
		for _, jCode := range stationIndexMap[j].codes {
			if iCode[:2] == jCode[:2] {
				return iCode[:2]
			}
		}
	}
	return ""
}

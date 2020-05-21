package repository

import "strconv"

type byStationCode []string

func (arr byStationCode) Len() int {
	return len(arr)
}

func (arr byStationCode) Less(i, j int) bool {
	val1, _ := strconv.Atoi(arr[i][2:])
	val2, _ := strconv.Atoi(arr[j][2:])
	return val1 < val2
}

func (arr byStationCode) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

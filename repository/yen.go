package repository

import (
	"math"
	"sort"
	"time"

	"github.com/rahulbharuka/train-route-finder/types"
)

// potential is an object for potential shortest path.
type potential struct {
	dist int
	path []int
}

// Yen returns top-k shortest path from src to dst using Yen's algorithm.
func (h *handlerImpl) yen(adjMatrix adjacencyMatrix, src int, dst int, topK int, journeyTime time.Time, computeTimeCost bool) ([]int, [][]int, error) {
	var potentials []potential
	distTopK := make([]int, topK)
	pathTopK := make([][]int, topK)
	for i := 0; i < topK; i++ {
		distTopK[i] = math.MaxInt32
	}

	ht := types.HTInvalid
	if computeTimeCost {
		ht = types.GetHourType(journeyTime)
	}

	// find the first shortest path
	dist, dijkstraPrev, err := h.dijkstra(adjMatrix, src, dst, ht, computeTimeCost)
	if err != nil {
		return nil, nil, err
	}
	distTopK[0] = dist
	path := []int{}
	h.prepareDijkstraPath(dst, dijkstraPrev, &path)
	pathTopK[0] = path // store first shortest path

	// now run Yen's algorithm for topK-1 times
	for k := 1; k < topK; {
		for i := 0; i < len(pathTopK[k-1])-1; i++ {
			for j := 0; j < k; j++ {
				if isShareRootPath(pathTopK[j], pathTopK[k-1][:i+1]) {
					adjMatrix[pathTopK[j][i]][pathTopK[j][i+1]].disabled = true
				}
			}
			h.disablePath(adjMatrix, pathTopK[k-1][:i])

			dist, dijkstraPrev, _ := h.dijkstra(adjMatrix, pathTopK[k-1][i], dst, ht, computeTimeCost)
			if dist != math.MaxInt32 {
				sPath := []int{}
				h.prepareDijkstraPath(dst, dijkstraPrev, &sPath)
				spurPath := mergePath(pathTopK[k-1][:i], sPath)
				spurWeight := h.getPathWeight(adjMatrix, spurPath, ht, computeTimeCost)
				existed := false
				for _, each := range potentials {
					if isSamePath(each.path, spurPath) {
						existed = true
						break
					}
				}
				if !existed {
					potentials = append(potentials, potential{
						spurWeight,
						spurPath,
					})
				}
			}

			h.reset(adjMatrix)
		}

		if len(potentials) == 0 {
			break
		}
		sort.Slice(potentials, func(i, j int) bool {
			return potentials[i].dist < potentials[j].dist
		})

		if len(potentials) >= topK-k {
			for l := 0; k < topK; l++ {
				distTopK[k] = potentials[l].dist
				pathTopK[k] = potentials[l].path
				k++
			}
			break
		} else {
			distTopK[k] = potentials[0].dist
			pathTopK[k] = potentials[0].path
			potentials = potentials[1:]
			k++
		}
	}

	return distTopK, pathTopK, nil
}

// DisablePath disables all the vertices in the path for further calculation.
func (h *handlerImpl) disablePath(adjMatrix adjacencyMatrix, path []int) {
	for _, vertex := range path {
		h.disableVertex(adjMatrix, vertex)
	}
}

// DisableVertex disables the vertex for further calculation.
func (h *handlerImpl) disableVertex(adjMatrix adjacencyMatrix, vertex int) {
	for to := range adjMatrix[vertex] {
		adjMatrix[vertex][to].disabled = true
	}
}

// getPathWeight returns weight of given path.
func (h *handlerImpl) getPathWeight(adjMatrix adjacencyMatrix, path []int, ht types.HourType, computeTimeCost bool) int {
	if len(path) == 0 {
		return math.MinInt32
	}

	if _, ok := adjMatrix[path[0]]; !ok {
		return math.MinInt32
	}

	pathWeight := 0
	prevLine := ""
	nextLine := ""
	for i := 0; i < len(path)-1; i++ {
		if _, ok := adjMatrix[path[i+1]]; !ok {
			return math.MinInt32
		}

		if _, ok := adjMatrix[path[i]][path[i+1]]; ok {
			interchangeCost := 0
			if computeTimeCost {
				nextLine = h.getTrainLine(path[i], path[i+1])
				if prevLine != "" && nextLine != prevLine {
					interchangeCost = interchangeCostMap[ht]
				}
				prevLine = nextLine
			}
			pathWeight = pathWeight + h.getEdgeWeight(adjMatrix, path[i], path[i+1], ht) + interchangeCost
		} else {
			return math.MaxInt32
		}
	}

	return pathWeight
}

// getEdgeWeight returns weight of an edge.
func (h *handlerImpl) getEdgeWeight(adjMatrix adjacencyMatrix, i, j int, ht types.HourType) int {
	var weight int
	switch ht {
	case types.HTNonPeak:
		weight = adjMatrix[i][j].weight.nonPeakHour
	case types.HTPeak:
		weight = adjMatrix[i][j].weight.peakHour
	case types.HTNight:
		weight = adjMatrix[i][j].weight.nightHour
	default:
		weight = adjMatrix[i][j].weight.defaults
	}
	return weight
}

// reset enables all vertices and edges for further calculation.
func (h *handlerImpl) reset(adjMatrix adjacencyMatrix) {
	for from := range adjMatrix {
		for to := range adjMatrix[from] {
			adjMatrix[from][to].disabled = false
		}
	}
}

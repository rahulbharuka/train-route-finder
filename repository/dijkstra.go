package repository

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/rahulbharuka/train-route-finder/types"
)

// prevVertex objects stores previous vertex attributes.
type prevVertex struct {
	stationIdx int
	line       string
}

// dijkstra finds shortest path from src to dst using Dijkstra's algorithm.
func (h *handlerImpl) dijkstra(adjMatrix adjacencyMatrix, src int, dst int, ht types.HourType, computeTimeCost bool) (int, map[int]*prevVertex, error) {
	if _, ok := adjMatrix[src]; !ok {
		return math.MaxInt32, nil, fmt.Errorf("Vertex %v does not exist", src)
	}

	prevMap := make(map[int]*prevVertex)     // map to store previous node of a minHeapNode.
	minHeap := minHeap{}                     // min heap to find unvisited vertext with min distance.
	minHeapNodeMap := map[int]*minHeapNode{} // heap index to node map

	// initialize the min-heap
	i := 0
	for stationIdx := range adjMatrix {
		prevMap[stationIdx] = &prevVertex{stationIdx: -1}
		if stationIdx != src {
			n := &minHeapNode{stationIdx: stationIdx, dist: math.MaxInt32, idx: i}
			minHeap = append(minHeap, n)
			minHeapNodeMap[stationIdx] = n
		} else {
			n := &minHeapNode{stationIdx: stationIdx, dist: 0, idx: i}
			minHeap = append(minHeap, n)
			minHeapNodeMap[stationIdx] = n
		}
		i++
	}
	heap.Init(&minHeap)

	// now run Dijkstra's algorithm
	for minHeap.Len() != 0 {
		fromNode := heap.Pop(&minHeap).(*minHeapNode)
		if fromNode.dist == math.MaxInt32 {
			// As current minHeapNode is unreachable, 'dst' will also be unreachable. So break early.
			return math.MaxInt32, nil, ErrRouteNotFound
		}

		from := fromNode.stationIdx
		if from == dst {
			// route found. so break early.
			return fromNode.dist, prevMap, nil
		}

		// update distance for every vertex directly reachable from current vertex.
		for to, edge := range adjMatrix[from] {
			if edge.disabled {
				continue // edge is disabled; so skip it.
			}
			weight := edge.weight.defaults
			interchangeCost := 0
			nextLine := h.getTrainLine(from, to)

			if computeTimeCost {
				weight = h.getEdgeWeight(adjMatrix, from, to, ht)

				if prevMap[from].line != "" && nextLine != prevMap[from].line {
					interchangeCost = interchangeCostMap[ht]
				}
			}

			toNode := minHeapNodeMap[to]
			newDist := fromNode.dist + weight + interchangeCost
			if newDist < toNode.dist {
				toNode.dist = newDist
				heap.Fix(&minHeap, toNode.idx)
				prevMap[to].stationIdx = from
				prevMap[to].line = nextLine
			}
		}
	}

	return minHeapNodeMap[dst].dist, prevMap, nil
}

// prepareDijkstraPath preapares path to passed destination.
func (h *handlerImpl) prepareDijkstraPath(dst int, prevMap map[int]*prevVertex, route *[]int) error {
	if prevMap[dst].stationIdx >= 0 {
		h.prepareDijkstraPath(prevMap[dst].stationIdx, prevMap, route)
	}
	*route = append(*route, dst)
	return nil
}

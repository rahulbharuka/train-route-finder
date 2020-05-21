package repository

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	// ErrRouteNotFound ...
	ErrRouteNotFound = errors.New("no route exist")
	// ErrInvalidRequest ...
	ErrInvalidRequest = errors.New("invalid request")
)

// Route is the route response object
type Route struct {
	Heading string `json:"heading"`
	Steps   string `json:"steps"`
}

// Handler is the repository handler interface
type Handler interface {
	FindRoutes(source string, destination string, journeyTime time.Time, computeTimeCost bool) ([]*Route, error)
}

// handlerImpl is a implementation of Handler interface
type handlerImpl struct{}

// GetHandler initializes and returns the repository layer handler.
func GetHandler() Handler {
	return &handlerImpl{}
}

// FindRoutes find shortest top-k routes from source to destionation.
func (h *handlerImpl) FindRoutes(source string, destination string, journeyTime time.Time, computeTimeCost bool) ([]*Route, error) {
	srcStation, ok1 := stationNameMap[source]
	dstStation, ok2 := stationNameMap[destination]
	if !ok1 || !ok2 {
		log.Println("invalid source or destination station")
		return nil, ErrInvalidRequest
	}

	dist, prev, err := h.yen(createAdjacencyMatrixCopy(), srcStation.idx, dstStation.idx, topK, journeyTime, computeTimeCost)
	if err != nil {
		log.Printf("failed to find route from %v to dst %v, err: %v", srcStation.name, dstStation.name, err)
		return nil, err
	}

	resp := make([]*Route, len(dist))

	var headingTemplate string
	if computeTimeCost {
		headingTemplate = "Expected Travel time: %v"
	} else {
		headingTemplate = "Number of stops to destination: %v"
	}
	// preapare response
	for i := 0; i < len(dist); i++ {
		resp[i] = &Route{
			Heading: fmt.Sprintf(headingTemplate, dist[i]),
			Steps:   h.prepareRouteSteps(prev[i]),
		}
	}

	return resp, nil
}

// createAdjacencyMatrix returns a deep copy (except edge weights) of adjacency matrix.
func createAdjacencyMatrixCopy() adjacencyMatrix {
	adjCopy := make(adjacencyMatrix)
	for from, adjacency := range railNetworkAdjacencyMatrix {
		adjMap := make(map[int]*edge)
		for to, e := range adjacency {
			adjMap[to] = &edge{
				disabled: e.disabled,
				weight:   e.weight, // edge weight is not modified during route calculation. So separate copy is not needed.
			}
		}
		adjCopy[from] = adjMap
	}
	return adjCopy
}

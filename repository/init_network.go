package repository

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/rahulbharuka/train-route-finder/types"
)

var (
	// auxillary data structures with package wide usage.
	stationNameMap             = map[string]*station{}    // maps a station-name to station.
	stationIndexMap            = map[int]*station{}       // maps station index to station.
	interchangeCostMap         = map[types.HourType]int{} // map of hourtype to interchange cost
	railNetworkAdjacencyMatrix = adjacencyMatrix{}        // graph of whole train network.
	topK                       int                        // max number of shortest routes to return

	// temporary data structures for initialization
	lineStationMap = map[string]*lineStation{} // maps stationCode to line-station.
	lineCostMap    = map[string][3]int{}       // map of train line to travel time cost per station
)

type adjacencyMatrix map[int]map[int]*edge

// edge object stores attributes of an edge
type edge struct {
	disabled bool    // whether its enabled or disabled. default: False
	weight   *weight // weights for the edge
}

type weight struct {
	nonPeakHour int // non-peak hours cost
	peakHour    int // peak hours cost
	nightHour   int // night hours cost
	defaults    int // default cost. Its always 1.
}

// lineStation is an object for a station on a specific line.
type lineStation struct {
	name           string
	openingDate    time.Time
	neighbours     map[int]*edge
	lineStationIdx int
}

// station stores attributes of a station
type station struct {
	name  string
	codes []string
	idx   int
}

// RailNetworkInit initializes the rail network by reading specified csv file.
func RailNetworkInit() {

	// read station-map file
	trainLines := readStationMapFile()

	// read trainline cost file
	readTrainlineCostFile()

	// read interchange cost file
	readInterchangeCostFile()

	// populate neighbours for every line-station.
	populateNeighbours(trainLines)

	// create adjacency matrix
	initRailNetworkAdjacencyMatrix()

	// set topK value
	setTopKValue()
}

// readStationMapFile reads station-map file and initializes the multiple auxillary data structures.
func readStationMapFile() map[string][]string {
	csvFile := os.Getenv("STATION_MAP_FILE")
	if csvFile == "" {
		log.Fatal("$STATION_MAP_FILE must be set")
	}

	file, err := os.Open(csvFile)
	if err != nil {
		panic("error while opening file")
	}

	r := csv.NewReader(file)

	trainLines := map[string][]string{} // temporary map of line to stations. Its used to order stations on a line.

	line := -1
	stationIdx := 0 // auto generated index for every station.
	// Iterate through the records
	for {
		line++
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line == 0 {
			// skip the header line
			continue
		}

		stationCode := record[0]
		stationName := record[1]
		openingTime, err := time.Parse("2 January 2006", record[2])
		if err != nil {
			log.Printf("wrong date format for station %v. Skipping it\n", stationName)
			continue
		}

		if _, ok := lineStationMap[stationCode]; ok {
			log.Println("duplicate station entry for stationCode ", stationCode)
			continue
		}

		// create line-station for the record.
		lineStationMap[stationCode] = &lineStation{
			name:        stationName,
			openingDate: openingTime,
			neighbours:  map[int]*edge{},
		}

		// if its an existing station, just append station code. Otherwise, create a new station object.
		if s, ok := stationNameMap[stationName]; ok {
			s.codes = append(s.codes, stationCode)
		} else {
			station := &station{
				name:  stationName,
				codes: []string{stationCode},
				idx:   stationIdx,
			}
			stationNameMap[stationName] = station
			stationIndexMap[stationIdx] = station
			stationIdx++
		}

		lineCode := stationCode[:2] // extract train line.
		if _, ok := trainLines[lineCode]; ok {
			trainLines[lineCode] = append(trainLines[lineCode], stationCode)
		} else {
			trainLines[lineCode] = []string{stationCode}
		}
	}
	return trainLines
}

// readTrainlineCostFile reads trainline-cost file and initializes lineCostMap.
func readTrainlineCostFile() {
	csvFile := os.Getenv("TRAINLINE_COST_FILE")
	if csvFile == "" {
		log.Fatal("$TRAINLINE_COST_FILE must be set")
	}

	file, err := os.Open(csvFile)
	if err != nil {
		panic("error while opening $TRAINLINE_COST_FILE file")
	}

	r := csv.NewReader(file)

	line := -1
	// Iterate through the records
	for {
		line++
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line == 0 {
			continue // skip the header line
		}

		nonPeakHoursCost, err := strconv.Atoi(record[1])
		if err != nil {
			panic("invalid non-peak hours travel time cost")
		}
		if nonPeakHoursCost == -1 {
			nonPeakHoursCost = math.MaxInt32
		}

		peakHoursCost, err := strconv.Atoi(record[2])
		if err != nil {
			panic("invalid peak hours travel time cost")
		}
		if peakHoursCost == -1 {
			peakHoursCost = math.MaxInt32
		}

		nightHoursCost, err := strconv.Atoi(record[3])
		if err != nil {
			panic("invalid night hours travel time cost")
		}
		if nightHoursCost == -1 {
			nightHoursCost = math.MaxInt32
		}
		lineCostMap[record[0]] = [...]int{nonPeakHoursCost, peakHoursCost, nightHoursCost}
	}
}

// readInterchangeCostFile reads interchange-cost file and initializes interchangeCostMap.
func readInterchangeCostFile() {
	csvFile := os.Getenv("INTERCHANGE_COST_FILE")
	if csvFile == "" {
		log.Fatal("$INTERCHANGE_COST_FILE must be set")
	}

	file, err := os.Open(csvFile)
	if err != nil {
		panic("error while opening $INTERCHANGE_COST_FILE file")
	}

	r := csv.NewReader(file)

	line := -1
	// Iterate through the records
	for {
		line++
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line == 0 {
			continue // skip the header line
		}

		cost, err := strconv.Atoi(record[1])
		if err != nil || cost < 0 {
			panic("invalid interchange time cost")
		}
		interchangeCostMap[types.ConvertToHourType(record[0])] = cost
	}
}

// populateNeighbours populates neighbours (prev and next) for every line-station.
func populateNeighbours(trainLines map[string][]string) {
	for lineCode, stationCodes := range trainLines {
		// to find neighbours, first sort line stations by line-station number.
		sort.Sort(byStationCode(stationCodes))

		edge := createEdge(lineCode)

		for i, stationCode := range stationCodes {
			lineStationMap[stationCode].lineStationIdx = i

			if i-1 >= 0 {
				lineStationMap[stationCode].neighbours[stationNameMap[lineStationMap[stationCodes[i-1]].name].idx] = edge
			}
			if i+1 < len(stationCodes) {
				lineStationMap[stationCode].neighbours[stationNameMap[lineStationMap[stationCodes[i+1]].name].idx] = edge
			}
		}
	}
}

// createEdge creates an edge (connection) for a given train line.
func createEdge(lineCode string) *edge {
	lineCosts, ok := lineCostMap[lineCode]
	if !ok {
		log.Printf("Trainline %v cost is not available\n", lineCode)
		panic("trainline cost not available")
	}

	return &edge{
		disabled: false,
		weight: &weight{
			nonPeakHour: lineCosts[0],
			peakHour:    lineCosts[1],
			nightHour:   lineCosts[2],
			defaults:    1,
		},
	}
}

// initRailNetworkAdjacencyMatrix initializes adjacency matrix of given rail network.
func initRailNetworkAdjacencyMatrix() {
	for _, station := range stationNameMap {
		adj := map[int]*edge{}
		for _, stationCode := range station.codes {
			for k, v := range lineStationMap[stationCode].neighbours {
				adj[k] = v
			}
		}
		railNetworkAdjacencyMatrix[station.idx] = adj
	}
}

// setTopKValue sets the topK value configured from environment variable.
// If environment variable is not set or invalid, it falls back to default value of 1.
func setTopKValue() {
	var err error
	maxRoutes := os.Getenv("MAX_ROUTES")
	topK, err = strconv.Atoi(maxRoutes)
	if err != nil {
		log.Println("invalid MAX_ROUTES. Falling back to default value of 1")
		topK = 1
	}
}
